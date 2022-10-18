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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/core"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

type addPipelineProjects struct{}

type GitlabPipelineProjects20220907 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	PipelineId   int    `gorm:"primaryKey"`
	ProjectId    int    `gorm:"primaryKey"`
	Ref          string `gorm:"type:varchar(255)"`
	Sha          string `gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (GitlabPipelineProjects20220907) TableName() string {
	return "_tool_gitlab_pipeline_projects"
}

func (*addPipelineProjects) Up(baseRes core.BasicRes) errors.Error {
	err := migrationhelper.AutoMigrateTables(baseRes, &GitlabPipelineProjects20220907{})
	if err != nil {
		return err
	}
	return nil
}

func (*addPipelineProjects) Version() uint64 {
	return 20220907230912
}

func (*addPipelineProjects) Name() string {
	return "gitlab add _tool_gitlab_pipeline_projects table"
}
