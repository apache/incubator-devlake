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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
)

type addProjectPrMetric struct{}

func (u *addProjectPrMetric) Up(baseRes core.BasicRes) errors.Error {
	db := baseRes.GetDal()
	err := migrationhelper.AutoMigrateTables(
		baseRes,
		&archived.ProjectPrMetrics{},
	)
	if err != nil {
		return err
	}
	prColums := []string{
		`coding_timespan`,
		`review_lag`,
		`review_timespan`,
		`deploy_timespan`,
		`change_timespan`,
		`orig_coding_timespan`,
		`orig_review_lag`,
		`orig_review_timespan`,
		`orig_deploy_timespan`,
	}
	err = db.DropColumns(`pull_requests`, prColums...)
	if err != nil {
		return err
	}
	err = db.DropColumns(`cicd_pipeline_commits`, "repo_url")
	return err
}

func (*addProjectPrMetric) Version() uint64 {
	return 20221111000001
}

func (*addProjectPrMetric) Name() string {
	return "add project pr metric tables"
}
