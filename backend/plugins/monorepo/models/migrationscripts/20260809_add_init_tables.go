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

var _ plugin.MigrationScript = (*addInitTables)(nil)

// subProjectDeployment20260809 is a version-frozen snapshot of
// models.SubProjectDeployment as it looked when this migration was written.
// Migration scripts must not import their plugin's live models package
// (see core/migration/linter), so the shape is duplicated here on purpose.
type subProjectDeployment20260809 struct {
	archived.NoPKModel
	ProjectName      string `gorm:"primaryKey;type:varchar(100)"`
	SubProject       string `gorm:"primaryKey;type:varchar(100)"`
	CicdDeploymentId string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha        string `gorm:"primaryKey;type:varchar(64)"`
	JobName          string `gorm:"type:varchar(255)"`
	Result           string `gorm:"type:varchar(100)"`
	Environment      string `gorm:"type:varchar(255)"`
	FinishedDate     *time.Time
}

func (subProjectDeployment20260809) TableName() string {
	return "monorepo_subproject_deployments"
}

// subProjectPrMetric20260809 is a version-frozen snapshot of
// models.SubProjectPrMetric as it looked when this migration was written.
type subProjectPrMetric20260809 struct {
	archived.NoPKModel
	ProjectName   string `gorm:"primaryKey;type:varchar(100)"`
	PullRequestId string `gorm:"primaryKey;type:varchar(255)"`
	SubProject    string `gorm:"index;type:varchar(255)"`

	CodingTime *int64
	PickupTime *int64
	ReviewTime *int64
	DeployTime *int64
	CycleTime  *int64

	DeploymentId string `gorm:"type:varchar(255)"`

	PrCreatedDate *time.Time
	PrMergedDate  *time.Time
	DeployedDate  *time.Time
}

func (subProjectPrMetric20260809) TableName() string {
	return "monorepo_subproject_pr_metrics"
}

type addInitTables struct{}

func (script *addInitTables) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&subProjectDeployment20260809{},
		&subProjectPrMetric20260809{},
	)
}

func (*addInitTables) Version() uint64 { return 20260809100000 }

func (*addInitTables) Name() string {
	return "create init tables for the monorepo plugin"
}
