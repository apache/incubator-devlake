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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type SonarqubeIssue struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	IssueKey     string `gorm:"primaryKey;type:varchar(100)"`
	Rule         string `gorm:"type:varchar(255)"`
	Severity     string `gorm:"type:varchar(100)"`
	Component    string `gorm:"type:varchar(255)"`
	ProjectKey   string `gorm:"index;type:varchar(255)"` //domain project key
	Line         int
	Status       string `gorm:"type:varchar(20)"`
	Message      string
	Debt         int
	Effort       int
	Author       string `gorm:"type:varchar(100)"`
	Hash         string `gorm:"type:varchar(100)"`
	Tags         string
	Type         string `gorm:"type:varchar(100)"`
	Scope        string `gorm:"type:varchar(255)"`
	StartLine    int
	EndLine      int
	StartOffset  int
	EndOffset    int
	CreationDate *api.Iso8601Time
	UpdateDate   *api.Iso8601Time
	common.NoPKModel
}

func (SonarqubeIssue) TableName() string {
	return "_tool_sonarqube_issues"
}
