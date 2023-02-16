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

	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type SonarqubeProject struct {
	ConnectionId     uint64     `json:"connectionId" gorm:"primaryKey"`
	ProjectKey       string     `json:"projectKey" gorm:"type:varchar(64);primaryKey"`
	Name             string     `json:"name" gorm:"type:varchar(255)"`
	Qualifier        string     `json:"qualifier" gorm:"type:varchar(255)"`
	Visibility       string     `json:"visibility" gorm:"type:varchar(64)"`
	LastAnalysisDate *time.Time `json:"lastAnalysisDate"`
	Revision         string     `json:"revision" gorm:"type:varchar(128)"`
	archived.NoPKModel
}

func (SonarqubeProject) TableName() string {
	return "_tool_sonarqube_projects"
}
