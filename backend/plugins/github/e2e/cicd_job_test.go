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
	"github.com/apache/incubator-devlake/core/models/common"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestGithubCICDJobDataFlow(t *testing.T) {
	var github impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", github)
	regexEnricher := helper.NewRegexEnricher()
	_ = regexEnricher.TryAdd(devops.DEPLOYMENT, "deploywindows.*")
	_ = regexEnricher.TryAdd(devops.PRODUCTION, "deploywindows.*")
	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Name:         "panjf2000/ants",
			GithubId:     134018330,
		},
		RegexEnricher: regexEnricher,
	}

	// import raw data table
	// SELECT * FROM _raw_github_api_jobs INTO OUTFILE "/tmp/_raw_github_api_jobs.csv" FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY '\r\n';
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_jobs.csv", "_raw_github_api_jobs")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubJob{})
	dataflowTester.FlushTabler(&devops.CICDTask{})

	dataflowTester.Subtask(tasks.ExtractJobsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.GithubJob{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_github_jobs.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.Subtask(tasks.ConvertJobsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CICDTask{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_tasks.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
