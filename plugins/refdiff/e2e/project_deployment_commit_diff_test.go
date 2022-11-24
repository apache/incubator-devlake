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

	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/refdiff/impl"
	"github.com/apache/incubator-devlake/plugins/refdiff/tasks"
)

func TestRepoDataFlow(t *testing.T) {

	var plugin impl.RefDiff
	dataflowTester := e2ehelper.NewDataFlowTester(t, "refdiff", plugin)

	taskData := &tasks.RefdiffTaskData{
		Options: &tasks.RefdiffOptions{
			ProjectName: "project1",
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoTabler("./raw_tables/project_mapping.csv", &crossdomain.ProjectMapping{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/cicd_tasks.csv", &devops.CICDTask{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/cicd_pipelines.csv", &devops.CICDPipeline{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/cicd_pipeline_commits.csv", &devops.CiCDPipelineCommit{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/repo_commits.csv", &code.RepoCommit{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/commit_parents.csv", &code.CommitParent{})

	// verify extraction
	dataflowTester.FlushTabler(&code.CommitsDiff{})
	dataflowTester.FlushTabler(&code.FinishedCommitsDiff{})

	dataflowTester.Subtask(tasks.CalculateProjectDeploymentCommitsDiffMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&code.CommitsDiff{}, e2ehelper.TableOptions{
		CSVRelPath: "./snapshot_tables/commits_diffs.csv",
	})

	dataflowTester.VerifyTableWithOptions(&code.FinishedCommitsDiff{}, e2ehelper.TableOptions{
		CSVRelPath: "./snapshot_tables/finished_commits_diffs.csv",
	})
}
