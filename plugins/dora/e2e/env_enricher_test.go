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
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/dora/impl"
	"github.com/apache/incubator-devlake/plugins/dora/tasks"
)

func TestEnrichEnvDataFlow(t *testing.T) {
	var plugin impl.Dora
	dataflowTester := e2ehelper.NewDataFlowTester(t, "dora", plugin)

	taskData := &tasks.DoraTaskData{
		Options: &tasks.DoraOptions{
			RepoId: "github:GithubRepo:1:384111310",
			TransformationRules: tasks.TransformationRules{
				ProductionPattern: "(?i)deploy",
				StagingPattern:    "(?i)stag",
				TestingPattern:    "(?i)test",
			},
		},
	}

	dataflowTester.FlushTabler(&devops.CICDTask{})

	// import raw data table
	dataflowTester.ImportCsvIntoTabler("./raw_tables/lake_cicd_pipeline_commits.csv", &devops.CiCDPipelineCommit{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/lake_cicd_tasks.csv", &devops.CICDTask{})

	// verify enrich with repoId
	dataflowTester.Subtask(tasks.EnrichTaskEnvMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CICDTask{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/lake_cicd_tasks.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// verify enrich with prefix
	dataflowTester.FlushTabler(&devops.CICDTask{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/lake_cicd_tasks_for_prefix.csv", &devops.CICDTask{})
	taskDataPrefix := &tasks.DoraTaskData{
		Options: &tasks.DoraOptions{
			TransformationRules: tasks.TransformationRules{
				ProductionPattern: "(?i)deploy",
				StagingPattern:    "(?i)stag",
				TestingPattern:    "(?i)test",
			},
			Prefix: "jenkins",
		},
	}
	dataflowTester.Subtask(tasks.EnrichTaskEnvMeta, taskDataPrefix)
	dataflowTester.VerifyTableWithOptions(&devops.CICDTask{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/lake_cicd_tasks_prefix.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
