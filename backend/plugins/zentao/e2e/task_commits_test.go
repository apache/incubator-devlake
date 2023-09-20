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
	"github.com/apache/incubator-devlake/plugins/zentao/impl"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/apache/incubator-devlake/plugins/zentao/tasks"
)

func TestZentaoTaskCommitsDataFlow(t *testing.T) {

	var zentao impl.Zentao
	dataflowTester := e2ehelper.NewDataFlowTester(t, "zentao", zentao)

	taskData := &tasks.ZentaoTaskData{
		Options: &tasks.ZentaoOptions{
			ConnectionId: 1,
			ProjectId:    48,
		},
		ApiClient: getFakeAPIClient(),
	}

	// import _raw_zentao_api_task_commits raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_zentao_api_task_commits.csv",
		"_raw_zentao_api_task_commits")

	// verify task commit extraction
	dataflowTester.FlushTabler(&models.ZentaoTaskCommit{})
	dataflowTester.Subtask(tasks.ExtractTaskCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.ZentaoTaskCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_zentao_task_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// import _raw_zentao_api_task_repo_commits raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_zentao_api_task_repo_commits.csv",
		"_raw_zentao_api_task_repo_commits")

	// verify task repo commit extraction
	dataflowTester.FlushTabler(&models.ZentaoTaskRepoCommit{})
	dataflowTester.Subtask(tasks.ExtractTaskRepoCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.ZentaoTaskRepoCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_zentao_task_repo_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// verify task repo commit conversion
	dataflowTester.FlushTabler(&crossdomain.IssueRepoCommit{})
	dataflowTester.Subtask(tasks.ConvertTaskRepoCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&crossdomain.IssueRepoCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/issue_task_repo_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

}
