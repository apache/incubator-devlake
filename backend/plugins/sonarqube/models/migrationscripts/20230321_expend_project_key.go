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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type SonarqubeProject20230206Before struct {
	common.NoPKModel `json:"-" mapstructure:"-"`
	ConnectionId     uint64           `json:"connectionId" validate:"required" gorm:"primaryKey"`
	ProjectKey       string           `json:"projectKey" validate:"required" gorm:"type:varchar(64);primaryKey"`
	Name             string           `json:"name" gorm:"type:varchar(255)"`
	Qualifier        string           `json:"qualifier" gorm:"type:varchar(255)"`
	Visibility       string           `json:"visibility" gorm:"type:varchar(64)"`
	LastAnalysisDate *api.Iso8601Time `json:"lastAnalysisDate"`
	Revision         string           `json:"revision" gorm:"type:varchar(128)"`
}

func (SonarqubeProject20230206Before) TableName() string {
	return "_tool_sonarqube_projects"
}

type SonarqubeProject20230206After struct {
	common.NoPKModel `json:"-" mapstructure:"-"`
	ConnectionId     uint64           `json:"connectionId" validate:"required" gorm:"primaryKey"`
	ProjectKey       string           `json:"projectKey" validate:"required" gorm:"type:varchar(255);primaryKey"` // expand this
	Name             string           `json:"name" gorm:"type:varchar(255)"`
	Qualifier        string           `json:"qualifier" gorm:"type:varchar(255)"`
	Visibility       string           `json:"visibility" gorm:"type:varchar(64)"`
	LastAnalysisDate *api.Iso8601Time `json:"lastAnalysisDate"`
	Revision         string           `json:"revision" gorm:"type:varchar(128)"`
}

func (SonarqubeProject20230206After) TableName() string {
	return "_tool_sonarqube_projects"
}

type SonarqubeIssue20230206Before struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	IssueKey     string `gorm:"primaryKey;type:varchar(100)"`
	Rule         string `gorm:"type:varchar(255)"`
	Severity     string `gorm:"type:varchar(100)"`
	Component    string `gorm:"type:varchar(255)"`
	ProjectKey   string `gorm:"index;type:varchar(100)"`
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

func (SonarqubeIssue20230206Before) TableName() string {
	return "_tool_sonarqube_issues"
}

type SonarqubeIssue20230206After struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	IssueKey     string `gorm:"primaryKey;type:varchar(100)"`
	Rule         string `gorm:"type:varchar(255)"`
	Severity     string `gorm:"type:varchar(100)"`
	Component    string `gorm:"type:varchar(255)"`
	ProjectKey   string `gorm:"index;type:varchar(255)"` // expand this
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

func (SonarqubeIssue20230206After) TableName() string {
	return "_tool_sonarqube_issues"
}

type expandProjectKey20230206 struct{}

func (script *expandProjectKey20230206) Up(basicRes context.BasicRes) errors.Error {
	// expand `ProjectKey` from varchar(64) to varchar(255)
	err := migrationhelper.TransformTable(
		basicRes,
		script,
		"_tool_sonarqube_projects",
		func(s *SonarqubeProject20230206Before) (*SonarqubeProject20230206After, errors.Error) {
			dst := (*SonarqubeProject20230206After)(s)
			return dst, nil
		},
	)
	if err != nil {
		return err
	}

	// expand `ProjectKey` from varchar(100) to varchar(255)
	err = migrationhelper.TransformTable(
		basicRes,
		script,
		"_tool_sonarqube_issues",
		func(s *SonarqubeIssue20230206Before) (*SonarqubeIssue20230206After, errors.Error) {
			dst := (*SonarqubeIssue20230206After)(s)
			return dst, nil
		},
	)
	if err != nil {
		return err
	}

	// also, SonarqubeFileMetrics and SonarqubeHotspot have ProjectKey.
	// But I think varchar(191) in mysql and text in pg is enough.
	return nil
}

func (*expandProjectKey20230206) Version() uint64 {
	return 20230321000003
}

func (*expandProjectKey20230206) Name() string {
	return "expend project_key"
}
