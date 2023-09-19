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

var _ plugin.MigrationScript = (*addEnvNamePattern)(nil)

type deployBuild20220807 struct {
	PlanResultKey  string `gorm:"primaryKey"`
	PlanBranchName string `gorm:"type:varchar(255)"`
}

func (deployBuild20220807) TableName() string {
	return "_tool_bamboo_deploy_build"
}

type planBuildVcsRevision20220807 struct {
	PlanResultKey string `gorm:"primaryKey"`
}

func (planBuildVcsRevision20220807) TableName() string {
	return "_tool_bamboo_plan_build_commits"
}

type addPlanResultKey struct{}

func (script *addPlanResultKey) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&deployBuild20220807{},
		&planBuildVcsRevision20220807{},
	)
}

func (*addPlanResultKey) Version() uint64 {
	return 20230807165119
}

func (script *addPlanResultKey) Name() string {
	return "add plan_result_key to _tool_bamboo_deploy_builds and _tool_bamboo_plan_build_commits"
}
