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

	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/bitbucket/impl"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
)

func TestRepoDataFlow(t *testing.T) {
	var plugin impl.Bitbucket
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bitbucket", plugin)

	bitbucketRepository := &models.BitbucketRepo{
		BitbucketId: "panjf2000/ants",
	}
	taskData := &tasks.BitbucketTaskData{
		Options: &tasks.BitbucketOptions{
			ConnectionId: 1,
			Owner:        "panjf2000",
			Repo:         "ants",
			TransformationRules: models.TransformationRules{
				PrType:      "type/(.*)$",
				PrComponent: "component/(.*)$",
			},
		},
		Repo: bitbucketRepository,
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_bitbucket_api_repositories.csv", "_raw_bitbucket_api_repositories")

	// verify extraction
	dataflowTester.FlushTabler(&models.BitbucketRepo{})
	dataflowTester.FlushTabler(&models.BitbucketAccount{})
	dataflowTester.Subtask(tasks.ExtractApiRepoMeta, taskData)
	dataflowTester.VerifyTable(
		models.BitbucketRepo{},
		"./snapshot_tables/_tool_bitbucket_repos.csv",
		[]string{
			"connection_id",
			"bitbucket_id",
			"name",
			"html_url",
			"description",
			"owner_id",
			"language",
			"created_date",
			"updated_date",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.BitbucketAccount{},
		"./snapshot_tables/_tool_bitbucket_accounts.csv",
		[]string{
			"connection_id",
			"user_name",
			"account_id",
			"account_status",
			"display_name",
			"avatar_url",
			"html_url",
			"uuid",
			"has2_fa_enabled",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify extraction
	dataflowTester.FlushTabler(&code.Repo{})
	dataflowTester.FlushTabler(&ticket.Board{})
	dataflowTester.Subtask(tasks.ConvertRepoMeta, taskData)
	dataflowTester.VerifyTable(
		code.Repo{},
		"./snapshot_tables/repos.csv",
		[]string{
			"id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"name",
			"url",
			"description",
			"owner_id",
			"language",
			"forked_from",
			"created_date",
			"updated_date",
			"deleted",
		},
	)
	dataflowTester.VerifyTable(
		ticket.Board{},
		"./snapshot_tables/boards.csv",
		[]string{
			"id",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"name",
			"description",
			"url",
			"created_date",
		},
	)
}
