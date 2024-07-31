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
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/circleci/impl"
	"github.com/apache/incubator-devlake/plugins/circleci/models"
	"github.com/apache/incubator-devlake/plugins/circleci/tasks"
)

func TestCircleciPipeline(t *testing.T) {
	var circleci impl.Circleci

	dataflowTester := e2ehelper.NewDataFlowTester(t, "circleci", circleci)
	taskData := &tasks.CircleciTaskData{
		Options: &tasks.CircleciOptions{
			ConnectionId: 1,
			ProjectSlug:  "github/coldgust/coldgust.github.io",
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_circleci_api_pipelines.csv",
		"_raw_circleci_api_pipelines")
	dataflowTester.FlushTabler(&models.CircleciPipeline{})

	// verify extraction
	dataflowTester.Subtask(tasks.ExtractPipelinesMeta, taskData)
	dataflowTester.VerifyTableWithOptions(
		models.CircleciPipeline{},
		e2ehelper.TableOptions{
			CSVRelPath:   "./snapshot_tables/_tool_circleci_pipelines.csv",
			IgnoreTypes:  []interface{}{common.NoPKModel{}},
			IgnoreFields: []string{"stopped_date"},
		},
	)
}
