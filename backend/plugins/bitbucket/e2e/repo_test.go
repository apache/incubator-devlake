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
	"encoding/json"
	"testing"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/helpers/pluginhelper"
	"github.com/apache/incubator-devlake/plugins/bitbucket/impl"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
	"github.com/stretchr/testify/assert"
)

func TestRepoDataFlow(t *testing.T) {
	var plugin impl.Bitbucket
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bitbucket", plugin)

	taskData := &tasks.BitbucketTaskData{
		Options: &tasks.BitbucketOptions{
			ConnectionId: 1,
			FullName:     "likyh/likyhphp",
		},
	}

	// import raw data table
	csvIter, _ := pluginhelper.NewCsvFileIterator("./raw_tables/_raw_bitbucket_api_repositories.csv")
	defer csvIter.Close()
	apiRepo := &models.BitbucketApiRepo{}
	// load rows and insert into target table
	for csvIter.HasNext() {
		toInsertValues := csvIter.Fetch()
		data := json.RawMessage(toInsertValues[`data`].(string))
		err := errors.Convert(json.Unmarshal(data, apiRepo))
		assert.Nil(t, err)
		break
	}

	// verify extraction
	dataflowTester.FlushTabler(&models.BitbucketRepo{})
	scope := apiRepo.ConvertApiScope().(*models.BitbucketRepo)
	scope.ConnectionId = 1
	err := dataflowTester.Dal.CreateIfNotExist(scope)
	assert.Nil(t, err)
	dataflowTester.VerifyTable(
		models.BitbucketRepo{},
		"./snapshot_tables/_tool_bitbucket_repos.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"bitbucket_id",
			"name",
			"html_url",
			"description",
			"owner",
			"language",
			"clone_url",
		),
	)

	// verify extraction
	dataflowTester.FlushTabler(&code.Repo{})
	dataflowTester.FlushTabler(&ticket.Board{})
	dataflowTester.FlushTabler(&devops.CicdScope{})
	dataflowTester.Subtask(tasks.ConvertRepoMeta, taskData)
	dataflowTester.VerifyTable(
		code.Repo{},
		"./snapshot_tables/repos.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"name",
			"url",
			"description",
			"owner_id",
			"language",
			"forked_from",
			"deleted",
		),
	)
	dataflowTester.VerifyTable(
		ticket.Board{},
		"./snapshot_tables/boards.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"name",
			"description",
			"url",
			"created_date",
		),
	)
	dataflowTester.VerifyTable(
		devops.CicdScope{},
		"./snapshot_tables/cicd_scopes.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"name",
			"description",
			"url",
			"created_date",
		),
	)
}
