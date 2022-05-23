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

type Issue struct {
	DomainEntity
	Url                     string `gorm:"type:varchar(255)"`
	Number                  string `gorm:"type:varchar(255)"`
	Title                   string
	Description             string
	EpicKey                 string `gorm:"type:varchar(255)"`
	Type                    string `gorm:"type:varchar(100)"`
	Status                  string `gorm:"type:varchar(100)"`
	OriginalStatus          string `gorm:"type:varchar(100)"`
	StoryPoint              uint
	ResolutionDate          *time.Time
	CreatedDate             *time.Time
	UpdatedDate             *time.Time
	LeadTimeMinutes         uint
	ParentIssueId           string `gorm:"type:varchar(255)"`
	Priority                string `gorm:"type:varchar(255)"`
	OriginalEstimateMinutes int64
	TimeSpentMinutes        int64
	TimeRemainingMinutes    int64
	CreatorId               string `gorm:"type:varchar(255)"`
	AssigneeId              string `gorm:"type:varchar(255)"`
	AssigneeName            string `gorm:"type:varchar(255)"`
	Severity                string `gorm:"type:varchar(255)"`
	Component               string `gorm:"type:varchar(255)"`
}

type IssueCommit struct {
	NoPKModel
	IssueId   string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha string `gorm:"primaryKey;type:varchar(255)"`
}

type IssueLabel struct {
	IssueId   string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	LabelName string `gorm:"primaryKey;type:varchar(255)"`
	NoPKModel
}

type IssueComment struct {
	DomainEntity
	IssueId     string `gorm:"index"`
	Body        string
	UserId      string `gorm:"type:varchar(255)"`
	CreatedDate time.Time
}
