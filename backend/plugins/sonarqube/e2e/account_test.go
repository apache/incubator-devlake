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
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/sonarqube/impl"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"github.com/apache/incubator-devlake/plugins/sonarqube/tasks"
)

func TestSonarqubeAccountDataFlow(t *testing.T) {

	var sonarqube impl.Sonarqube
	dataflowTester := e2ehelper.NewDataFlowTester(t, "sonarqube", sonarqube)

	taskData := &tasks.SonarqubeTaskData{
		Options: &tasks.SonarqubeOptions{
			ConnectionId: 1,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_sonarqube_api_accounts.csv",
		"_raw_sonarqube_api_accounts")

	// verify extraction
	dataflowTester.FlushTabler(&models.SonarqubeAccount{})
	dataflowTester.Subtask(tasks.ExtractAccountsMeta, taskData)

	taskData2 := &tasks.SonarqubeTaskData{
		Options: &tasks.SonarqubeOptions{
			ConnectionId: 2,
		},
	}

	dataflowTester.Subtask(tasks.ExtractAccountsMeta, taskData2)
	dataflowTester.VerifyTableWithOptions(&models.SonarqubeAccount{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_sonarqube_accounts.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.FlushTabler(&crossdomain.Account{})
	dataflowTester.Subtask(tasks.ConvertAccountsMeta, taskData2)
	dataflowTester.VerifyTableWithOptions(&crossdomain.Account{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/accounts.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
