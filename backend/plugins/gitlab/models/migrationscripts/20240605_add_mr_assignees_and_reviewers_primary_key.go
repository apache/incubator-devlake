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
	archivedCore "github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/gitlab/models/migrationscripts/archived"
)

type addGitlabAssigneeAndReviewerPrimaryKey struct{}

type GitlabAssignee20240605 struct {
	ConnectionId   uint64 `gorm:"primaryKey"`
	AssigneeId     int    `gorm:"primaryKey"`
	MergeRequestId int    `gorm:"primaryKey"`
	ProjectId      int    `gorm:"index"`
	Name           string `gorm:"type:varchar(255)"`
	Username       string `gorm:"type:varchar(255)"`
	State          string `gorm:"type:varchar(255)"`
	AvatarUrl      string `gorm:"type:varchar(255)"`
	WebUrl         string `gorm:"type:varchar(255)"`
	archivedCore.NoPKModel
}

func (GitlabAssignee20240605) TableName() string {
	return "_tool_gitlab_assignees"
}

type GitlabReviewer20240605 struct {
	ConnectionId   uint64 `gorm:"primaryKey"`
	ReviewerId     int    `gorm:"primaryKey"`
	MergeRequestId int    `gorm:"primaryKey"`
	ProjectId      int    `gorm:"index"`
	Name           string `gorm:"type:varchar(255)"`
	Username       string `gorm:"type:varchar(255)"`
	State          string `gorm:"type:varchar(255)"`
	AvatarUrl      string `gorm:"type:varchar(255)"`
	WebUrl         string `gorm:"type:varchar(255)"`
	archivedCore.NoPKModel
}

func (GitlabReviewer20240605) TableName() string {
	return "_tool_gitlab_reviewers"
}

func (*addGitlabAssigneeAndReviewerPrimaryKey) Up(baseRes context.BasicRes) errors.Error {
	err := baseRes.GetDal().DropTables(archived.GitlabAssignee{}, archived.GitlabReviewer{})
	if err != nil {
		return err
	}

	err = migrationhelper.AutoMigrateTables(
		baseRes,
		&GitlabAssignee20240605{},
		&GitlabReviewer20240605{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (*addGitlabAssigneeAndReviewerPrimaryKey) Version() uint64 {
	return 20240605110339
}

func (*addGitlabAssigneeAndReviewerPrimaryKey) Name() string {
	return "add primary key to _tool_gitlab_assignees and _tool_gitlab_reviewers tables"
}
