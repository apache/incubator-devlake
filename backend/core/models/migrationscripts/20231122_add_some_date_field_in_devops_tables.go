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
	"time"
)

var _ plugin.MigrationScript = (*addSomeDateFieldsToDevopsTables)(nil)

type cicdDeployment20231122 struct {
	QueuedDate *time.Time
}

func (cicdDeployment20231122) TableName() string {
	return "cicd_deployments"
}

type cicdDeploymentCommit20231122 struct {
	QueuedDate *time.Time
}

func (cicdDeploymentCommit20231122) TableName() string {
	return "cicd_deployment_commits"
}

type cicdPipeline20231122 struct {
	StartedDate *time.Time
	QueuedDate  *time.Time
}

func (cicdPipeline20231122) TableName() string {
	return "cicd_pipelines"
}

type cicdTask20231122 struct {
	CreatedDate time.Time
	QueuedDate  *time.Time
}

func (cicdTask20231122) TableName() string {
	return "cicd_tasks"
}

type addSomeDateFieldsToDevopsTables struct{}

func (u *addSomeDateFieldsToDevopsTables) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes,
		&cicdPipeline20231122{},
		&cicdTask20231122{},
		&cicdDeployment20231122{},
		&cicdDeploymentCommit20231122{},
	)
}

func (*addSomeDateFieldsToDevopsTables) Version() uint64 {
	return 20231122140000
}

func (*addSomeDateFieldsToDevopsTables) Name() string {
	return "change duration_sec field to float64 in all related tables"
}
