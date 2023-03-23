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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/teambition/impl"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
	"github.com/apache/incubator-devlake/plugins/teambition/tasks"
	"testing"
)

func TestTeambitionProject(t *testing.T) {

	var teambition impl.Teambition
	dataflowTester := e2ehelper.NewDataFlowTester(t, "teambition", teambition)

	taskData := &tasks.TeambitionTaskData{
		Options: &tasks.TeambitionOptions{
			ConnectionId: 1,
			ProjectId:    "64132c94f0d59df1c9825ab8",
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_teambition_api_projects.csv",
		"_raw_teambition_api_projects")
	dataflowTester.FlushTabler(&models.TeambitionProject{})

	// verify extraction
	dataflowTester.Subtask(tasks.ExtractProjectsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(
		models.TeambitionProject{},
		e2ehelper.TableOptions{
			CSVRelPath:   "./snapshot_tables/_tool_teambition_projects.csv",
			IgnoreTypes:  []interface{}{common.NoPKModel{}},
			IgnoreFields: []string{"created", "updated", "start_date", "end_date"},
		},
	)

	dataflowTester.FlushTabler(&ticket.Board{})
	dataflowTester.Subtask(tasks.ConvertProjectsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(
		ticket.Board{},
		e2ehelper.TableOptions{
			CSVRelPath:   "./snapshot_tables/boards.csv",
			IgnoreFields: []string{"created_date"},
			IgnoreTypes:  []interface{}{domainlayer.DomainEntity{}},
		},
	)
}
