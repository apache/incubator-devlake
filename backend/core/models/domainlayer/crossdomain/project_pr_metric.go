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

package crossdomain

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/domainlayer"
)

type ProjectPrMetric struct {
	domainlayer.DomainEntity
	ProjectName        string `gorm:"primaryKey;type:varchar(100)"`
	FirstCommitSha     string
	PrCodingTime       *int64
	FirstReviewId      string
	PrPickupTime       *int64
	PrReviewTime       *int64
	DeploymentCommitId string
	PrDeployTime       *int64
	PrCycleTime        *int64

	FirstCommitAuthoredDate *time.Time
	FirstCommentDate        *time.Time
	PrCreatedDate           *time.Time
	PrMergedDate            *time.Time
	PrDeployedDate          *time.Time
}

func (ProjectPrMetric) TableName() string {
	return "project_pr_metrics"
}
