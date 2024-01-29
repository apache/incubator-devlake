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
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

var ConvertDeployBuildsToDeploymentMeta = plugin.SubTaskMeta{
	Name:             "convertDeployBuildsToDeployments",
	EntryPoint:       ConvertDeployBuildsToDeployments,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bamboo_deploy_builds into domain layer table cicd_deployments",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

type bambooDeployBuildEx struct {
	models.BambooDeployBuild
	ProjectPlanName string
	ProjectName     string
}

func ConvertDeployBuildsToDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_JOB_BUILD_TABLE)
	cursor, err := db.Cursor(
		dal.Select("db.*, p.name as project_plan_name, p.project_name"),
		dal.From("_tool_bamboo_deploy_builds AS db"),
		dal.Join("LEFT JOIN _tool_bamboo_plans as p ON db.plan_key = p.plan_key"),
		dal.Where("db.connection_id = ? and db.plan_key = ?", data.Options.ConnectionId, data.Options.PlanKey),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	planIdGen := didgen.NewDomainIdGenerator(&models.BambooPlan{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(bambooDeployBuildEx{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			input := inputRow.(*bambooDeployBuildEx)
			deploymentCommitId := didgen.NewDomainIdGenerator(&bambooDeployBuildEx{}).Generate(data.Options.ConnectionId, input.DeployBuildId)
			createdDate := time.Now()
			if input.StartedDate != nil {
				createdDate = *input.StartedDate
			}
			name := input.DeploymentVersionName
			if input.ProjectPlanName != "" {
				name = fmt.Sprintf("%s/%s", input.ProjectPlanName, input.DeploymentVersionName)
			}

			deployment := &devops.CICDDeployment{
				DomainEntity: domainlayer.DomainEntity{
					Id: deploymentCommitId,
				},
				CicdScopeId: planIdGen.Generate(data.Options.ConnectionId, data.Options.PlanKey),
				Name:        name,
				Result: devops.GetResult(&devops.ResultRule{
					Success: []string{ResultSuccess, ResultSuccessful},
					Failure: []string{ResultFailed},
					Default: devops.RESULT_DEFAULT,
				}, input.DeploymentState),
				OriginalResult: input.DeploymentState,
				Status: devops.GetStatus(&devops.StatusRule{
					Done:       []string{StatusFinished},
					InProgress: []string{StatusInProgress, StatusPending, StatusQueued},
					Default:    devops.STATUS_OTHER,
				}, input.LifeCycleState),
				OriginalStatus:      input.LifeCycleState,
				Environment:         input.Environment,
				OriginalEnvironment: input.Environment,
				TaskDatesInfo: devops.TaskDatesInfo{
					CreatedDate:  createdDate,
					QueuedDate:   input.QueuedDate,
					StartedDate:  input.ExecutedDate,
					FinishedDate: input.FinishedDate,
				},
			}
			if data.RegexEnricher.ReturnNameIfMatched(devops.ENV_NAME_PATTERN, input.Environment) != "" {
				deployment.Environment = devops.PRODUCTION
			}
			if input.FinishedDate != nil && input.ExecutedDate != nil {
				duration := float64(input.FinishedDate.Sub(*input.ExecutedDate).Milliseconds() / 1e3)
				deployment.DurationSec = &duration
			}
			return []interface{}{deployment}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
