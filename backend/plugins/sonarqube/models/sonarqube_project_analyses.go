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

package models

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

type SonarqubeProjectAnalysis struct {
	ConnectionId   uint64    `gorm:"primaryKey"`
	ProjectKey     string    `gorm:"primaryKey;type:varchar(255)"`
	AnalysisKey    string    `gorm:"primaryKey;type:varchar(255)"`
	AnalysisDate   time.Time `gorm:"index"`
	ProjectVersion string    `gorm:"type:varchar(255)"`
	Revision       string    `gorm:"type:varchar(255)"`
	BuildString    string    `gorm:"type:varchar(255)"`
	DetectedCI     string    `gorm:"type:varchar(100)"`
	common.NoPKModel
}

func (SonarqubeProjectAnalysis) TableName() string {
	return "_tool_sonarqube_project_analyses"
}
