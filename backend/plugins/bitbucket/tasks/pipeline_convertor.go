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
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

var ConvertPipelineMeta = plugin.SubTaskMeta{
	Name:             "convertPipelines",
	EntryPoint:       ConvertPipelines,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bitbucket_pipeline into domain layer table pipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ConvertPipelines(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)
	db := taskCtx.GetDal()

	repo := &models.BitbucketRepo{}
	err := db.First(repo, dal.Where("connection_id = ? AND bitbucket_id = ?", data.Options.ConnectionId, data.Options.FullName))
	if err != nil {
		return err
	}
	repoId := didgen.NewDomainIdGenerator(&models.BitbucketRepo{}).Generate(repo.ConnectionId, repo.BitbucketId)

	cursor, err := db.Cursor(
		dal.From(models.BitbucketPipeline{}),
		dal.Where("connection_id = ? AND repo_id = ?", data.Options.ConnectionId, data.Options.FullName),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	pipelineIdGen := didgen.NewDomainIdGenerator(&models.BitbucketPipeline{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.BitbucketPipeline{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			bitbucketPipeline := inputRow.(*models.BitbucketPipeline)

			createdAt := time.Now()
			if bitbucketPipeline.BitbucketCreatedOn != nil {
				createdAt = *bitbucketPipeline.BitbucketCreatedOn
			}
			results := make([]interface{}, 0, 2)
			domainPipelineCommit := &devops.CiCDPipelineCommit{
				PipelineId: pipelineIdGen.Generate(data.Options.ConnectionId, bitbucketPipeline.BitbucketId),
				RepoId:     repoId,
				CommitSha:  bitbucketPipeline.CommitSha,
				Branch:     bitbucketPipeline.RefName,
				RepoUrl:    repo.HTMLUrl,
			}
			domainPipeline := &devops.CICDPipeline{
				DomainEntity: domainlayer.DomainEntity{
					Id: pipelineIdGen.Generate(data.Options.ConnectionId, bitbucketPipeline.BitbucketId),
				},
				Name: didgen.NewDomainIdGenerator(&models.BitbucketPipeline{}).
					Generate(data.Options.ConnectionId, bitbucketPipeline.RefName),
				Result: devops.GetResult(&devops.ResultRule{
					Failed:  []string{models.FAILED, models.ERROR, models.UNDEPLOYED},
					Abort:   []string{models.STOPPED, models.SKIPPED},
					Success: []string{models.SUCCESSFUL, models.COMPLETED},
					Manual:  []string{models.PAUSED, models.HALTED},
					Default: devops.SUCCESS,
				}, bitbucketPipeline.Result),
				Status: devops.GetStatus(&devops.StatusRule{
					InProgress: []string{models.IN_PROGRESS, models.PENDING, models.BUILDING},
					Default:    devops.DONE,
				}, bitbucketPipeline.Status),
				Type:         bitbucketPipeline.Type,
				Environment:  bitbucketPipeline.Environment,
				CreatedDate:  createdAt,
				DurationSec:  bitbucketPipeline.DurationInSeconds,
				FinishedDate: bitbucketPipeline.BitbucketCompleteOn,
				CicdScopeId:  repoId,
			}
			results = append(results, domainPipelineCommit, domainPipeline)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
