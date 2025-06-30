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

var ConvertTestsMeta = plugin.SubTaskMeta{
	Name:             "convertTests",
	EntryPoint:       ConvertTests,
	EnabledByDefault: true,
	Description:      "Convert tool layer table testmo_tests into domain layer table test_cases",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ConvertTests(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*TestmoTaskData)
	db := taskCtx.GetDal()

	cursor, err := db.Cursor(dal.From(&testmoModels.TestmoTest{}), dal.Where("connection_id = ? AND project_id = ?", data.Options.ConnectionId, data.Options.ProjectId))
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
			Table: RAW_TEST_TABLE,
		},
		InputRowType: reflect.TypeOf(testmoModels.TestmoTest{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			test := inputRow.(*testmoModels.TestmoTest)

			// Convert to domain layer QA test case
			qaTestCase := &qa.QaTestCase{
				DomainEntityExtended: domainlayer.DomainEntityExtended{
					Id: didgen.NewDomainIdGenerator(&testmoModels.TestmoTest{}).Generate(data.Options.ConnectionId, test.Id),
				},
				Name:        test.Name,
				Type:        getTestType(test),
				QaProjectId: didgen.NewDomainIdGenerator(&testmoModels.TestmoAutomationRun{}).Generate(data.Options.ConnectionId, test.AutomationRunId),
			}

			if test.TestmoCreatedAt != nil {
				qaTestCase.CreateTime = *test.TestmoCreatedAt
			}

			// Create test case execution
			qaExecution := &qa.QaTestCaseExecution{
				DomainEntityExtended: domainlayer.DomainEntityExtended{
					Id: didgen.NewDomainIdGenerator(&testmoModels.TestmoTest{}).Generate(data.Options.ConnectionId, test.Id) + ":execution",
				},
				QaProjectId:  didgen.NewDomainIdGenerator(&testmoModels.TestmoAutomationRun{}).Generate(data.Options.ConnectionId, test.AutomationRunId),
				QaTestCaseId: didgen.NewDomainIdGenerator(&testmoModels.TestmoTest{}).Generate(data.Options.ConnectionId, test.Id),
				Status:       convertTestStatus(test.Status),
			}

			if test.TestmoCreatedAt != nil {
				qaExecution.CreateTime = *test.TestmoCreatedAt
				qaExecution.StartTime = *test.TestmoCreatedAt
			}
			if test.TestmoUpdatedAt != nil {
				qaExecution.FinishTime = *test.TestmoUpdatedAt
			}

			return []interface{}{qaTestCase, qaExecution}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}

func convertTestStatus(status int32) string {
	// Testmo status mapping - this may need adjustment based on actual Testmo status values
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

func getTestType(test *testmoModels.TestmoTest) string {
	if test.IsAcceptanceTest {
		return "functional"
	}
	if test.IsSmokeTest {
		return "functional"
	}
	return "functional"
}
