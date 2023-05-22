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

func TestGenerateDeploymentCommitsDataFlow(t *testing.T) {
	var plugin impl.Dora
	dataflowTester := e2ehelper.NewDataFlowTester(t, "dora", plugin)

	taskData := &tasks.DoraTaskData{
		Options: &tasks.DoraOptions{
			ProjectName: "project1",
		},
	}
	// import raw data table
	dataflowTester.ImportCsvIntoTabler("./deployment_generator/project_mapping.csv", &crossdomain.ProjectMapping{})
	dataflowTester.ImportCsvIntoTabler("./deployment_generator/cicd_pipeline_commits.csv", &devops.CiCDPipelineCommit{})
	dataflowTester.ImportCsvIntoTabler("./deployment_generator/cicd_pipelines.csv", &devops.CICDPipeline{})
	dataflowTester.ImportCsvIntoTabler("./deployment_generator/cicd_tasks.csv", &devops.CICDTask{})

	// verify converter
	dataflowTester.FlushTabler(&devops.CicdDeploymentCommit{})
	dataflowTester.Subtask(tasks.DeploymentCommitsGeneratorMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CicdDeploymentCommit{}, e2ehelper.TableOptions{
		CSVRelPath: "./deployment_generator/cicd_deployment_commits.csv",
		TargetFields: []string{
			"id",
			"cicd_scope_id",
			"cicd_deployment_id",
			"result",
			"repo_url",
			"environment",
			"started_date",
			"finished_date",
		},
	})
}
