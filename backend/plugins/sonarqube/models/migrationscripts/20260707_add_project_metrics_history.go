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
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addProjectMetricsHistory)(nil)

type projectMetricsHistory20260707 struct {
	ConnectionId uint64    `gorm:"primaryKey"`
	ProjectKey   string    `gorm:"primaryKey;type:varchar(255)"`
	AnalysisDate time.Time `gorm:"primaryKey"`
	MetricKey    string    `gorm:"primaryKey;type:varchar(100)"`
	MetricValue  string    `gorm:"type:varchar(50)"`
	archived.NoPKModel
}

func (projectMetricsHistory20260707) TableName() string {
	return "_tool_sonarqube_project_metrics_history"
}

type projectAnalyses20260707 struct {
	ConnectionId   uint64    `gorm:"primaryKey"`
	ProjectKey     string    `gorm:"primaryKey;type:varchar(255)"`
	AnalysisKey    string    `gorm:"primaryKey;type:varchar(255)"`
	AnalysisDate   time.Time `gorm:"index"`
	ProjectVersion string    `gorm:"type:varchar(255)"`
	Revision       string    `gorm:"type:varchar(255)"`
	BuildString    string    `gorm:"type:varchar(255)"`
	DetectedCI     string    `gorm:"type:varchar(100)"`
	archived.NoPKModel
}

func (projectAnalyses20260707) TableName() string {
	return "_tool_sonarqube_project_analyses"
}

type addProjectMetricsHistory struct{}

func (script *addProjectMetricsHistory) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&projectMetricsHistory20260707{},
		&projectAnalyses20260707{},
	)
}

func (*addProjectMetricsHistory) Version() uint64 {
	return 20260707153200
}

func (*addProjectMetricsHistory) Name() string {
	return "add project_metrics_history and project_analyses tables"
}
