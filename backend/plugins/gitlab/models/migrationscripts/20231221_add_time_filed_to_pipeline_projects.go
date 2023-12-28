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
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type gitlabPipelineProject20231221 struct {
	GitlabCreatedAt *time.Time
	GitlabUpdatedAt *time.Time
}

func (gitlabPipelineProject20231221) TableName() string {
	return "_tool_gitlab_pipeline_projects"
}

type addTimeToGitlabPipelineProject struct{}

func (script *addTimeToGitlabPipelineProject) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&gitlabPipelineProject20231221{},
	)
}

func (*addTimeToGitlabPipelineProject) Version() uint64 {
	return 20231221161031
}

func (*addTimeToGitlabPipelineProject) Name() string {
	return "add time filed to _tool_gitlab_pipeline_projects"
}
