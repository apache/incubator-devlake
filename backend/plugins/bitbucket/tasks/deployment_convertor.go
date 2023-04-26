/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

var ConvertiDeploymentMeta = plugin.SubTaskMeta{
	Name:             "convertDeployments",
	EntryPoint:       ConvertDeployments,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bitbucket_deployment into domain layer tables",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

type bitbucketDeploymentWithRefName struct {
	models.BitbucketDeployment
	RefName string
}

func ConvertDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)
	db := taskCtx.GetDal()

	repo := &models.BitbucketRepo{}
	err := db.First(repo, dal.Where("connection_id = ? AND bitbucket_id = ?", data.Options.ConnectionId, data.Options.FullName))
	if err != nil {
		return err
	}
	repoId := didgen.NewDomainIdGenerator(&models.BitbucketRepo{}).Generate(data.Options.ConnectionId, repo.BitbucketId)

	cursor, err := db.Cursor(
		dal.Select("d.*, p.ref_name"),
		dal.From("_tool_bitbucket_deployments d"),
		dal.Join("LEFT JOIN _tool_bitbucket_pipelines p ON (p.connection_id = d.connection_id AND p.bitbucket_id = d.pipeline_id)"),
		dal.Where("d.connection_id = ? AND p.repo_id = ? ", data.Options.ConnectionId, data.Options.FullName),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	idGen := didgen.NewDomainIdGenerator(&models.BitbucketDeployment{})
	pipelineIdGen := didgen.NewDomainIdGenerator(&models.BitbucketPipeline{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(bitbucketDeploymentWithRefName{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			bitbucketDeployment := inputRow.(*bitbucketDeploymentWithRefName)

			var duration *uint64
			if bitbucketDeployment.CompletedOn != nil {
				d := uint64(bitbucketDeployment.CompletedOn.Sub(*bitbucketDeployment.StartedOn).Seconds())
				duration = &d
			}
			domainDeployCommit := &devops.CicdDeploymentCommit{
				DomainEntity: domainlayer.DomainEntity{
					Id: idGen.Generate(data.Options.ConnectionId, bitbucketDeployment.BitbucketId),
				},
				CicdScopeId:      repoId,
				CicdDeploymentId: pipelineIdGen.Generate(data.Options.ConnectionId, bitbucketDeployment.PipelineId),
				Name:             bitbucketDeployment.Name,
				Result: devops.GetResult(&devops.ResultRule{
					Failed:  []string{"UNDEPLOYED"},
					Success: []string{"COMPLETED"},
					Default: "",
				}, bitbucketDeployment.Status),
				Status: devops.GetStatus(&devops.StatusRule{
					Done:    []string{"COMPLETED", "UNDEPLOYED"},
					Default: devops.IN_PROGRESS,
				}, bitbucketDeployment.Status),
				Environment:  bitbucketDeployment.Environment,
				CreatedDate:  *bitbucketDeployment.CreatedOn,
				StartedDate:  bitbucketDeployment.StartedOn,
				FinishedDate: bitbucketDeployment.CompletedOn,
				DurationSec:  duration,
				CommitSha:    bitbucketDeployment.CommitSha,
				RefName:      bitbucketDeployment.RefName,
				RepoId:       repoId,
				RepoUrl:      repo.HTMLUrl,
			}
			return []interface{}{domainDeployCommit}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
