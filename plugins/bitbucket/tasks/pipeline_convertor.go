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
	"github.com/apache/incubator-devlake/errors"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	bitbucketModels "github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertPipelineMeta = core.SubTaskMeta{
	Name:             "convertPipelines",
	EntryPoint:       ConvertPipelines,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bitbucket_pipeline into domain layer table pipeline",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func ConvertPipelines(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)

	cursor, err := db.Cursor(dal.From(bitbucketModels.BitbucketPipeline{}))
	if err != nil {
		return err
	}
	defer cursor.Close()

	pipelineIdGen := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketPipeline{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(bitbucketModels.BitbucketPipeline{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_PIPELINE_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			bitbucketPipeline := inputRow.(*bitbucketModels.BitbucketPipeline)

			createdAt := time.Now()
			if bitbucketPipeline.BitbucketCreatedOn != nil {
				createdAt = *bitbucketPipeline.BitbucketCreatedOn
			}
			results := make([]interface{}, 0, 2)
			domainPipelineCommit := &devops.CiCDPipelineCommit{
				PipelineId: pipelineIdGen.Generate(data.Options.ConnectionId, bitbucketPipeline.BitbucketId),
				RepoId: didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketRepo{}).
					Generate(bitbucketPipeline.ConnectionId, bitbucketPipeline.RepoId),
				CommitSha: bitbucketPipeline.CommitSha,
				Branch:    bitbucketPipeline.RefName,
				RepoUrl:   bitbucketPipeline.WebUrl,
			}
			domainPipeline := &devops.CICDPipeline{
				DomainEntity: domainlayer.DomainEntity{
					Id: pipelineIdGen.Generate(data.Options.ConnectionId, bitbucketPipeline.BitbucketId),
				},
				Name: didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketPipeline{}).
					Generate(data.Options.ConnectionId, bitbucketPipeline.RefName),
				Result: devops.GetResult(&devops.ResultRule{
					Failed:  []string{bitbucketModels.FAILED, bitbucketModels.ERROR, bitbucketModels.UNDEPLOYED},
					Abort:   []string{bitbucketModels.STOPPED, bitbucketModels.SKIPPED},
					Success: []string{bitbucketModels.SUCCESSFUL, bitbucketModels.COMPLETED},
					Manual:  []string{bitbucketModels.PAUSED, bitbucketModels.HALTED},
					Default: devops.SUCCESS,
				}, bitbucketPipeline.Result),
				Status: devops.GetStatus(&devops.StatusRule{
					InProgress: []string{bitbucketModels.IN_PROGRESS, bitbucketModels.PENDING, bitbucketModels.BUILDING},
					Default:    devops.DONE,
				}, bitbucketPipeline.Status),
				Type:         "CI/CD",
				CreatedDate:  createdAt,
				DurationSec:  bitbucketPipeline.DurationInSeconds,
				FinishedDate: bitbucketPipeline.BitbucketCompleteOn,
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
