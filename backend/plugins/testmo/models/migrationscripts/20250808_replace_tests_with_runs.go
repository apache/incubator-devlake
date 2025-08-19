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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/testmo/models/migrationscripts/archived"
)

type replaceTestsWithRuns struct{}

func (*replaceTestsWithRuns) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	// Check if the old tool layer tests table exists and drop it if it does
	if db.HasTable("_tool_testmo_tests") {
		err := db.DropTables("_tool_testmo_tests")
		if err != nil {
			return err
		}
	}

	// Check if the old raw tests table exists and drop it if it does
	if db.HasTable("_raw_testmo_tests") {
		err := db.DropTables("_raw_testmo_tests")
		if err != nil {
			return err
		}
	}

	// Create the new runs table
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.TestmoRun{},
	)
}

func (*replaceTestsWithRuns) Version() uint64 {
	return 20250808000001
}

func (*replaceTestsWithRuns) Name() string {
	return "Replace testmo tests tables with runs table (both tool and raw layers)"
}

var _ plugin.MigrationScript = (*replaceTestsWithRuns)(nil)
