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
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/pagerduty/impl"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"github.com/apache/incubator-devlake/plugins/pagerduty/tasks"
	"testing"
)

func TestIncidentDataFlow(t *testing.T) {
	var plugin impl.PagerDuty
	dataflowTester := e2ehelper.NewDataFlowTester(t, "pagerduty", plugin)

	taskData := &tasks.PagerDutyTaskData{
		Options: &tasks.PagerDutyOptions{
			ConnectionId: 1,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_pagerduty_incidents.csv", "_raw_pagerduty_incidents")

	// verify worklog extraction
	dataflowTester.FlushTabler(&models.Incident{})
	dataflowTester.FlushTabler(&models.User{})
	dataflowTester.FlushTabler(&models.Service{})
	dataflowTester.FlushTabler(&models.Assignment{})
	dataflowTester.Subtask(tasks.ExtractIncidentsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(
		models.Incident{},
		e2ehelper.TableOptions{
			CSVRelPath:  "./snapshot_tables/_tool_pagerduty_incident.csv",
			IgnoreTypes: []any{common.Model{}},
		},
	)
	dataflowTester.VerifyTableWithOptions(
		models.User{},
		e2ehelper.TableOptions{
			CSVRelPath:  "./snapshot_tables/_tool_pagerduty_user.csv",
			IgnoreTypes: []any{common.Model{}},
		},
	)
	dataflowTester.VerifyTableWithOptions(
		models.Assignment{},
		e2ehelper.TableOptions{
			CSVRelPath:  "./snapshot_tables/_tool_pagerduty_assignment.csv",
			IgnoreTypes: []any{common.Model{}},
		},
	)
	dataflowTester.VerifyTableWithOptions(
		models.Service{},
		e2ehelper.TableOptions{
			CSVRelPath:  "./snapshot_tables/_tool_pagerduty_service.csv",
			IgnoreTypes: []any{common.Model{}},
		},
	)
	dataflowTester.FlushTabler(&ticket.Issue{})
	dataflowTester.Subtask(tasks.ConvertIncidentsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(
		ticket.Issue{},
		e2ehelper.TableOptions{
			CSVRelPath:  "./snapshot_tables/issues.csv",
			IgnoreTypes: []any{common.NoPKModel{}},
		},
	)
}
