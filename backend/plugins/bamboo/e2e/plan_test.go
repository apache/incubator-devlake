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
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/bamboo/impl"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"github.com/apache/incubator-devlake/plugins/bamboo/tasks"
)

func TestBambooPlanDataFlow(t *testing.T) {

	var bamboo impl.Bamboo
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bamboo", bamboo)

	taskData := &tasks.BambooOptions{
		Options: &models.BambooOptions{
			ConnectionId:      3,
			PlanKey:           "TEST1",
			BambooScopeConfig: new(models.BambooScopeConfig),
		},
		ApiClient: getFakeAPIClient(),
	}

	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_bamboo_plans.csv", models.BambooPlan{})
	// verify extraction
	dataflowTester.FlushTabler(&devops.CicdScope{})
	dataflowTester.Subtask(tasks.ConvertPlansMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CicdScope{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_scopes.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
