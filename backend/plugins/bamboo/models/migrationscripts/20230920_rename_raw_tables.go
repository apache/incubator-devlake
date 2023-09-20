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

var _ plugin.MigrationScript = (*renameMultiBambooRawTables20230920)(nil)

type renameMultiBambooRawTables20230920 struct{}

func (*renameMultiBambooRawTables20230920) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	changedTables := map[string]string{
		"_raw_bamboo_api_deploy_build": "_raw_bamboo_api_deploy_builds",
		"_raw_bamboo_api_deploy":       "_raw_bamboo_api_deploys",
		"_raw_bamboo_api_job_build":    "_raw_bamboo_api_job_builds",
		"_raw_bamboo_api_job":          "_raw_bamboo_api_jobs",
		"_raw_bamboo_api_plan_build":   "_raw_bamboo_api_plan_builds",
	}
	for oldName, newName := range changedTables {
		if db.HasTable(oldName) {
			if err := db.RenameTable(oldName, newName); err != nil {
				return err
			}
		}
	}
	return nil
}

func (*renameMultiBambooRawTables20230920) Version() uint64 {
	return 20230920000001
}

func (*renameMultiBambooRawTables20230920) Name() string {
	return "rename _raw_bamboo_api_deploy_build to _raw_bamboo_api_deploy_builds," +
		" _raw_bamboo_api_deploy to _raw_bamboo_api_deploys," +
		" _raw_bamboo_api_job_build to _raw_bamboo_api_job_builds," +
		" _raw_bamboo_api_job to _raw_bamboo_api_jobs," +
		" _raw_bamboo_api_plan_build to _raw_bamboo_api_plan_builds"
}
