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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"reflect"
	"strings"
	"time"
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

// ConvertDeployments should be split into two task theoretically
// But in BitBucket, all deployments have commits, and we use "LEFT JOIN" to get "ref_name" only, so there is no need to change it.
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

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(bitbucketDeploymentWithRefName{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			bitbucketDeployment := inputRow.(*bitbucketDeploymentWithRefName)

			var duration *float64
			if bitbucketDeployment.CompletedOn != nil && bitbucketDeployment.StartedOn != nil {
				d := float64(bitbucketDeployment.CompletedOn.Sub(*bitbucketDeployment.StartedOn).Milliseconds() / 1e3)
				duration = &d
			}
			createdAt := time.Now()
			if bitbucketDeployment.CreatedOn != nil {
				createdAt = *bitbucketDeployment.CreatedOn
			}
			domainDeployCommit := &devops.CicdDeploymentCommit{
				DomainEntity: domainlayer.DomainEntity{
					Id: idGen.Generate(data.Options.ConnectionId, bitbucketDeployment.BitbucketId),
				},
				CicdScopeId: repoId,
				Name:        bitbucketDeployment.Name,
				Result: devops.GetResult(&devops.ResultRule{
					Success: []string{models.COMPLETED, models.SUCCESSFUL},
					Failure: []string{models.FAILED, models.STOPPED, models.CANCELLED},
					Default: devops.RESULT_DEFAULT,
				}, bitbucketDeployment.Status),
				Status: devops.GetStatus(&devops.StatusRule{
					Done:       []string{models.COMPLETED, models.SUCCESSFUL, models.FAILED, models.STOPPED, models.CANCELLED},
					InProgress: []string{models.IN_PROGRESS},
					Default:    devops.STATUS_OTHER,
				}, bitbucketDeployment.Status),
				OriginalStatus:      bitbucketDeployment.Status,
				Environment:         strings.ToUpper(bitbucketDeployment.Environment), // or bitbucketDeployment.EnvironmentType, they are same so far.
				OriginalEnvironment: strings.ToUpper(bitbucketDeployment.Environment),
				TaskDatesInfo: devops.TaskDatesInfo{
					CreatedDate:  createdAt,
					StartedDate:  bitbucketDeployment.StartedOn,
					FinishedDate: bitbucketDeployment.CompletedOn,
				},
				DurationSec: duration,
				CommitSha:   bitbucketDeployment.CommitSha,
				RefName:     bitbucketDeployment.RefName,
				RepoId:      repoId,
				RepoUrl:     repo.HTMLUrl,
			}
			if domainDeployCommit.Environment == devops.TEST {
				// Theoretically, environment cannot be "Test" according to
				// https://developer.atlassian.com/server/bitbucket/rest/v814/api-group-builds-and-deployments/#api-api-latest-projects-projectkey-repos-repositoryslug-commits-commitid-deployments-get
				// but in practice, we found environment is "Test".
				// So convert it to DevLake's definition.
				domainDeployCommit.Environment = devops.TESTING
			}

			domainDeployCommit.CicdDeploymentId = domainDeployCommit.Id
			return []interface{}{domainDeployCommit, domainDeployCommit.ToDeployment()}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
