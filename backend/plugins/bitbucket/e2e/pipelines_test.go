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

	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/impl"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
)

func TestBitbucketPipelineDataFlow(t *testing.T) {

	var bitbucket impl.Bitbucket
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bitbucket", bitbucket)

	regexEnricher := helper.NewRegexEnricher()
	_ = regexEnricher.TryAdd(devops.DEPLOYMENT, "main")
	_ = regexEnricher.TryAdd(devops.PRODUCTION, "pipeline")
	taskData := &tasks.BitbucketTaskData{
		Options: &tasks.BitbucketOptions{
			ConnectionId: 1,
			FullName:     "likyh/likyhphp",
		},
		RegexEnricher: regexEnricher,
	}
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_bitbucket_api_pipelines.csv", "_raw_bitbucket_api_pipelines")

	// verify extraction
	dataflowTester.FlushTabler(&models.BitbucketPipeline{})
	dataflowTester.Subtask(tasks.ExtractApiPipelinesMeta, taskData)
	dataflowTester.VerifyTable(
		models.BitbucketPipeline{},
		"./snapshot_tables/_tool_bitbucket_pipelines.csv",
		[]string{
			"connection_id",
			"bitbucket_id",
			"status",
			"result",
			"ref_name",
			"web_url",
			"duration_in_seconds",
			"type",
			"repo_id",
			"environment",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify conversion
	dataflowTester.FlushTabler(&devops.CiCDPipelineCommit{})
	dataflowTester.FlushTabler(&devops.CICDPipeline{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_bitbucket_repos.csv", &models.BitbucketRepo{})
	dataflowTester.Subtask(tasks.ConvertPipelineMeta, taskData)
	dataflowTester.VerifyTable(
		devops.CICDPipeline{},
		"./snapshot_tables/cicd_pipelines.csv",
		[]string{
			"id",
			"name",
			"result",
			"status",
			"original_result",
			"original_status",
			"type",
			"duration_sec",
			"environment",
			"display_title",
			"url",
		},
	)

	dataflowTester.VerifyTable(
		devops.CiCDPipelineCommit{},
		"./snapshot_tables/cicd_pipeline_commits.csv",
		[]string{
			"pipeline_id",
			"commit_sha",
			"branch",
			"repo_id",
			"repo_url",
			"display_title",
			"url",
		},
	)
}
