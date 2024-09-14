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
	"github.com/apache/incubator-devlake/core/models/domainlayer"
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

const ENV_NAME_PATTERN = "ENV_NAME_PATTERN"

type TransformDeployment struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type TransformDeploymentResponse struct {
	Total int                   `json:"total"`
	Data  []TransformDeployment `json:"data"`
}

type CICDTask struct {
	domainlayer.DomainEntity
	Name              string `gorm:"type:varchar(255)"`
	PipelineId        string `gorm:"index;type:varchar(255)"`
	Result            string `gorm:"type:varchar(100)"`
	Status            string `gorm:"type:varchar(100)"`
	OriginalStatus    string `gorm:"type:varchar(100)"`
	OriginalResult    string `gorm:"type:varchar(100)"`
	Type              string `gorm:"type:varchar(100);comment: to indicate this is CI or CD"`
	Environment       string `gorm:"type:varchar(255)"`
	DurationSec       float64
	QueuedDurationSec *float64
	TaskDatesInfo
	//StartedDate  time.Time  // notice here
	CicdScopeId string `gorm:"index;type:varchar(255)"`
}

func (CICDTask) TableName() string {
	return "cicd_tasks"
}
