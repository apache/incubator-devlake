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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/impl"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
)

func TestGitlabPipelineDetailDataFlow(t *testing.T) {

	var gitlab impl.Gitlab
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", gitlab)

	regexEnricher := api.NewRegexEnricher()
	_ = regexEnricher.TryAdd(devops.DEPLOYMENT, "EE-7121")
	_ = regexEnricher.TryAdd(devops.PRODUCTION, "EE-7121")

	taskData := &tasks.GitlabTaskData{
		Options: &tasks.GitlabOptions{
			ConnectionId: 1,
			ProjectId:    12345678,
			ScopeConfig:  new(models.GitlabScopeConfig),
		},
		RegexEnricher: regexEnricher,
	}

	// import raw data table
	dataflowTester.FlushTabler(&models.GitlabPipelineProject{})
	dataflowTester.FlushTabler(models.GitlabPipeline{})
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_gitlab_api_pipeline_details.csv", "_raw_gitlab_api_pipeline_details")
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_gitlab_pipeline_projects.csv", &models.GitlabPipelineProject{})

	// verify extraction
	dataflowTester.Subtask(tasks.ExtractApiPipelineDetailsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GitlabPipeline{},
		"./snapshot_tables/_tool_gitlab_pipelines.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"gitlab_id",
			"gitlab_created_at",
			"project_id",
			"status",
			"ref",
			"sha",
			"web_url",
			"duration",
			"queued_duration",
			"started_at",
			"finished_at",
			"coverage",
			"type",
			"environment",
		),
	)

	// verify conversion
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_gitlab_projects.csv", &models.GitlabProject{})
	dataflowTester.FlushTabler(&devops.CICDPipeline{})
	dataflowTester.FlushTabler(&devops.CiCDPipelineCommit{})
	dataflowTester.Subtask(tasks.ConvertDetailPipelineMeta, taskData)
	dataflowTester.Subtask(tasks.ConvertPipelineCommitMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CICDPipeline{}, e2ehelper.TableOptions{
		CSVRelPath:   "./snapshot_tables/cicd_pipelines.csv",
		IgnoreTypes:  []interface{}{common.NoPKModel{}},
		IgnoreFields: []string{"is_child"},
	})

	dataflowTester.VerifyTableWithOptions(&devops.CiCDPipelineCommit{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_pipeline_commits.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
