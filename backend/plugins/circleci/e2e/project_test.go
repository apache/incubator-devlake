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
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/circleci/impl"
	"github.com/apache/incubator-devlake/plugins/circleci/models"
	"github.com/apache/incubator-devlake/plugins/circleci/tasks"
)

func TestCircleciProject(t *testing.T) {
	var circleci impl.Circleci

	dataflowTester := e2ehelper.NewDataFlowTester(t, "circleci", circleci)
	taskData := &tasks.CircleciTaskData{
		Options: &tasks.CircleciOptions{
			ConnectionId: 1,
			ProjectSlug:  "github/coldgust/coldgust.github.io",
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_circleci_api_projects.csv",
		"_raw_circleci_api_projects")
	dataflowTester.FlushTabler(&models.CircleciProject{})

	// verify extraction
	dataflowTester.Subtask(tasks.ExtractProjectsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(
		models.CircleciProject{},
		e2ehelper.TableOptions{
			CSVRelPath:   "./snapshot_tables/_tool_circleci_projects.csv",
			IgnoreTypes:  []interface{}{common.NoPKModel{}},
			IgnoreFields: []string{"started_at", "stopped_at"},
		},
	)

	dataflowTester.FlushTabler(&devops.CicdScope{})
	dataflowTester.Subtask(tasks.ConvertProjectsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(
		devops.CicdScope{},
		e2ehelper.TableOptions{
			CSVRelPath:   "./snapshot_tables/cicd_scopes.csv",
			IgnoreFields: []string{"created_date"},
			IgnoreTypes:  []interface{}{domainlayer.DomainEntity{}},
		},
	)
}
