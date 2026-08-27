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
	"github.com/stretchr/testify/require"
)

// TestMonorepoAttributionDataFlow exercises all three subtasks against a monorepo
// containing serviceA and serviceB, each with its own deploy job.
//
// The fixtures deliberately include the cases that motivated this plugin:
//   - pr2 (serviceB) merges at 09:00 while serviceA deploys at 10:00 and serviceB only at
//     12:00. A naive "nearest deployment repo-wide" heuristic would link pr2 to the 10:00
//     deployment; it must instead link to 12:00 (via DORA's project_pr_metrics, which the
//     fixture seeds directly rather than re-deriving from commit ancestry).
//   - pipeline3 runs both deploy jobs, so it must yield one row per sub-project in
//     cicd_deployment_subprojects, and pr3 (shipped by pipeline3) triggers an *expected*
//     label/deployment mismatch against the sub-project pipeline3 also deploys.
//   - a failed deployment and a staging deployment are still attributed by job name alone
//     (attribution never looks at result/environment).
//   - pipeline8 runs a job that matches no configured sub-project, so it - and pr4 - land
//     in the 'unattributed' bucket.
//   - pr6 belongs to a different project ("other") with no monorepo configuration in this
//     run, so it must be left with sub_project = NULL entirely untouched.
//   - pr7 is open (never merged), attributed anyway per Fix 2, but absent from
//     project_pr_metrics (DORA never computed metrics for it) and therefore absent from
//     the monorepo_subproject_pr_metrics backfill too.
//   - pr8 is labelled serviceA but its project_pr_metrics row says DORA shipped it via
//     pipeline2, which only deploys serviceB: a genuine label/deployJobPattern
//     misconfiguration that FindSubProjectMismatches must surface.
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
			// Exercise the default explicitly rather than relying on a nil pointer, so
			// this test keeps working if the default ever changes.
			IncludeUnattributed: boolPtr(true),
		},
		Matcher: matcher,
	}

	// seed the domain layer
	dataflowTester.FlushTabler(&crossdomain.ProjectMapping{})
	dataflowTester.FlushTabler(&devops.CICDTask{})
	dataflowTester.FlushTabler(&devops.CicdDeploymentCommit{})
	dataflowTester.FlushTabler(&code.PullRequest{})
	dataflowTester.FlushTabler(&code.PullRequestLabel{})
	dataflowTester.FlushTabler(&code.PullRequestCommit{})
	dataflowTester.FlushTabler(&crossdomain.ProjectPrMetric{})
	dataflowTester.FlushTabler(&devops.CicdDeploymentSubproject{})

	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/project_mapping.csv", &crossdomain.ProjectMapping{})
	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/cicd_tasks.csv", &devops.CICDTask{})
	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/cicd_deployment_commits.csv", &devops.CicdDeploymentCommit{})
	dataflowTester.ImportNullableCsvIntoTabler("./monorepo_attribution/pull_requests.csv", &code.PullRequest{})
	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/pull_request_labels.csv", &code.PullRequestLabel{})
	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/pull_request_commits.csv", &code.PullRequestCommit{})
	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/project_pr_metrics.csv", &crossdomain.ProjectPrMetric{})

	// 1. attributeDeployments must run first: updateProjectPrMetricsSubProject's
	// cross-check reads cicd_deployment_subprojects back.
	dataflowTester.FlushTabler(&models.SubProjectDeployment{})
	dataflowTester.Subtask(tasks.AttributeDeploymentsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CicdDeploymentSubproject{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_deployment_subprojects.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
	dataflowTester.VerifyTableWithOptions(&models.SubProjectDeployment{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/monorepo_subproject_deployments.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// 2. attributePullRequests: attribution only. Assert it does NOT touch
	// monorepo_subproject_pr_metrics (that responsibility moved entirely to
	// updateProjectPrMetricsSubProject) by flushing the compat table to empty first and
	// confirming it is still empty right after this subtask runs.
	dataflowTester.FlushTabler(&models.SubProjectPrMetric{})
	dataflowTester.Subtask(tasks.AttributePullRequestsMeta, taskData)

	var prMetricsAfterAttribution []models.SubProjectPrMetric
	require.NoError(t, dataflowTester.Dal.All(&prMetricsAfterAttribution))
	assert.Empty(t, prMetricsAfterAttribution,
		"attributePullRequests must not write monorepo_subproject_pr_metrics - that belongs to updateProjectPrMetricsSubProject")

	dataflowTester.VerifyTableWithOptions(&code.PullRequest{}, e2ehelper.TableOptions{
		CSVRelPath:   "./snapshot_tables/pull_requests_sub_project.csv",
		TargetFields: []string{"sub_project"},
	})
	dataflowTester.VerifyTableWithOptions(&code.PullRequestCommit{}, e2ehelper.TableOptions{
		CSVRelPath:   "./snapshot_tables/pull_request_commits_sub_project.csv",
		TargetFields: []string{"sub_project"},
	})

	// 3. updateProjectPrMetricsSubProject: tag project_pr_metrics, cross-check against
	// deployment attribution, and (only now) backfill monorepo_subproject_pr_metrics.
	dataflowTester.Subtask(tasks.UpdateProjectPrMetricsSubProjectMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&crossdomain.ProjectPrMetric{}, e2ehelper.TableOptions{
		CSVRelPath:   "./snapshot_tables/project_pr_metrics_sub_project.csv",
		TargetFields: []string{"sub_project"},
	})
	dataflowTester.VerifyTableWithOptions(&models.SubProjectPrMetric{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/monorepo_subproject_pr_metrics.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// pr3 is shipped by pipeline3, which deploys both serviceA and serviceB; it is
	// labelled serviceA, so it necessarily disagrees with pipeline3's serviceB mapping -
	// that is expected noise from a multi-sub-project deployment, not a misconfiguration.
	// pr8 is a genuine misconfiguration: labelled serviceA but shipped by pipeline2, which
	// only ever deploys serviceB. Labels win either way: both PRs keep sub_project =
	// serviceA in project_pr_metrics (verified above), the mismatch is only logged.
	mismatches, mErr := tasks.FindSubProjectMismatches(dataflowTester.Dal, "monorepo")
	require.NoError(t, mErr)
	require.Len(t, mismatches, 2)
	byPr := map[string]tasks.SubProjectMismatchRow{}
	for _, m := range mismatches {
		byPr[m.PrId] = m
	}
	require.Contains(t, byPr, "pr3")
	assert.Equal(t, "serviceA", byPr["pr3"].PrSubProject)
	assert.Equal(t, "serviceB", byPr["pr3"].DeploymentSubProject)
	require.Contains(t, byPr, "pr8")
	assert.Equal(t, "serviceA", byPr["pr8"].PrSubProject)
	assert.Equal(t, "serviceB", byPr["pr8"].DeploymentSubProject)
}

// TestMonorepoAttributeDeploymentsIncludeUnattributedFalse exercises the
// includeUnattributed=false path (design doc decision 3) against attributeDeployments:
// pipeline8's job matches no configured sub-project, so with the "old" behaviour
// restored, it should be skipped entirely rather than getting an 'unattributed' row in
// either cicd_deployment_subprojects or the compat monorepo_subproject_deployments table.
func TestMonorepoAttributeDeploymentsIncludeUnattributedFalse(t *testing.T) {
	var plugin impl.Monorepo
	dataflowTester := e2ehelper.NewDataFlowTester(t, "monorepo", plugin)

	subProjects := []tasks.SubProjectConfig{
		{Name: "serviceA", DeployJobPattern: "^deploy-serviceA$"},
		{Name: "serviceB", DeployJobPattern: "^deploy-serviceB$"},
	}
	matcher, err := tasks.NewSubProjectMatcher(subProjects)
	require.NoError(t, err)

	taskData := &tasks.MonorepoTaskData{
		Options: &tasks.MonorepoOptions{
			ProjectName:         "monorepo",
			SubProjects:         subProjects,
			IncludeUnattributed: boolPtr(false),
		},
		Matcher: matcher,
	}

	dataflowTester.FlushTabler(&crossdomain.ProjectMapping{})
	dataflowTester.FlushTabler(&devops.CICDTask{})
	dataflowTester.FlushTabler(&devops.CicdDeploymentCommit{})
	dataflowTester.FlushTabler(&devops.CicdDeploymentSubproject{})
	dataflowTester.FlushTabler(&models.SubProjectDeployment{})

	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/project_mapping.csv", &crossdomain.ProjectMapping{})
	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/cicd_tasks.csv", &devops.CICDTask{})
	dataflowTester.ImportCsvIntoTabler("./monorepo_attribution/cicd_deployment_commits.csv", &devops.CicdDeploymentCommit{})

	dataflowTester.Subtask(tasks.AttributeDeploymentsMeta, taskData)

	var mappings []devops.CicdDeploymentSubproject
	require.NoError(t, dataflowTester.Dal.All(&mappings))
	for _, m := range mappings {
		assert.NotEqual(t, tasks.UnattributedSubProject, m.SubProject,
			"pipeline8 must be skipped, not marked unattributed, when includeUnattributed is false")
	}
	// pipeline1/2/3(x2, one per sub-project)/6/7 are still attributed normally; only
	// pipeline8 (which matches nothing) is affected by the flag.
	assert.Len(t, mappings, 6)

	var compat []models.SubProjectDeployment
	require.NoError(t, dataflowTester.Dal.All(&compat))
	assert.Len(t, compat, 6)
}

func boolPtr(b bool) *bool {
	return &b
}
