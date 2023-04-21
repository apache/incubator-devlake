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

type addEnvToRunAndJob struct{}

type run20230418 struct {
	Environment string `gorm:"type:varchar(255)"`
}

func (run20230418) TableName() string {
	return "_tool_github_runs"
}

type runJob20230418 struct {
	Environment string `gorm:"type:varchar(255)"`
}

func (runJob20230418) TableName() string {
	return "_tool_github_jobs"
}

func (u *addEnvToRunAndJob) Up(baseRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(baseRes, &run20230418{}, &runJob20230418{})
}

func (*addEnvToRunAndJob) Version() uint64 {
	return 20230322150357
}

func (*addEnvToRunAndJob) Name() string {
	return "add type/env to github runs and run_jobs"
}
