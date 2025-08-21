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
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

var ConvertJobBuildsMeta = plugin.SubTaskMeta{
	Name:             "convertJobBuilds",
	EntryPoint:       ConvertJobBuilds,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bamboo_jobBuilds into  domain layer table jobBuilds",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ConvertJobBuilds(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_JOB_BUILD_TABLE)
	cursor, err := db.Cursor(
		dal.From(&models.BambooJobBuild{}),
		dal.Where("connection_id = ? and plan_key = ?", data.Options.ConnectionId, data.Options.PlanKey))
	if err != nil {
		return err
	}
	defer cursor.Close()

	jobBuildIdGen := didgen.NewDomainIdGenerator(&models.BambooJobBuild{})
	planBuildIdGen := didgen.NewDomainIdGenerator(&models.BambooPlanBuild{})
	planIdGen := didgen.NewDomainIdGenerator(&models.BambooPlan{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.BambooJobBuild{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			line := inputRow.(*models.BambooJobBuild)
			if line.BuildStartedTime == nil {
				return nil, nil
			}
			createdDate := time.Now()
			if line.BuildStartedTime != nil {
				createdDate = *line.BuildStartedTime
			}
			queuedDurationSec := float64(line.QueueDuration) / 1e3
			domainJobBuild := &devops.CICDTask{
				DomainEntity: domainlayer.DomainEntity{Id: jobBuildIdGen.Generate(data.Options.ConnectionId, line.JobBuildKey)},
				Name:         line.JobName,
				DurationSec:  float64(line.BuildDuration) / 1e3,
				TaskDatesInfo: devops.TaskDatesInfo{
					CreatedDate:  createdDate,
					QueuedDate:   line.QueueStartedTime,
					StartedDate:  line.BuildStartedTime,
					FinishedDate: line.BuildCompletedDate,
				},
				QueuedDurationSec: &queuedDurationSec,
				PipelineId:        planBuildIdGen.Generate(data.Options.ConnectionId, line.PlanBuildKey),
				CicdScopeId:       planIdGen.Generate(data.Options.ConnectionId, data.Options.PlanKey),
				Result: devops.GetResult(&devops.ResultRule{
					Success: []string{ResultSuccess, ResultSuccessful},
					Failure: []string{ResultFailed},
					Default: devops.RESULT_DEFAULT,
				}, line.BuildState),
				OriginalResult: line.BuildState,
				Status: devops.GetStatus(&devops.StatusRule{
					Done:       []string{StatusFinished},
					InProgress: []string{StatusInProgress, StatusPending, StatusQueued},
					Default:    devops.STATUS_OTHER,
				}, line.LifeCycleState),
				OriginalStatus: line.LifeCycleState,
			}

			domainJobBuild.Type = line.Type
			domainJobBuild.Environment = line.Environment

			return []interface{}{
				domainJobBuild,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
