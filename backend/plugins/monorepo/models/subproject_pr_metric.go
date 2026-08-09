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
// CodingTime/PickupTime/ReviewTime are carried over from DORA's project_pr_metrics:
// they depend only on the pull request itself, so DORA already computes them correctly
// for a monorepo. Only DeployTime (and therefore CycleTime) is recomputed here, against
// the deployments of this sub-project rather than the whole repository's.
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

	// DeploymentId is the sub-project deployment this PR was linked to, if any.
	DeploymentId string `gorm:"type:varchar(255)"`

	PrCreatedDate *time.Time
	PrMergedDate  *time.Time
	DeployedDate  *time.Time
}

func (SubProjectPrMetric) TableName() string {
	return "monorepo_subproject_pr_metrics"
}
