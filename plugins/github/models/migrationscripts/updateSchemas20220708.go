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
	"context"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
	"time"
)

// GithubMilestone20220620 new struct for milestones
type GithubMilestone20220620 struct {
	archived.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	MilestoneId  int    `gorm:"primaryKey;autoIncrement:false"`
	RepoId       int
	Number       int
	URL          string
	OpenIssues   int
	ClosedIssues int
	State        string
	Title        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ClosedAt     time.Time
}

// GithubIssue20220708 new field for models.GithubIssue
type GithubIssue20220708 struct {
	MilestoneId int
}

type UpdateSchemas20220708 struct{}

func (GithubMilestone20220620) TableName() string {
	return "_tool_github_milestones"
}

func (GithubIssue20220708) TableName() string {
	return "_tool_github_issues"
}

func (*UpdateSchemas20220708) Up(_ context.Context, db *gorm.DB) error {
	err := db.Migrator().AddColumn(GithubIssue20220708{}, "milestone_id")
	if err != nil {
		return err
	}
	return db.Migrator().CreateTable(GithubMilestone20220620{})
}

func (*UpdateSchemas20220708) Version() uint64 {
	return 20220708000001
}

func (*UpdateSchemas20220708) Name() string {
	return "Add milestone for github"
}
