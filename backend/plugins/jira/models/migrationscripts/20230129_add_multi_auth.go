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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type jiraMultiAuth20230129 struct {
	AuthMethod string `gorm:"type:varchar(20)"`
	Token      string `gorm:"type:varchar(255)"`
}

func (jiraMultiAuth20230129) TableName() string {
	return "_tool_jira_connections"
}

type addJiraMultiAuth20230129 struct{}

func (script *addJiraMultiAuth20230129) Up(basicRes context.BasicRes) errors.Error {
	err := migrationhelper.AutoMigrateTables(basicRes, &jiraMultiAuth20230129{})
	if err != nil {
		return err
	}
	return basicRes.GetDal().UpdateColumn(
		&jiraMultiAuth20230129{},
		"auth_method", plugin.AUTH_METHOD_BASIC,
		dal.Where("auth_method IS NULL"),
	)
}

func (*addJiraMultiAuth20230129) Version() uint64 {
	return 20230129115901
}

func (*addJiraMultiAuth20230129) Name() string {
	return "add multiauth to _tool_jira_connections"
}
