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

type fixRawTableNamesAndSchemas struct{}

func (*fixRawTableNamesAndSchemas) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	legacy := []string{
		"testmo_projects",
		"testmo_milestones",
		"testmo_automation_runs",
		"testmo_runs",
	}
	for _, tbl := range legacy {
		if db.HasTable(tbl) {
			if err := db.DropTables(tbl); err != nil {
				return err
			}
		}
	}

	current := []string{
		"_raw_testmo_projects",
		"_raw_testmo_milestones",
		"_raw_testmo_automation_runs",
		"_raw_testmo_runs",
	}
	for _, tbl := range current {
		if db.HasTable(tbl) {
			err := db.All(&[]struct{}{}, dal.From(tbl), dal.Where("_raw_data_table = '' AND 1 = 0"))
			if err != nil {
				if dropErr := db.DropTables(tbl); dropErr != nil {
					return dropErr
				}
			}
		}
	}

	return nil
}

func (*fixRawTableNamesAndSchemas) Version() uint64 { return 20250905000001 }
func (*fixRawTableNamesAndSchemas) Name() string    { return "Fix raw table names and schemas" }

var _ plugin.MigrationScript = (*fixRawTableNamesAndSchemas)(nil)
