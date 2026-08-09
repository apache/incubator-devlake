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
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/monorepo/impl"
	"github.com/apache/incubator-devlake/plugins/monorepo/models"
	"github.com/apache/incubator-devlake/plugins/monorepo/tasks"
	"github.com/stretchr/testify/assert"
)

// TestMonorepoAttributionDataFlow exercises both subtasks against a monorepo containing
// serviceA and serviceB, each with its own deploy job.
//
// The fixtures deliberately include the cases that motivated this plugin:
//   - pr2 (serviceB) merges at 09:00 while serviceA deploys at 10:00 and serviceB only at
//     12:00. DORA would link pr2 to the 10:00 deployment because it searches the whole
//     repository; pr2 must instead link to 12:00.
//   - pipeline3 runs both deploy jobs, so it must yield one row per sub-project.
//   - a failed deployment and a staging deployment sit between pr1's merge and the
//     deployment that actually shipped it, so neither may be linked.
func TestMonorepoAttributionDataFlow(t *testing.T) {
	var plugin impl.Monorepo
	dataflowTester := e2ehelper.NewDataFlowTester(t, "monorepo", plugin)

	subProjects := []tasks.SubProjectConfig{
		{
			Name:             "serviceA",
			PrLabels:         []string{"serviceA"},
			DeployJobPattern: "^deploy-serviceA$",
		},
		{
			Name:             "serviceB",
			PrLabels:         []string{"serviceB"},
			DeployJobPattern: "^deploy-serviceB$",
		},
	}
	matcher, err := tasks.NewSubProjectMatcher(subProjects)
	assert.Nil(t, err)

	taskData := &tasks.MonorepoTaskData{
		Options: &tasks.MonorepoOptions{
			ProjectName: "monorepo",
			SubProjects: subProjects,
		},
		Matcher: matcher,
	}

	// seed the domain layer
	dataflowTester.FlushTabler(&crossdomain.ProjectMapping{})
	dataflowTester.FlushTabler(&devops.CICDTask{})
	dataflowTester.FlushTabler(&devops.CicdDeploymentCommit{})
	dataflowTester.FlushTabler(&code.PullRequest{})
	dataflowTester.FlushTabler(&code.PullRequestLabel{})
	dataflowTester.FlushTabler(&crossdomain.ProjectPrMetric{})

	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/project_mapping.csv", &crossdomain.ProjectMapping{})
	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/cicd_tasks.csv", &devops.CICDTask{})
	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/cicd_deployment_commits.csv", &devops.CicdDeploymentCommit{})
	dataflowTester.ImportNullableCsvIntoTabler("./monorepo_attribution/pull_requests.csv", &code.PullRequest{})
	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/pull_request_labels.csv", &code.PullRequestLabel{})
	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/project_pr_metrics.csv", &crossdomain.ProjectPrMetric{})

	// deployments must be attributed first: the pull request subtask reads them back to
	// work out which deployment shipped each merged pull request.
	dataflowTester.FlushTabler(&models.SubProjectDeployment{})
	dataflowTester.Subtask(tasks.AttributeDeploymentsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.SubProjectDeployment{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/monorepo_subproject_deployments.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.FlushTabler(&models.SubProjectPrMetric{})
	dataflowTester.Subtask(tasks.AttributePullRequestsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.SubProjectPrMetric{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/monorepo_subproject_pr_metrics.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
