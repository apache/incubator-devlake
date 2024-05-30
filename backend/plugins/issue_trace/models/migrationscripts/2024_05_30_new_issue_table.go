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

type NewIssueTable struct {
}

func (*NewIssueTable) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &IssueStatusHistory20240530{}, &IssueAssigneeHistory20240530{})
}

func (*NewIssueTable) Version() uint64 {
	return 20240530144400
}

func (*NewIssueTable) Name() string {
	return "add issue_status_history and issue_assignee_history"
}

type IssueStatusHistory20240530 struct {
	archived.NoPKModel
	IssueId           string     `gorm:"primaryKey;type:varchar(255)"`
	Status            string     `gorm:"type:varchar(100)"`
	OriginalStatus    string     `gorm:"primaryKey;type:varchar(255)"`
	StartDate         time.Time  `gorm:"primaryKey"`
	EndDate           *time.Time `gorm:"type:timestamp"`
	IsCurrentStatus   bool       `gorm:"type:boolean"`
	IsFirstStatus     bool       `gorm:"type:boolean"`
	StatusTimeMinutes int32      `gorm:"type:integer"`
}

func (IssueStatusHistory20240530) TableName() string {
	return "issue_status_history"
}

type IssueAssigneeHistory20240530 struct {
	archived.NoPKModel
	IssueId   string    `gorm:"primaryKey;type:varchar(255)"`
	Assignee  string    `gorm:"primaryKey;type:varchar(255)"`
	StartDate time.Time `gorm:"primaryKey"`
	EndDate   *time.Time
}

func (IssueAssigneeHistory20240530) TableName() string {
	return "issue_assignee_history"
}
