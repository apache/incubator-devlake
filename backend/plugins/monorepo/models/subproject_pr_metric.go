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

// SubProjectPrMetric holds the change-lead-time breakdown for a merged pull request,
// attributed to exactly one sub-project of a monorepo.
//
// Deprecated: this table is kept populated for one release for backward compatibility
// with dashboards/integrations built against it, but new dashboards should read
// project_pr_metrics.sub_project (joined through cicd_deployment_commits /
// cicd_deployment_subprojects for deployment-side data) instead. It is a candidate for
// removal in a follow-up release once the compat window closes.
//
// All five metric fields (CodingTime/PickupTime/ReviewTime/DeployTime/CycleTime) are now
// written by project_pr_metrics_updater.go's updateProjectPrMetricsSubProject subtask,
// copied verbatim from DORA's project_pr_metrics rather than recomputed here. In
// particular, DeployTime/CycleTime used to be computed by AttributePullRequests using a
// merge-date-nearest-deployment heuristic; that heuristic has been retired because it is
// less accurate than DORA's own commit-based PR-to-deployment attribution
// (project_pr_metrics.deployment_commit_id). Existing monorepo users will see these two
// values change (improve) on upgrade - this is a correction, not a regression.
//
// All durations are in minutes, matching DORA's convention.
type SubProjectPrMetric struct {
	common.NoPKModel
	ProjectName   string `gorm:"primaryKey;type:varchar(100)"`
	PullRequestId string `gorm:"primaryKey;type:varchar(255)"`
	SubProject    string `gorm:"index;type:varchar(255)"`

	CodingTime *int64
	PickupTime *int64
	ReviewTime *int64
	DeployTime *int64
	CycleTime  *int64

	// DeploymentId is the cicd_deployment_id (pipeline id) of the deployment DORA
	// attributed this PR to via project_pr_metrics.deployment_commit_id, if any.
	DeploymentId string `gorm:"type:varchar(255)"`

	PrCreatedDate *time.Time
	PrMergedDate  *time.Time
	DeployedDate  *time.Time
}

func (SubProjectPrMetric) TableName() string {
	return "monorepo_subproject_pr_metrics"
}
