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

type gitlabAccount20231120 struct {
	CreatedUserAt *time.Time
}

func (gitlabAccount20231120) TableName() string {
	return "_tool_gitlab_accounts"
}

type addUserCreationAt20231120 struct{}

func (*addUserCreationAt20231120) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &gitlabAccount20231120{})
}

func (*addUserCreationAt20231120) Version() uint64 {
	return 20231226140000
}

func (*addUserCreationAt20231120) Name() string {
	return "add created_user_at column to _tool_gitlab_accounts"
}
