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

package archived

type ProjectPrMetric struct {
	DomainEntity
	ProjectName    string `gorm:"primaryKey;type:varchar(100)"`
	FirstReviewId  string
	FirstCommitSha string
	CodingTimespan *int64
	ReviewLag      *int64
	ReviewTimespan *int64
	DeploymentId   string
	DeployTimespan *int64
	ChangeTimespan *int64
}

func (ProjectPrMetric) TableName() string {
	return "project_pr_metrics"
}
