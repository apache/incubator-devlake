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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

var ConvertPipelineStepMeta = plugin.SubTaskMeta{
	Name:             "Convert Pipeline Steps",
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

			// don't save to domain layer if `StartedOn` is nil
			if bitbucketPipelineStep.StartedOn == nil {
				return nil, nil
			}

			createdAt := time.Now()
			if bitbucketPipelineStep.StartedOn != nil {
				createdAt = *bitbucketPipelineStep.StartedOn
			}
			domainTask := &devops.CICDTask{
				DomainEntity: domainlayer.DomainEntity{
					Id: pipelineStepIdGen.Generate(data.Options.ConnectionId, bitbucketPipelineStep.BitbucketId),
				},
				Name:       bitbucketPipelineStep.Name,
				PipelineId: pipelineIdGen.Generate(data.Options.ConnectionId, bitbucketPipelineStep.PipelineId),
				Result: devops.GetResult(&devops.ResultRule{
					Success: []string{models.SUCCESSFUL, models.COMPLETED},
					Failure: []string{models.FAILED, models.ERROR, models.STOPPED},
					Default: devops.RESULT_DEFAULT,
				}, bitbucketPipelineStep.Result),
				OriginalResult: bitbucketPipelineStep.Result,
				Status: devops.GetStatus(&devops.StatusRule{
					Done:       []string{models.COMPLETED, models.SUCCESSFUL, models.FAILED, models.ERROR, models.STOPPED},
					InProgress: []string{models.IN_PROGRESS, models.PENDING, models.BUILDING, models.READY},
					Default:    devops.STATUS_OTHER,
				}, bitbucketPipelineStep.State),
				OriginalStatus: bitbucketPipelineStep.State,
				CicdScopeId:    repoIdGen.Generate(data.Options.ConnectionId, data.Options.FullName),
				DurationSec:    float64(bitbucketPipelineStep.DurationInSeconds),
				TaskDatesInfo: devops.TaskDatesInfo{
					CreatedDate:  createdAt,
					StartedDate:  bitbucketPipelineStep.StartedOn,
					FinishedDate: bitbucketPipelineStep.CompletedOn,
				},
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
