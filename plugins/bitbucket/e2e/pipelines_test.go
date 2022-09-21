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

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/bitbucket/impl"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
)

func TestBitbucketPipelineDataFlow(t *testing.T) {

	var bitbucket impl.Bitbucket
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bitbucket", bitbucket)

	bitbucketRepository := &models.BitbucketRepo{
		BitbucketId: "thenicetgp/ptest",
	}
	taskData := &tasks.BitbucketTaskData{
		Options: &tasks.BitbucketOptions{
			ConnectionId: 1,
			Owner:        "thenicetgp",
			Repo:         "ptest",
		},
		Repo: bitbucketRepository,
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
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify conversion
	dataflowTester.FlushTabler(&devops.CICDPipeline{})
	dataflowTester.Subtask(tasks.ConvertPipelineMeta, taskData)
	dataflowTester.VerifyTable(
		devops.CICDPipeline{},
		"./snapshot_tables/cicd_pipelines.csv",
		[]string{
			"id",
			"name",
			"result",
			"status",
			"type",
			"duration_sec",
			"created_date",
			"finished_date",
			"environment",
		},
	)
}
