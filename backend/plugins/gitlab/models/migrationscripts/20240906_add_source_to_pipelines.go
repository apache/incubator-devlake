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
)

type gitlabPipelineProject240906 struct {
	Source string
}

func (gitlabPipelineProject240906) TableName() string {
	return "_tool_gitlab_pipeline_projects"
}

type gitlabPipeline240906 struct {
	Source string
}

func (gitlabPipeline240906) TableName() string {
	return "_tool_gitlab_pipelines"
}

type addIsChildToPipelines240906 struct{}

func (*addIsChildToPipelines240906) Up(baseRes context.BasicRes) errors.Error {
	db := baseRes.GetDal()
	if err := db.AutoMigrate(&gitlabPipeline240906{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&gitlabPipelineProject240906{}); err != nil {
		return err
	}
	return nil
}

func (*addIsChildToPipelines240906) Version() uint64 {
	return 20240906150000
}

func (*addIsChildToPipelines240906) Name() string {
	return "add is_child to table _tool_gitlab_pipelines and _tool_gitlab_pipeline_projects"
}
