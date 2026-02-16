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

type githubConnection20260211 struct {
	TokenExpiresAt        *time.Time
	RefreshTokenExpiresAt *time.Time
}

func (githubConnection20260211) TableName() string {
	return "_tool_github_connections"
}

type modifyTokenExpiresAtToNullable struct{}

func (*modifyTokenExpiresAtToNullable) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&githubConnection20260211{},
	)
}

func (*modifyTokenExpiresAtToNullable) Version() uint64 {
	return 20260211000001
}

func (*modifyTokenExpiresAtToNullable) Name() string {
	return "modify token_expires_at and refresh_token_expires_at to nullable"
}
