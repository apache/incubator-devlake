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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	archivedModel "github.com/apache/incubator-devlake/plugins/bitbucket/models/migrationscripts/archived"
	"time"
)

type BitbucketRepo20230206 struct {
	TransformationRuleId uint64 `json:"transformationRuleId,omitempty" mapstructure:"transformationRuleId,omitempty"`
	CloneUrl             string `json:"cloneUrl" gorm:"type:varchar(255)" mapstructure:"cloneUrl,omitempty"`
	Owner                string `json:"owner" mapstructure:"owner,omitempty"`
}

func (BitbucketRepo20230206) TableName() string {
	return "_tool_bitbucket_repos"
}

type BitbucketIssue20230206 struct {
	ConnectionId       uint64 `gorm:"primaryKey"`
	RepoId             string `gorm:"primaryKey;type:varchar(255)"`
	BitbucketId        int    `gorm:"primaryKey"`
	Number             int    `gorm:"index;comment:Used in API requests ex. api/issues/<THIS_NUMBER>"`
	State              string `gorm:"type:varchar(255)"`
	StdState           string `gorm:"type:varchar(255)"`
	Title              string `gorm:"type:varchar(255)"`
	Body               string
	Priority           string `gorm:"type:varchar(255)"`
	Type               string `gorm:"type:varchar(100)"`
	AuthorId           string `gorm:"type:varchar(255)"`
	AuthorName         string `gorm:"type:varchar(255)"`
	AssigneeId         string `gorm:"type:varchar(255)"`
	AssigneeName       string `gorm:"type:varchar(255)"`
	MilestoneId        int    `gorm:"index"`
	LeadTimeMinutes    uint
	Url                string `gorm:"type:varchar(255)"`
	ClosedAt           *time.Time
	BitbucketCreatedAt time.Time
	BitbucketUpdatedAt time.Time `gorm:"index"`
	Severity           string    `gorm:"type:varchar(255)"`
	Component          string    `gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (BitbucketIssue20230206) TableName() string {
	return "_tool_bitbucket_issues"
}

type BitbucketPullRequest20230206 struct {
	ConnectionId       uint64 `gorm:"primaryKey"`
	RepoId             string `gorm:"primaryKey;type:varchar(255)"`
	BitbucketId        int    `gorm:"primaryKey"`
	Number             int    `gorm:"index"` // This number is used in GET requests to the API associated to reviewers / comments / etc.
	BaseRepoId         string
	HeadRepoId         string
	State              string `gorm:"type:varchar(255)"`
	Title              string
	Description        string
	BitbucketCreatedAt time.Time
	BitbucketUpdatedAt time.Time `gorm:"index"`
	ClosedAt           *time.Time
	CommentCount       int
	Commits            int
	MergedAt           *time.Time
	Body               string
	Type               string `gorm:"type:varchar(255)"`
	Component          string `gorm:"type:varchar(255)"`
	MergeCommitSha     string `gorm:"type:varchar(40)"`
	HeadRef            string `gorm:"type:varchar(255)"`
	BaseRef            string `gorm:"type:varchar(255)"`
	BaseCommitSha      string `gorm:"type:varchar(255)"`
	HeadCommitSha      string `gorm:"type:varchar(255)"`
	Url                string `gorm:"type:varchar(255)"`
	AuthorName         string `gorm:"type:varchar(255)"`
	AuthorId           string `gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (BitbucketPullRequest20230206) TableName() string {
	return "_tool_bitbucket_pull_requests"
}

type addScope20230206 struct{}

func (script *addScope20230206) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.RenameColumn("_tool_bitbucket_repos", "owner_id", "owner")
	if err != nil {
		return err
	}

	// add `RepoId` as primary key
	err = migrationhelper.TransformTable(
		basicRes,
		script,
		"_tool_bitbucket_issues",
		func(s *archivedModel.BitbucketIssue) (*BitbucketIssue20230206, errors.Error) {
			dst := &BitbucketIssue20230206{
				ConnectionId:       s.ConnectionId,
				RepoId:             s.RepoId,
				BitbucketId:        s.BitbucketId,
				Number:             s.Number,
				State:              ``,
				StdState:           s.State,
				Title:              s.Title,
				Body:               s.Body,
				Priority:           s.Priority,
				Type:               s.Type,
				AuthorId:           s.AuthorId,
				AuthorName:         s.AuthorName,
				AssigneeId:         s.AssigneeId,
				AssigneeName:       s.AssigneeName,
				MilestoneId:        s.MilestoneId,
				LeadTimeMinutes:    s.LeadTimeMinutes,
				Url:                s.Url,
				ClosedAt:           s.ClosedAt,
				BitbucketCreatedAt: s.BitbucketCreatedAt,
				BitbucketUpdatedAt: s.BitbucketUpdatedAt,
				Severity:           s.Severity,
				Component:          s.Component,
				NoPKModel:          s.NoPKModel,
			}
			return dst, nil
		},
	)
	if err != nil {
		return err
	}

	// add `RepoId` as primary key
	err = migrationhelper.TransformTable(
		basicRes,
		script,
		"_tool_bitbucket_pull_requests",
		func(s *archivedModel.BitbucketPullRequest) (*BitbucketPullRequest20230206, errors.Error) {
			dst := &BitbucketPullRequest20230206{
				ConnectionId:       s.ConnectionId,
				RepoId:             s.RepoId,
				BitbucketId:        s.BitbucketId,
				Number:             s.Number,
				BaseRepoId:         s.BaseRepoId,
				HeadRepoId:         s.HeadRepoId,
				State:              s.State,
				Title:              s.Title,
				Description:        s.Description,
				BitbucketCreatedAt: s.BitbucketCreatedAt,
				BitbucketUpdatedAt: s.BitbucketUpdatedAt,
				ClosedAt:           s.ClosedAt,
				CommentCount:       s.CommentCount,
				Commits:            s.Commits,
				MergedAt:           s.MergedAt,
				Body:               s.Body,
				Type:               s.Type,
				Component:          s.Component,
				MergeCommitSha:     s.MergeCommitSha,
				HeadRef:            s.HeadRef,
				BaseRef:            s.BaseRef,
				BaseCommitSha:      s.BaseCommitSha,
				HeadCommitSha:      s.HeadCommitSha,
				Url:                s.Url,
				AuthorName:         s.AuthorName,
				AuthorId:           s.AuthorId,
				NoPKModel:          s.NoPKModel,
			}
			return dst, nil
		},
	)
	if err != nil {
		return err
	}

	return migrationhelper.AutoMigrateTables(
		basicRes,
		&BitbucketRepo20230206{},
		&archivedModel.BitbucketTransformationRule{},
	)
}

func (*addScope20230206) Version() uint64 {
	return 20230206000008
}

func (*addScope20230206) Name() string {
	return "add scope and table _tool_bitbucket_transformation_rules"
}
