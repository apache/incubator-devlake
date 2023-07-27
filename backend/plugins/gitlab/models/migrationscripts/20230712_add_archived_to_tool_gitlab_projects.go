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

type gitlabProject20230711 struct {
	Archived bool `json:"archived"`
}

func (gitlabProject20230711) TableName() string {
	return "_tool_gitlab_projects"
}

type addProjectArchived struct{}

func (script *addProjectArchived) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&gitlabProject20230711{})
}

func (*addProjectArchived) Version() uint64 {
	return 20230712095900
}

func (*addProjectArchived) Name() string {
	return "add archived to _tool_gitlab_projects"
}
