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
)

var _ plugin.MigrationScript = (*addCicdDeploymentSubprojects)(nil)

// cicdDeploymentSubproject20260810 is a version-frozen snapshot of
// devops.CicdDeploymentSubproject as it looked when this migration was written.
// Migration scripts must not import live model packages (see core/migration/linter),
// so the shape is duplicated here on purpose rather than imported from the domain layer.
type cicdDeploymentSubproject20260810 struct {
	archived.NoPKModel
	ProjectName      string `gorm:"primaryKey;type:varchar(100)"`
	CicdDeploymentId string `gorm:"primaryKey;type:varchar(255);index:idx_cds_deployment"`
	SubProject       string `gorm:"primaryKey;type:varchar(100)"`
}

func (cicdDeploymentSubproject20260810) TableName() string {
	return "cicd_deployment_subprojects"
}

type addCicdDeploymentSubprojects struct{}

// Up creates the new cicd_deployment_subprojects mapping table.
func (script *addCicdDeploymentSubprojects) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&cicdDeploymentSubproject20260810{},
	)
}

func (*addCicdDeploymentSubprojects) Version() uint64 {
	return 20260810100100
}

func (*addCicdDeploymentSubprojects) Name() string {
	return "create cicd_deployment_subprojects mapping table for monorepo support"
}
