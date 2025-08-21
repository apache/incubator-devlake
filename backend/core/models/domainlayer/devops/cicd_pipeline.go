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

package devops

import (
	"strings"

	"github.com/spf13/cast"

	"github.com/apache/incubator-devlake/core/models/domainlayer"
)

type CICDPipeline struct {
	domainlayer.DomainEntity
	Name              string `gorm:"type:varchar(255)"`
	DisplayTitle      string
	Url               string
	Result            string `gorm:"type:varchar(100)"`
	Status            string `gorm:"type:varchar(100)"`
	OriginalStatus    string `gorm:"type:varchar(100)"`
	OriginalResult    string `gorm:"type:varchar(100)"`
	Type              string `gorm:"type:varchar(100);comment: to indicate this is CI or CD"`
	DurationSec       float64
	QueuedDurationSec *float64
	Environment       string `gorm:"type:varchar(255)"`
	TaskDatesInfo
	CicdScopeId string `gorm:"index;type:varchar(255)"`
	IsChild     bool
}

func (CICDPipeline) TableName() string {
	return "cicd_pipelines"
}

// this is for the field `result` in table.cicd_pipelines and table.cicd_tasks
const (
	RESULT_SUCCESS = "SUCCESS"
	RESULT_FAILURE = "FAILURE"
	RESULT_DEFAULT = ""
)

// this is for the field `status` in table.cicd_pipelines and table.cicd_tasks
const (
	STATUS_IN_PROGRESS = "IN_PROGRESS"
	STATUS_DONE        = "DONE"
	STATUS_OTHER       = "OTHER"
)

type ResultRule struct {
	Success []string
	Failure []string
	Default string
}

type StatusRule struct {
	InProgress []string
	Done       []string
	Default    string
}

type StatusRuleCommon[T comparable] struct {
	InProgress []T
	Done       []T
	Default    string
}

// GetResult compare the input with rule for return the enum value of result case-insensitively.
func GetResult(rule *ResultRule, input interface{}) string {
	for _, suc := range rule.Success {
		if strings.EqualFold(suc, cast.ToString(input)) {
			return RESULT_SUCCESS
		}
	}
	for _, fail := range rule.Failure {
		if strings.EqualFold(fail, cast.ToString(input)) {
			return RESULT_FAILURE
		}
	}
	return rule.Default
}

// GetStatus compare the input with rule for return the enum value of status
func GetStatus(rule *StatusRule, input interface{}) string {
	for _, inProgress := range rule.InProgress {
		if strings.EqualFold(inProgress, cast.ToString(input)) {
			return STATUS_IN_PROGRESS
		}
	}
	for _, done := range rule.Done {
		if strings.EqualFold(done, cast.ToString(input)) {
			return STATUS_DONE
		}
	}
	return rule.Default
}

// GetStatusCommon compare the input with rule for return the enum value of status.
// If T is string, it is case-sensitivity.
func GetStatusCommon[T comparable](rule *StatusRuleCommon[T], input T) string {
	if rule.Default == "" {
		rule.Default = STATUS_OTHER
	}
	for _, inp := range rule.InProgress {
		if inp == input {
			return STATUS_IN_PROGRESS
		}
	}
	for _, done := range rule.Done {
		if done == input {
			return STATUS_DONE
		}
	}
	return rule.Default
}
