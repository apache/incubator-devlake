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

func TestZentaoStoryCommitsDataFlow(t *testing.T) {

	var zentao impl.Zentao
	dataflowTester := e2ehelper.NewDataFlowTester(t, "zentao", zentao)

	taskData := &tasks.ZentaoTaskData{
		Options: &tasks.ZentaoOptions{
			ConnectionId: 1,
			ProjectId:    1,
		},
		ApiClient: getFakeAPIClient(),
	}

	// import _raw_zentao_api_story_commits raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_zentao_api_story_commits.csv",
		"_raw_zentao_api_story_commits")

	// verify story commit extraction
	dataflowTester.FlushTabler(&models.ZentaoStoryCommit{})
	dataflowTester.Subtask(tasks.ExtractStoryCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.ZentaoStoryCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_zentao_story_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// import _raw_zentao_api_story_repo_commits raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_zentao_api_story_repo_commits.csv",
		"_raw_zentao_api_story_repo_commits")

	// verify story repo commit extraction
	dataflowTester.FlushTabler(&models.ZentaoStoryRepoCommit{})
	dataflowTester.Subtask(tasks.ExtractStoryRepoCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.ZentaoStoryRepoCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_zentao_story_repo_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// verify story repo commit conversion
	dataflowTester.FlushTabler(&crossdomain.IssueRepoCommit{})
	dataflowTester.Subtask(tasks.ConvertStoryRepoCommitsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&crossdomain.IssueRepoCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/issue_story_repo_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

}
