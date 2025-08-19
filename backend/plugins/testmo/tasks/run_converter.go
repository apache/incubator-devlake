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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/qa"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	testmoModels "github.com/apache/incubator-devlake/plugins/testmo/models"
)

var ConvertRunsMeta = plugin.SubTaskMeta{
	Name:             "convertRuns",
	EntryPoint:       ConvertRuns,
	EnabledByDefault: true,
	Description:      "Convert tool layer table testmo_runs into domain layer table test_cases",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ConvertRuns(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*TestmoTaskData)
	db := taskCtx.GetDal()

	cursor, err := db.Cursor(dal.From(&testmoModels.TestmoRun{}), dal.Where("connection_id = ? AND project_id = ?", data.Options.ConnectionId, data.Options.ProjectId))
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TestmoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_RUN_TABLE,
		},
		InputRowType: reflect.TypeOf(testmoModels.TestmoRun{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			run := inputRow.(*testmoModels.TestmoRun)

			// Convert to domain layer QA test case (treating runs as test cases)
			qaTestCase := &qa.QaTestCase{
				DomainEntityExtended: domainlayer.DomainEntityExtended{
					Id: didgen.NewDomainIdGenerator(&testmoModels.TestmoRun{}).Generate(data.Options.ConnectionId, run.Id),
				},
				Name:        run.Name,
				Type:        getRunType(run),
				QaProjectId: didgen.NewDomainIdGenerator(&testmoModels.TestmoProject{}).Generate(data.Options.ConnectionId, run.ProjectId),
			}

			if run.TestmoCreatedAt != nil && !run.TestmoCreatedAt.IsZero() {
				qaTestCase.CreateTime = *run.TestmoCreatedAt
			}

			// Create test case execution
			qaExecution := &qa.QaTestCaseExecution{
				DomainEntityExtended: domainlayer.DomainEntityExtended{
					Id: didgen.NewDomainIdGenerator(&testmoModels.TestmoRun{}).Generate(data.Options.ConnectionId, run.Id) + ":execution",
				},
				QaProjectId:  didgen.NewDomainIdGenerator(&testmoModels.TestmoProject{}).Generate(data.Options.ConnectionId, run.ProjectId),
				QaTestCaseId: didgen.NewDomainIdGenerator(&testmoModels.TestmoRun{}).Generate(data.Options.ConnectionId, run.Id),
				Status:       convertRunStatus(run.Status),
			}

			if run.TestmoCreatedAt != nil && !run.TestmoCreatedAt.IsZero() {
				qaExecution.CreateTime = *run.TestmoCreatedAt
				qaExecution.StartTime = *run.TestmoCreatedAt
			}
			if run.TestmoUpdatedAt != nil && !run.TestmoUpdatedAt.IsZero() {
				qaExecution.FinishTime = *run.TestmoUpdatedAt
			}

			return []interface{}{qaTestCase, qaExecution}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}

func convertRunStatus(status int32) string {
	// Testmo run status mapping
	switch status {
	case 1:
		return "SUCCESS"
	case 2:
		return "FAILED"
	case 3:
		return "PENDING"
	default:
		return "PENDING"
	}
}

func getRunType(run *testmoModels.TestmoRun) string {
	if run.IsAcceptanceTest {
		return "functional"
	}
	if run.IsSmokeTest {
		return "smoke"
	}
	return "functional"
}
