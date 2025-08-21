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
	"fmt"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

type JenkinsBuildWithRepoStage struct {
	// collected fields
	ConnectionId        uint64 `gorm:"primaryKey"`
	ID                  string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name                string `json:"name" gorm:"type:varchar(255)"`
	ExecNode            string `json:"execNode" gorm:"type:varchar(255)"`
	CommitSha           string `gorm:"type:varchar(255)"`
	Result              string // Result
	Status              string `json:"status" gorm:"type:varchar(255)"`
	StartTimeMillis     int64  `json:"startTimeMillis"`
	DurationMillis      int    `json:"durationMillis"`
	PauseDurationMillis int    `json:"pauseDurationMillis"`
	Type                string `gorm:"index;type:varchar(255)"`
	BuildName           string `gorm:"primaryKey;type:varchar(255)"`
	Branch              string `gorm:"type:varchar(255)"`
	RepoUrl             string `gorm:"type:varchar(255)"`
	common.NoPKModel
}

var ConvertStagesMeta = plugin.SubTaskMeta{
	Name:             "convertStages",
	EntryPoint:       ConvertStages,
	EnabledByDefault: true,
	Description:      "convert jenkins_stages",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ConvertStages(taskCtx plugin.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JenkinsTaskData)

	clauses := []dal.Clause{
		dal.Select(`tjb.connection_id, tjs.build_name, tjs.id, tjs._raw_data_remark, tjs.name,
			tjs._raw_data_id, tjs._raw_data_table, tjs._raw_data_params,
			tjs.status, tjs.start_time_millis, tjs.duration_millis,
			tjs.pause_duration_millis, tjs.type, tjb.result,
			tjb.triggered_by, tjb.building`),
		dal.From("_tool_jenkins_stages tjs"),
		dal.Join("left join _tool_jenkins_builds tjb on tjs.build_name = tjb.full_name"),
	}

	if data.Options.Class == WORKFLOW_MULTI_BRANCH_PROJECT {
		clauses = append(clauses,
			dal.Where(`tjb.connection_id = ? and tjb.full_name like ?`,
				data.Options.ConnectionId, fmt.Sprintf("%s%%", data.Options.JobFullName)))
	} else {
		clauses = append(clauses,
			dal.Where("tjb.connection_id = ? and tjb.job_path = ? and tjb.job_name = ? ",
				data.Options.ConnectionId, data.Options.JobPath, data.Options.JobName))
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	stageIdGen := didgen.NewDomainIdGenerator(&models.JenkinsStage{})
	buildIdGen := didgen.NewDomainIdGenerator(&models.JenkinsBuild{})
	jobIdGen := didgen.NewDomainIdGenerator(&models.JenkinsJob{})

	convertor, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(JenkinsBuildWithRepoStage{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_STAGE_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			body := inputRow.(*JenkinsBuildWithRepoStage)
			if body.Name == "" {
				return nil, err
			}
			var durationMillis int64
			if body.DurationMillis > 0 {
				durationMillis = int64(body.DurationMillis)
			} else {
				durationMillis = int64(0)
			}
			durationSec := float64(durationMillis / 1e3)
			jenkinsTaskResult := devops.GetResult(&devops.ResultRule{
				Success: []string{SUCCESS},
				Failure: []string{FAILED, FAILURE, ABORTED},
				Default: devops.RESULT_DEFAULT,
			}, body.Result)

			jenkinsTaskStatus := devops.GetStatus(&devops.StatusRule{
				Done:       []string{SUCCESS},
				InProgress: []string{},
				Default:    devops.STATUS_OTHER,
			}, body.Status)

			var jenkinsTaskFinishedDate *time.Time
			results := make([]interface{}, 0)
			finishedDateMillis := body.StartTimeMillis + durationMillis
			finishedDate := time.Unix(finishedDateMillis/1e3, (finishedDateMillis%1e3)*int64(time.Millisecond))
			jenkinsTaskFinishedDate = &finishedDate
			startedDate := time.Unix(body.StartTimeMillis/1e3, 0)

			jenkinsTask := &devops.CICDTask{
				DomainEntity: domainlayer.DomainEntity{
					Id: stageIdGen.Generate(body.ConnectionId, body.BuildName, body.ID),
				},
				Name:        body.Name,
				PipelineId:  buildIdGen.Generate(body.ConnectionId, body.BuildName),
				Result:      jenkinsTaskResult,
				Status:      jenkinsTaskStatus,
				DurationSec: durationSec,
				TaskDatesInfo: devops.TaskDatesInfo{
					CreatedDate:  startedDate,
					StartedDate:  &startedDate,
					FinishedDate: jenkinsTaskFinishedDate,
				},
				OriginalResult: body.Result,
				OriginalStatus: body.Status,
				CicdScopeId:    jobIdGen.Generate(body.ConnectionId, data.Options.JobFullName),
				Type:           data.RegexEnricher.ReturnNameIfMatched(devops.DEPLOYMENT, body.Name),
				Environment:    data.RegexEnricher.ReturnNameIfOmittedOrMatched(devops.PRODUCTION, body.Name),
			}
			// if the task is not executed, set the result to default, so that it will not be calculated in the dora
			if jenkinsTask.OriginalStatus == "NOT_EXECUTED" {
				jenkinsTask.Result = devops.RESULT_DEFAULT
			}
			results = append(results, jenkinsTask)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return convertor.Execute()
}
