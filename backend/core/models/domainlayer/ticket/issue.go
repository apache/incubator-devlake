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

package ticket

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/domainlayer"
)

type Issue struct {
	domainlayer.DomainEntity
	Url                     string `gorm:"type:varchar(255)"`
	IconURL                 string `gorm:"type:varchar(255);column:icon_url"`
	IssueKey                string `gorm:"type:varchar(255)"`
	Title                   string
	Description             string
	EpicKey                 string `gorm:"type:varchar(255)"`
	Type                    string `gorm:"type:varchar(100)"`
	OriginalType            string `gorm:"type:varchar(100)"`
	Status                  string `gorm:"type:varchar(100)"`
	OriginalStatus          string `gorm:"type:varchar(100)"`
	StoryPoint              float64
	ResolutionDate          *time.Time
	CreatedDate             *time.Time
	UpdatedDate             *time.Time
	LeadTimeMinutes         int64
	ParentIssueId           string `gorm:"type:varchar(255)"`
	Priority                string `gorm:"type:varchar(255)"`
	OriginalEstimateMinutes int64
	TimeSpentMinutes        int64
	TimeRemainingMinutes    int64
	CreatorId               string `gorm:"type:varchar(255)"`
	CreatorName             string `gorm:"type:varchar(255)"`
	AssigneeId              string `gorm:"type:varchar(255)"`
	AssigneeName            string `gorm:"type:varchar(255)"`
	Severity                string `gorm:"type:varchar(255)"`
	Component               string `gorm:"type:varchar(255)"`
	OriginalProject         string `gorm:"type:varchar(255)"`
}

func (Issue) TableName() string {
	return "issues"
}

const (
	BUG         = "BUG"
	REQUIREMENT = "REQUIREMENT"
	INCIDENT    = "INCIDENT"
	TASK        = "TASK"

	// status
	TODO        = "TODO"
	DONE        = "DONE"
	IN_PROGRESS = "IN_PROGRESS"
	OTHER       = "OTHER"
)

type StatusRule struct {
	InProgress []string
	Todo       []string
	Done       []string
	Other      []string
	Default    string
}

// GetStatus compare the input with rule for return the enmu value of status
func GetStatus(rule *StatusRule, input interface{}) string {
	for _, inp := range rule.InProgress {
		if inp == input {
			return IN_PROGRESS
		}
	}
	for _, Todo := range rule.Todo {
		if Todo == input {
			return TODO
		}
	}
	for _, done := range rule.Done {
		if done == input {
			return DONE
		}
	}
	for _, Other := range rule.Other {
		if Other == input {
			return OTHER
		}
	}
	return rule.Default
}
