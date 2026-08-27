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
	"github.com/apache/incubator-devlake/core/models/common"
)

// CicdDeploymentSubproject maps a deployment (identified by its pipeline id,
// CicdDeploymentId) to the monorepo sub-project(s) it deployed. The relationship is
// many-to-many: a single pipeline can run the deploy jobs of several sub-projects, in
// which case it produces one row per sub-project rather than a delimited value, so
// dashboards can `GROUP BY sub_project` without double counting.
//
// Rows are written by the monorepo plugin's attributeDeployments subtask. Projects
// without monorepo configuration have no rows here at all; dashboards should treat a
// missing mapping as the single, implicit "All" group via COALESCE(sub_project, 'All').
type CicdDeploymentSubproject struct {
	common.NoPKModel
	ProjectName string `gorm:"primaryKey;type:varchar(100)"`
	// CicdDeploymentId is the pipeline id, matching cicd_deployment_commits.cicd_deployment_id.
	// It also carries its own secondary index (idx_cds_deployment) because it sits in the
	// middle of the composite primary key, so dashboards filtering by deployment id alone
	// cannot use the primary key's leftmost prefix.
	CicdDeploymentId string `gorm:"primaryKey;type:varchar(255);index:idx_cds_deployment"`
	SubProject       string `gorm:"primaryKey;type:varchar(100)"`
}

func (CicdDeploymentSubproject) TableName() string {
	return "cicd_deployment_subprojects"
}
