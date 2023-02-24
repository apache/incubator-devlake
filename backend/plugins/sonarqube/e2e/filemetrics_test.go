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

func TestSonarqubeFileMetricsDataFlow(t *testing.T) {

	var sonarqube impl.Sonarqube
	dataflowTester := e2ehelper.NewDataFlowTester(t, "sonarqube", sonarqube)

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_sonarqube_api_filemetrics.csv",
		"_raw_sonarqube_api_filemetrics")
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_sonarqube_api_filemetrics_additional.csv",
		"_raw_sonarqube_api_filemetrics_additional")

	// Standard data
	taskData := &tasks.SonarqubeTaskData{
		Options: &tasks.SonarqubeOptions{
			ConnectionId: 2,
			ProjectKey:   "testDevLake",
		},
	}
	// Interfered data
	taskData2 := &tasks.SonarqubeTaskData{
		Options: &tasks.SonarqubeOptions{
			ConnectionId: 1,
			ProjectKey:   "testNone",
		},
	}

	// verify extraction
	dataflowTester.FlushTabler(&models.SonarqubeWholeFileMetrics{})
	dataflowTester.Subtask(tasks.ExtractFilemetricsMeta, taskData)
	dataflowTester.Subtask(tasks.ExtractAdditionalFileMetricsMeta, taskData)

	dataflowTester.Subtask(tasks.ExtractFilemetricsMeta, taskData2)
	dataflowTester.Subtask(tasks.ExtractAdditionalFileMetricsMeta, taskData2)
	dataflowTester.VerifyTableWithOptions(&models.SonarqubeWholeFileMetrics{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_sonarqube_filemetrics.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// verify convertor
	dataflowTester.FlushTabler(&codequality.CqFileMetrics{})
	dataflowTester.Subtask(tasks.ConvertFileMetricsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&codequality.CqFileMetrics{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/filemetrics.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
