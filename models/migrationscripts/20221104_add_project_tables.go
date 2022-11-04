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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*addProjectTables)(nil)

type addProjectTables struct{}

type Project struct {
	Name     string `gorm:"primaryKey;type:varchar(255)"`
	Describe string `gorm:"type:text"`

	archived.NoPKModel
}

func (Project) TableName() string {
	return "_devlake_projects"
}

type ProjectMetric struct {
	ProjectName  string `gorm:"primaryKey;type:varchar(255)"`
	PluginName   string `gorm:"primaryKey;type:varchar(255)"`
	PluginOption string `gorm:"type:text"`

	archived.NoPKModel
}

func (ProjectMetric) TableName() string {
	return "_devlake_project_metrics"
}

type ProjectMapping struct {
	ProjectName string `gorm:"primaryKey;type:varchar(255)"`
	Table       string `gorm:"primaryKey;type:varchar(255)"`
	RowId       string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

func (ProjectMapping) TableName() string {
	return "_devlake_project_mapping"
}

func (script *addProjectTables) Up(basicRes core.BasicRes) errors.Error {
	// To create multiple tables with migrationhelper
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&Project{},
		&ProjectMetric{},
		&ProjectMapping{},
	)
}

func (*addProjectTables) Version() uint64 {
	return 20221104194734
}

func (*addProjectTables) Name() string {
	return "add project tables"
}
