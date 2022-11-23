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
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/dora/impl"
	"github.com/apache/incubator-devlake/plugins/dora/tasks"
	"testing"
)

func TestCalculateCLTimeDataFlow(t *testing.T) {
	var plugin impl.Dora
	dataflowTester := e2ehelper.NewDataFlowTester(t, "dora", plugin)

	taskData := &tasks.DoraTaskData{
		Options: &tasks.DoraOptions{
			ProjectName: "project1",
			TransformationRules: tasks.TransformationRules{
				ProductionPattern: "(?i)deploy",
			},
		},
	}

	dataflowTester.FlushTabler(&code.PullRequest{})

	// import raw data table
	dataflowTester.ImportCsvIntoTabler("./raw_tables/cicd_tasks_changeleadtime.csv", &devops.CICDTask{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/pull_requests.csv", &code.PullRequest{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/commits_diffs.csv", &code.CommitsDiff{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/cicd_pipeline_commits.csv", &devops.CiCDPipelineCommit{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/project_mapping.csv", &crossdomain.ProjectMapping{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/commits.csv", &code.Commit{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/pull_request_comments.csv", &code.PullRequestComment{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/pull_request_commits.csv", &code.PullRequestCommit{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/repos.csv", &code.Repo{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/cicd_scopes.csv", &devops.CicdScope{})

	// verify converter
	dataflowTester.FlushTabler(&crossdomain.ProjectPrMetric{})
	dataflowTester.Subtask(tasks.CalculateChangeLeadTimeMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&crossdomain.ProjectPrMetric{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/project_pr_metrics.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
