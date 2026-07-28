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

package e2e

import (
	"sync"
	"testing"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/migration"
	coreMigration "github.com/apache/incubator-devlake/core/models/migrationscripts"
	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/apache/incubator-devlake/plugins/jira/impl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/schema"
)

// TestMigrationSchema guards against schema drift between Jira's migration
// scripts and its runtime GORM models.
//
// Regression test for the Sprint Report bug (introduced by PR #8967/#9010):
// the migration that created `_tool_jira_sprint_reports` used a struct that did
// NOT embed common.NoPKModel, so the `_raw_data_table` / `_raw_data_params` /
// `_raw_data_id` / `_raw_data_remark` columns were missing. The runtime model
// DID embed common.NoPKModel, so the ApiExtractor's cleanup query
// (`WHERE _raw_data_table = ? AND _raw_data_params = ?`) failed at runtime with
// "Error 1054: Unknown column '_raw_data_table' in 'where clause'".
//
// The test runs the REAL migration scripts (framework + jira) to build the
// schema exactly the way a production install would — deliberately NOT via
// AutoMigrate on the runtime model, which would silently hide such drift — and
// then asserts that every column each runtime model expects actually exists in
// the migrated table. Any future migration that forgets to embed
// common.NoPKModel (or otherwise omits a column) will fail this test.
//
// The migrations run against a dedicated, empty database (see
// e2ehelper.NewIsolatedMigrationDb) because the shared e2e database is polluted
// by the other e2e tests, which AutoMigrate tables without recording anything
// in `_devlake_migration_history`.
//
// Requires E2E_DB_URL (runs under `make e2e-test` / `make e2e-test-go-plugins`).
func TestMigrationSchema(t *testing.T) {
	var pluginInstance impl.Jira

	db := e2ehelper.NewIsolatedMigrationDb(t, "jira_migration_schema")
	dalInstance := dalgorm.NewDalgorm(db)

	// Apply the real migration scripts so the schema matches a production install.
	basicRes := runner.CreateBasicRes(config.GetConfig(), logruslog.Global, db)
	migrator, err := migration.NewMigrator(basicRes)
	require.NoError(t, err)
	migrator.Register(coreMigration.All(), "Framework")
	migrator.Register(pluginInstance.MigrationScripts(), "jira")
	require.NoError(t, migrator.Execute())

	keepAll := func(dal.ColumnMeta) bool { return true }

	for _, table := range pluginInstance.GetTablesInfo() {
		table := table
		t.Run(table.TableName(), func(t *testing.T) {
			// Columns the runtime GORM model expects.
			sch, err := schema.Parse(table, &sync.Map{}, schema.NamingStrategy{})
			require.NoErrorf(t, err, "unable to parse schema for %T", table)

			// Columns that actually exist in the migrated table.
			actualColumns, err := dal.GetColumnNames(dalInstance, table, keepAll)
			if err != nil || len(actualColumns) == 0 {
				// No migration creates this table (e.g. API response models
				// that are listed in GetTablesInfo but never persisted) —
				// there is no schema to drift from.
				t.Skipf("table %q not present after migrations", table.TableName())
			}
			existing := make(map[string]struct{}, len(actualColumns))
			for _, c := range actualColumns {
				existing[c] = struct{}{}
			}

			for _, field := range sch.Fields {
				if field.DBName == "" || field.IgnoreMigration {
					continue
				}
				_, ok := existing[field.DBName]
				assert.Truef(t, ok,
					"table %q is missing column %q expected by model %T — "+
						"did a migration script forget to embed common.NoPKModel (raw-data columns) or add the field?",
					table.TableName(), field.DBName, table)
			}
		})
	}
}
