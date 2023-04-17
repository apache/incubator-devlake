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

var ConvertPipelineStepMeta = plugin.SubTaskMeta{
	Name:             "convertPipelineSteps",
	EntryPoint:       ConvertPipelineSteps,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bitbucket_pipeline into domain layer table pipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ConvertPipelineSteps(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_STEPS_TABLE)
	db := taskCtx.GetDal()

	cursor, err := db.Cursor(
		dal.From(models.BitbucketPipelineStep{}),
		dal.Where("connection_id = ? AND repo_id = ?", data.Options.ConnectionId, data.Options.FullName),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	pipelineStepIdGen := didgen.NewDomainIdGenerator(&models.BitbucketPipelineStep{})
	pipelineIdGen := didgen.NewDomainIdGenerator(&models.BitbucketPipeline{})
	repoIdGen := didgen.NewDomainIdGenerator(&models.BitbucketRepo{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.BitbucketPipelineStep{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			bitbucketPipelineStep := inputRow.(*models.BitbucketPipelineStep)

			domainTask := &devops.CICDTask{
				DomainEntity: domainlayer.DomainEntity{
					Id: pipelineStepIdGen.Generate(data.Options.ConnectionId, bitbucketPipelineStep.BitbucketId),
				},
				Name:       bitbucketPipelineStep.Name,
				PipelineId: pipelineIdGen.Generate(data.Options.ConnectionId, bitbucketPipelineStep.PipelineId),
				Result: devops.GetResult(&devops.ResultRule{
					Failed:  []string{models.FAILED, models.ERROR, models.UNDEPLOYED},
					Abort:   []string{models.STOPPED, models.SKIPPED},
					Success: []string{models.SUCCESSFUL, models.COMPLETED},
					Manual:  []string{models.PAUSED, models.HALTED},
					Default: devops.SUCCESS,
				}, bitbucketPipelineStep.Result),
				Status: devops.GetStatus(&devops.StatusRule{
					InProgress: []string{models.IN_PROGRESS, models.PENDING, models.BUILDING},
					Default:    devops.DONE,
				}, bitbucketPipelineStep.State),
				CicdScopeId: repoIdGen.Generate(data.Options.ConnectionId, data.Options.FullName),
			}
			// not save to domain layer if StartedOn is empty
			if bitbucketPipelineStep.StartedOn == nil {
				return nil, nil
			}
			domainTask.StartedDate = *bitbucketPipelineStep.StartedOn
			// rebuild the FinishedDate
			if domainTask.Status == devops.DONE {
				domainTask.FinishedDate = bitbucketPipelineStep.CompletedOn
				domainTask.DurationSec = uint64(bitbucketPipelineStep.DurationInSeconds)
			}

			bitbucketDeployment := &models.BitbucketDeployment{}
			deploymentErr := db.First(bitbucketDeployment, dal.Where(`step_id=?`, bitbucketPipelineStep.BitbucketId))
			if deploymentErr == nil {
				domainTask.Type = devops.DEPLOYMENT
				if bitbucketDeployment.EnvironmentType == `Production` {
					domainTask.Environment = devops.PRODUCTION
				} else if bitbucketDeployment.EnvironmentType == `Staging` {
					domainTask.Environment = devops.STAGING
				} else if bitbucketDeployment.EnvironmentType == `Test` {
					domainTask.Environment = devops.TESTING
				}
			}
			return []interface{}{
				domainTask,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
