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
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addSubProjectToPrAndMetrics)(nil)

// pullRequest20260810 adds the sub_project column that the monorepo plugin's
// attributePullRequests subtask writes. Nullable, so single-repo projects (and rows not
// yet processed) simply read as NULL and dashboards fall back to the "All" group.
type pullRequest20260810 struct {
	SubProject string `gorm:"index;type:varchar(100)"`
}

func (pullRequest20260810) TableName() string {
	return "pull_requests"
}

// pullRequestCommit20260810 mirrors the owning pull request's sub_project.
type pullRequestCommit20260810 struct {
	SubProject string `gorm:"index;type:varchar(100)"`
}

func (pullRequestCommit20260810) TableName() string {
	return "pull_request_commits"
}

// projectPrMetric20260810 is tagged by the monorepo plugin's
// updateProjectPrMetricsSubProject subtask after DORA computes this row.
type projectPrMetric20260810 struct {
	SubProject string `gorm:"index;type:varchar(100)"`
}

func (projectPrMetric20260810) TableName() string {
	return "project_pr_metrics"
}

type addSubProjectToPrAndMetrics struct{}

func (script *addSubProjectToPrAndMetrics) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.AutoMigrate(&pullRequest20260810{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&pullRequestCommit20260810{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&projectPrMetric20260810{}); err != nil {
		return err
	}
	return nil
}

func (*addSubProjectToPrAndMetrics) Version() uint64 {
	return 20260810100000
}

func (*addSubProjectToPrAndMetrics) Name() string {
	return "add sub_project to pull_requests, pull_request_commits and project_pr_metrics for monorepo support"
}
