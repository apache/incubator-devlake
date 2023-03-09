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
	"github.com/apache/incubator-devlake/plugins/bitbucket/impl"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
)

func TestBitbucketPipelineStepsDataFlow(t *testing.T) {

	var bitbucket impl.Bitbucket
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bitbucket", bitbucket)

	taskData := &tasks.BitbucketTaskData{
		Options: &tasks.BitbucketOptions{
			ConnectionId: 4,
			FullName:     "thenicetgp/lake",
			BitbucketTransformationRule: &models.BitbucketTransformationRule{
				DeploymentPattern: "",
				ProductionPattern: "",
			},
		},
	}

	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_bitbucket_pipeline_steps.csv", &models.BitbucketPipelineStep{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_bitbucket_deployments.csv", &models.BitbucketDeployment{})
	dataflowTester.FlushTabler(&devops.CICDTask{})
	// verify extraction
	dataflowTester.Subtask(tasks.ConvertPipelineStepMeta, taskData)
	dataflowTester.VerifyTable(
		devops.CICDTask{},
		"./snapshot_tables/cicd_tasks.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"name",
			"pipeline_id",
			"result",
			"status",
			"type",
			"duration_sec",
			"started_date",
			"finished_date",
			"environment",
			"cicd_scope_id",
		),
	)
}
