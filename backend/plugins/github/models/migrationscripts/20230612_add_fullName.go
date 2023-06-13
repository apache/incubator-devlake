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

type githubRepo20230612 struct {
	FullName string `gorm:"type:varchar(255)"`
}

func (githubRepo20230612) TableName() string {
	return "_tool_github_repos"
}

type addFullName struct{}

func (*addFullName) Up(res context.BasicRes) errors.Error {
	db := res.GetDal()
	return db.AutoMigrate(&githubRepo20230612{})
}

func (*addFullName) Version() uint64 {
	return 20230612184110
}

func (*addFullName) Name() string {
	return "add full_name to _tool_github_repos"
}
