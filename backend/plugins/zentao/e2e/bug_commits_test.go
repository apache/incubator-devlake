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

func TestZentaoBugCommitsDataFlow(t *testing.T) {

	var zentao impl.Zentao
	dataflowTester := e2ehelper.NewDataFlowTester(t, "zentao", zentao)

	taskData := &tasks.ZentaoTaskData{
		Options: &tasks.ZentaoOptions{
			ConnectionId: 1,
			ProjectId:    22,
		},
	}

	// import _raw_zentao_api_bug_commits raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_zentao_api_bug_commits.csv",
		"_raw_zentao_api_bug_commits")

	// verify bug commit extraction
	dataflowTester.FlushTabler(&models.ZentaoBugCommit{})
	dataflowTester.Subtask(tasks.ExtractBugCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.ZentaoBugCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_zentao_bug_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// import _raw_zentao_api_bug_repo_commits raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_zentao_api_bug_repo_commits.csv",
		"_raw_zentao_api_bug_repo_commits")

	// verify bug repo commit extraction
	dataflowTester.FlushTabler(&models.ZentaoBugRepoCommit{})
	dataflowTester.Subtask(tasks.ExtractBugRepoCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.ZentaoBugRepoCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_zentao_bug_repo_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// verify bug repo commit conversion
	dataflowTester.FlushTabler(&crossdomain.IssueRepoCommit{})
	dataflowTester.Subtask(tasks.ConvertBugRepoCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&crossdomain.IssueRepoCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/issue_bug_repo_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

}
