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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addGiteeCommitAuthorInfo)(nil)

type GiteePrCommit20230419 struct {
	CommitAuthorName   string `gorm:"type:varchar(255)"` // Author name
	CommitAuthorEmail  string `gorm:"type:varchar(255)"` // Author email
	CommitAuthoredDate time.Time
}

func (GiteePrCommit20230419) TableName() string {
	return "_tool_gitee_pull_request_commits"
}

type addGiteeCommitAuthorInfo struct{}

func (script *addGiteeCommitAuthorInfo) Up(basicRes context.BasicRes) errors.Error {

	return migrationhelper.AutoMigrateTables(
		basicRes,
		&GiteePrCommit20230419{},
	)
}

func (*addGiteeCommitAuthorInfo) Version() uint64 {
	return 20230420135129
}

func (*addGiteeCommitAuthorInfo) Name() string {
	return "add commit author info to _tool_gitee_pull_request_commits table"
}
