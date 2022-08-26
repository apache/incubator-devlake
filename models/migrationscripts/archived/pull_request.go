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

package archived

import (
	"time"
)

type PullRequest struct {
	DomainEntity
	BaseRepoId  string `gorm:"index"`
	HeadRepoId  string `gorm:"index"`
	Status      string `gorm:"type:varchar(100);comment:open/closed or other"`
	Title       string
	Description string
	Url         string `gorm:"type:varchar(255)"`
	AuthorName  string `gorm:"type:varchar(100)"`
	//User		   domainUser.User `gorm:"foreignKey:AuthorId"`
	AuthorId       string `gorm:"type:varchar(100)"`
	ParentPrId     string `gorm:"index;type:varchar(100)"`
	PullRequestKey int
	CreatedDate    time.Time
	MergedDate     *time.Time
	ClosedDate     *time.Time
	Type           string `gorm:"type:varchar(100)"`
	Component      string `gorm:"type:varchar(100)"`
	MergeCommitSha string `gorm:"type:varchar(40)"`
	HeadRef        string `gorm:"type:varchar(255)"`
	BaseRef        string `gorm:"type:varchar(255)"`
	BaseCommitSha  string `gorm:"type:varchar(40)"`
	HeadCommitSha  string `gorm:"type:varchar(40)"`
}

func (PullRequest) TableName() string {
	return "pull_requests"
}

type PullRequestCommit struct {
	CommitSha     string `gorm:"primaryKey;type:varchar(40)"`
	PullRequestId string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	NoPKModel
}

func (PullRequestCommit) TableName() string {
	return "pull_request_commits"
}

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type PullRequestLabel struct {
	PullRequestId string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	LabelName     string `gorm:"primaryKey;type:varchar(255)"`
	NoPKModel
}

func (PullRequestLabel) TableName() string {
	return "pull_request_labels"
}

type PullRequestComment struct {
	DomainEntity
	PullRequestId string `gorm:"index"`
	Body          string
	AccountId     string `gorm:"type:varchar(255)"`
	CreatedDate   time.Time
	CommitSha     string `gorm:"type:varchar(255)"`
	Position      int
	Type          string `gorm:"type:varchar(255)"`
	ReviewId      string `gorm:"type:varchar(255)"`
	Status        string `gorm:"type:varchar(255)"`
}

func (PullRequestComment) TableName() string {
	return "pull_request_comments"
}

type PullRequestIssue struct {
	PullRequestId     string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	IssueId           string `gorm:"primaryKey;type:varchar(255)"`
	PullRequestNumber int
	IssueNumber       int
	NoPKModel
}

func (PullRequestIssue) TableName() string {
	return "pull_request_issues"
}
