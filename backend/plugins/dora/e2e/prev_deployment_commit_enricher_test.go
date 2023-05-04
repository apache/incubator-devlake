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

	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/dora/impl"
	"github.com/apache/incubator-devlake/plugins/dora/tasks"
)

func TestPrevSuccessDeploymentCommitEnricherDataFlow(t *testing.T) {
	var plugin impl.Dora
	dataflowTester := e2ehelper.NewDataFlowTester(t, "dora", plugin)

	taskData := &tasks.DoraTaskData{
		Options: &tasks.DoraOptions{
			ProjectName: "project1",
		},
	}
	// import raw data table
	dataflowTester.ImportCsvIntoTabler("./prev_success_deployment_commit/project_mapping.csv", &crossdomain.ProjectMapping{})
	dataflowTester.ImportCsvIntoTabler("./prev_success_deployment_commit/cicd_deployment_commits_before.csv", &devops.CicdDeploymentCommit{})

	// verify converter
	dataflowTester.Subtask(tasks.EnrichPrevSuccessDeploymentCommitMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CicdDeploymentCommit{}, e2ehelper.TableOptions{
		CSVRelPath: "./prev_success_deployment_commit/cicd_deployment_commits_after.csv",
		TargetFields: []string{
			"id",
			"result",
			"started_date",
			"finished_date",
			"cicd_deployment_id",
			"cicd_scope_id",
			"repo_url",
			"environment",
			"prev_success_deployment_commit_id",
		},
	})
}
