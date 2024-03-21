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
	"github.com/apache/incubator-devlake/plugins/jenkins/impl"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"

	api "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func TestJenkinsJobsDataFlow(t *testing.T) {

	var jenkins impl.Jenkins
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jenkins", jenkins)

	taskData := &tasks.JenkinsTaskData{
		Options: &tasks.JenkinsOptions{
			ConnectionId: 1,
			JobName:      `devlake`,
			JobFullName:  `Test-jenkins-dir/test-jenkins-sub-dir/test-sub-sub-dir/devlake`,
			JobPath:      `job/Test-jenkins-dir/job/test-jenkins-sub-dir/job/test-sub-sub-dir/`,
			ScopeConfig:  new(models.JenkinsScopeConfig),
		},
	}

	dataflowTester.FlushTabler(&devops.CicdScope{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/_tool_jenkins_jobs.csv", &models.JenkinsJob{})
	dataflowTester.Subtask(tasks.ConvertJobsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CicdScope{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_scopes.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}

func TestJenkinsMultibranchJobsDataFlow(t *testing.T) {
	var jenkins impl.Jenkins
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jenkins", jenkins)

	taskData := &tasks.JenkinsTaskData{
		Options: &tasks.JenkinsOptions{
			ConnectionId: 1,
			JobName:      `devlake-jenkins`,
			JobFullName:  `github_org/devlake-jenkins`,
			JobPath:      ``,
			Class:        `org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject`,
			ScopeConfig:  new(models.JenkinsScopeConfig),
		},
		Connection: &models.JenkinsConnection{
			JenkinsConn: models.JenkinsConn{
				RestConnection: api.RestConnection{
					Endpoint: "https://1457-62-195-68-26.ngrok-free.app",
				},
			},
		},
	}

	dataflowTester.FlushTabler(&models.JenkinsJob{})
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jenkins_api_jobs_multibranch.csv", "_raw_jenkins_api_jobs")
	dataflowTester.Subtask(tasks.ExtractApiJobsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.JenkinsJob{}, e2ehelper.TableOptions{
		CSVRelPath:  "./raw_tables/_tool_jenkins_jobs_multibranch.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.FlushTabler(&devops.CicdScope{})
	dataflowTester.Subtask(tasks.ConvertJobsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&devops.CicdScope{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/cicd_scopes_multibranch.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
