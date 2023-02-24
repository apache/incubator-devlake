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

package e2e

import (
	"testing"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/codequality"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/sonarqube/impl"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"github.com/apache/incubator-devlake/plugins/sonarqube/tasks"
)

func TestSonarqubeIssueDataFlow(t *testing.T) {

	var sonarqube impl.Sonarqube
	dataflowTester := e2ehelper.NewDataFlowTester(t, "sonarqube", sonarqube)

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_sonarqube_api_issues.csv",
		"_raw_sonarqube_api_issues")

	// Standard data
	taskData := &tasks.SonarqubeTaskData{
		Options: &tasks.SonarqubeOptions{
			ConnectionId: 1,
			ProjectKey:   "f5a50c63-2e8f-4107-9014-853f6f467757",
		},
	}
	// Interfered data
	taskData2 := &tasks.SonarqubeTaskData{
		Options: &tasks.SonarqubeOptions{
			ConnectionId: 2,
			ProjectKey:   "testWarrenEtcd",
		},
	}

	// verify extraction
	dataflowTester.FlushTabler(&models.SonarqubeIssue{})
	dataflowTester.FlushTabler(&models.SonarqubeIssueCodeBlock{})
	dataflowTester.Subtask(tasks.ExtractIssuesMeta, taskData)

	dataflowTester.Subtask(tasks.ExtractIssuesMeta, taskData2)
	dataflowTester.VerifyTableWithOptions(&models.SonarqubeIssue{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_sonarqube_issues.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
	dataflowTester.VerifyTableWithOptions(&models.SonarqubeIssueCodeBlock{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_sonarqube_issue_code_blocks.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// verify convertor
	dataflowTester.FlushTabler(&codequality.CqIssue{})
	dataflowTester.FlushTabler(&codequality.CqIssueCodeBlock{})
	dataflowTester.Subtask(tasks.ConvertIssuesMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&codequality.CqIssue{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/issues.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
	dataflowTester.Subtask(tasks.ConvertIssueCodeBlocksMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&codequality.CqIssueCodeBlock{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/issue_code_blocks.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
