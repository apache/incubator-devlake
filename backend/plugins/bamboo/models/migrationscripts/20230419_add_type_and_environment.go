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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addTypeAndEnvironment struct{}

type deployBuild20230419 struct {
	Environment string `gorm:"type:varchar(255)"`
}

func (deployBuild20230419) TableName() string {
	return "_tool_bamboo_deploy_build"
}

type jobBuild20230419 struct {
	Type        string `gorm:"type:varchar(255)"`
	Environment string `gorm:"type:varchar(255)"`
}

func (jobBuild20230419) TableName() string {
	return "_tool_bamboo_job_builds"
}

type planBuild20230419 struct {
	Type        string `gorm:"type:varchar(255)"`
	Environment string `gorm:"type:varchar(255)"`
}

func (planBuild20230419) TableName() string {
	return "_tool_bamboo_plan_builds"
}

func (u *addTypeAndEnvironment) Up(baseRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(baseRes, &jobBuild20230419{}, &planBuild20230419{}, &deployBuild20230419{})
}

func (*addTypeAndEnvironment) Version() uint64 {
	return 20230419141352
}

func (*addTypeAndEnvironment) Name() string {
	return "add type and environment to bamboo build tables"
}
