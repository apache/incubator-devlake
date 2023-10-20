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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	githubArchived "github.com/apache/incubator-devlake/plugins/github/models/migrationscripts/archived"
)

type modifyGithubMilestone struct{}

type GithubMilestone20231018 struct {
	ConnectionId    uint64 `gorm:"primaryKey"`
	MilestoneId     int    `gorm:"primaryKey;autoIncrement:false"`
	RepoId          int
	Number          int
	URL             string
	Title           string
	OpenIssues      int
	ClosedIssues    int
	State           string
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time
	ClosedAt        *time.Time

	archived.NoPKModel
}

func (GithubMilestone20231018) TableName() string {
	return "_tool_github_milestones"
}

func (script *modifyGithubMilestone) Up(basicRes context.BasicRes) errors.Error {
	err := basicRes.GetDal().DropTables(&githubArchived.GithubMilestone{})
	if err != nil {
		return err
	}
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&GithubMilestone20231018{},
	)
}

func (*modifyGithubMilestone) Version() uint64 {
	return 20231018122537
}

func (*modifyGithubMilestone) Name() string {
	return "modify _tool_github_milestones table created_at and updated_at"
}
