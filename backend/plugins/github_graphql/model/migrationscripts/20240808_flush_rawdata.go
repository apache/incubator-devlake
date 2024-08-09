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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*flushRawData)(nil)

type flushRawData struct{}

func (*flushRawData) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	errors.Must(db.Delete("_devlake_collector_latest_state", dal.Where("raw_data_table like ?", "_raw_github_graphql_%")))
	_ = db.DropTables(
		"_raw_github_graphql_accounts",
		"_raw_github_graphql_deployment",
		"_raw_github_graphql_issues",
		"_raw_github_graphql_jobs",
		"_raw_github_graphql_prs",
		"_raw_github_graphql_release",
	)
	return nil
}

func (*flushRawData) Version() uint64 {
	return 20240808141953
}

func (*flushRawData) Name() string {
	return "flush github graphql raw data due to storaging granularity changed from per page to per item"
}
