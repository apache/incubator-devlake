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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"time"
)

var _ plugin.MigrationScript = (*addCICDDeploymentsTable)(nil)

type addCICDDeploymentsTable struct {
}

type addCICDDeploymentsTable20230831 struct {
	archived.DomainEntity
	CicdScopeId      string `gorm:"index;type:varchar(255)"`
	CicdDeploymentId string `gorm:"type:varchar(255)"` // if it is converted from a cicd_pipeline_commit
	Name             string `gorm:"type:varchar(255)"`
	Result           string `gorm:"type:varchar(100)"`
	Status           string `gorm:"type:varchar(100)"`
	Environment      string `gorm:"type:varchar(255)"`
	CreatedDate      time.Time
	StartedDate      *time.Time
	FinishedDate     *time.Time
	DurationSec      *uint64
	RefName          string `gorm:"type:varchar(255)"` // to delete?
	RepoId           string `gorm:"type:varchar(255)"`
	RepoUrl          string `gorm:"index;not null"`
}

func (*addCICDDeploymentsTable20230831) TableName() string {
	return "cicd_deployments"
}

func (*addCICDDeploymentsTable) Up(basicRes context.BasicRes) errors.Error {
	// To create multiple tables with migration helper
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&addCICDDeploymentsTable20230831{},
	)
}

func (*addCICDDeploymentsTable) Version() uint64 {
	return 202307831162402
}

func (*addCICDDeploymentsTable) Name() string {
	return "add cicd deployments table"
}
