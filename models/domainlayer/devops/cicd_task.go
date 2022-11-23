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
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
)

const (
	TEST       = "TEST"
	LINT       = "LINT"
	BUILD      = "BUILD"
	DEPLOYMENT = "DEPLOYMENT"
)

const (
	PRODUCTION = "PRODUCTION"
	STAGING    = "STAGING"
	TESTING    = "TESTING"
)

type CICDTask struct {
	domainlayer.DomainEntity
	Name         string `gorm:"type:varchar(255)"`
	PipelineId   string `gorm:"index;type:varchar(255)"`
	Result       string `gorm:"type:varchar(100)"`
	Status       string `gorm:"type:varchar(100)"`
	Type         string `gorm:"type:varchar(100);comment: to indicate this is CI or CD"`
	Environment  string `gorm:"type:varchar(255)"`
	DurationSec  uint64
	StartedDate  time.Time
	FinishedDate *time.Time
	CicdScopeId  string `gorm:"index;type:varchar(255)"`
}

func (CICDTask) TableName() string {
	return "cicd_tasks"
}
