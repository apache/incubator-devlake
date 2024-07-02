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

type Incident struct {
	DomainEntity
	Url                     string `gorm:"type:varchar(255)"`
	IncidentKey             string `gorm:"type:varchar(255)"` // issue_key/pull_request_key
	Title                   string
	Description             string
	Status                  string `gorm:"type:varchar(100)"`
	OriginalStatus          string `gorm:"type:varchar(100)"`
	ResolutionDate          *time.Time
	CreatedDate             *time.Time
	UpdatedDate             *time.Time
	LeadTimeMinutes         *uint
	OriginalEstimateMinutes *int64
	TimeSpentMinutes        *int64
	TimeRemainingMinutes    *int64
	CreatorId               string `gorm:"type:varchar(255)"`
	CreatorName             string `gorm:"type:varchar(255)"`
	ParentIncidentId        string `gorm:"type:varchar(255)"`
	Priority                string `gorm:"type:varchar(255)"`
	Severity                string `gorm:"type:varchar(255)"`
	Urgency                 string `gorm:"type:varchar(255)"`
	Component               string `gorm:"type:varchar(255)"`

	OriginalProject string `gorm:"type:varchar(255)"`
}

func (Incident) TableName() string {
	return "incidents"
}
