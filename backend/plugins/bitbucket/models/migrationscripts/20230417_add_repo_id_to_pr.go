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
)

type addRepoIdToPr struct{}

type pr20230417 struct {
	ConnectionId       uint64 `gorm:"primaryKey"`
	RepoId             string `gorm:"primaryKey;type:varchar(255)"` // repo_id should be a part of the primary key
	BitbucketId        int    `gorm:"primaryKey"`
	Number             int    `gorm:"index"`
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

func (pr20230417) TableName() string {
	return "_tool_bitbucket_pull_requests"
}

type prComment20230417 struct {
	ConnectionId       uint64 `gorm:"primaryKey"`
	BitbucketId        int    `gorm:"primaryKey"`
	RepoId             string `gorm:"index:pr"` // PullRequestId is not unique across multiple repos of a connection
	PullRequestId      int    `gorm:"index:pr"`
	AuthorId           string `gorm:"type:varchar(255)"`
	AuthorName         string `gorm:"type:varchar(255)"`
	BitbucketCreatedAt time.Time
	BitbucketUpdatedAt *time.Time
	Type               string `gorm:"comment:if type=null, it is normal comment,if type=diffNote,it is diff comment"`
	Body               string
	archived.NoPKModel
}

func (prComment20230417) TableName() string {
	return "_tool_bitbucket_pull_request_comments"
}

type prCommit20230417 struct {
	ConnectionId  uint64 `gorm:"primaryKey"`
	RepoId        string `gorm:"primaryKey"` // PullRequestId is not unique across multiple repos of a connection
	PullRequestId int    `gorm:"primaryKey;autoIncrement:false"`
	CommitSha     string `gorm:"primaryKey;type:varchar(40)"`
	archived.NoPKModel
}

func (prCommit20230417) TableName() string {
	return "_tool_bitbucket_pull_request_commits"
}

func (u *addRepoIdToPr) Up(basicRes context.BasicRes) errors.Error {
	tables := []interface{}{
		&pr20230417{},
		&prComment20230417{},
		&prCommit20230417{},
	}
	err := basicRes.GetDal().DropTables(tables...)
	if err != nil {
		return err
	}
	return migrationhelper.AutoMigrateTables(basicRes, tables...)
}

func (*addRepoIdToPr) Version() uint64 {
	return 20230417150357
}

func (*addRepoIdToPr) Name() string {
	return "add repo_id to pr as part of the primary key"
}
