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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addSomeDateFieldsToDevopsTables)(nil)

type cicdDeployment20231128 struct {
	QueuedDurationSec *float64
}

func (cicdDeployment20231128) TableName() string {
	return "cicd_deployments"
}

type cicdDeploymentCommit20231128 struct {
	QueuedDurationSec *float64
}

func (cicdDeploymentCommit20231128) TableName() string {
	return "cicd_deployment_commits"
}

type cicdPipeline20231128 struct {
	QueuedDurationSec *float64
}

func (cicdPipeline20231128) TableName() string {
	return "cicd_pipelines"
}

type cicdTask20231128 struct {
	QueuedDurationSec *float64
}

func (cicdTask20231128) TableName() string {
	return "cicd_tasks"
}

type addQueuedDurationSecFieldToDevopsTables struct{}

func (u *addQueuedDurationSecFieldToDevopsTables) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes,
		&cicdPipeline20231128{},
		&cicdTask20231128{},
		&cicdDeployment20231128{},
		&cicdDeploymentCommit20231128{},
	)
}

func (*addQueuedDurationSecFieldToDevopsTables) Version() uint64 {
	return 20231128140000
}

func (*addQueuedDurationSecFieldToDevopsTables) Name() string {
	return "add queued_duration_sec field to devops tables"
}
