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

package schema_e2e

import (
	"fmt"
	"testing"

	"github.com/apache/devlake/core/config"
	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/models/migrationscripts/archived"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/core/runner"
	"github.com/apache/devlake/helpers/e2ehelper"
	"github.com/apache/devlake/impls/dalgorm"
	"github.com/apache/devlake/impls/logruslog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	copilotimpl "github.com/apache/devlake/plugins/gh-copilot/impl"
	jiraimpl "github.com/apache/devlake/plugins/jira/impl"
	taigaimpl "github.com/apache/devlake/plugins/taiga/impl"
	teambitionimpl "github.com/apache/devlake/plugins/teambition/impl"
	testmoimpl "github.com/apache/devlake/plugins/testmo/impl"
)

// ----------------------------------------------------------------------------
// Pre-repair table shapes
//
// Each struct reproduces a table EXACTLY as the buggy migration left it, i.e.
// without the columns that the repair migration adds. The repair script is then
// executed against that table *with rows in it*.
// ----------------------------------------------------------------------------

// upgradePreJiraSprintReport is _tool_jira_sprint_reports as created by
// 20260722 — no embedded NoPKModel, hence no _raw_data_* / created_at /
// updated_at columns.
type upgradePreJiraSprintReport struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	BoardId      uint64 `gorm:"primaryKey"`
	SprintId     uint64 `gorm:"primaryKey"`
	IssueId      uint64 `gorm:"primaryKey"`

	IssueKey                 string `gorm:"type:varchar(255)"`
	Bucket                   string `gorm:"type:varchar(32);index"`
	Done                     bool
	StoryPointsAtSprintStart *float64
	StoryPointsAtSprintEnd   *float64
}

func (upgradePreJiraSprintReport) TableName() string { return "_tool_jira_sprint_reports" }

// upgradePreTaigaScopeConfig lacks `type_mappings`.
type upgradePreTaigaScopeConfig struct {
	archived.Model
	Entities     []string `gorm:"type:json;serializer:json"`
	ConnectionId uint64   `gorm:"index"`
	Name         string   `gorm:"type:varchar(255);uniqueIndex"`
}

func (upgradePreTaigaScopeConfig) TableName() string { return "_tool_taiga_scope_configs" }

// upgradePreTeambitionScopeConfig lacks the embedded common.Model, i.e. the
// table has no primary key at all and no id / created_at / updated_at.
type upgradePreTeambitionScopeConfig struct {
	Entities          []string          `gorm:"type:json;serializer:json"`
	ConnectionId      uint64            `gorm:"index"`
	Name              string            `gorm:"type:varchar(255)"`
	TypeMappings      map[string]string `gorm:"serializer:json"`
	StatusMappings    map[string]string `gorm:"serializer:json"`
	BugDueDateField   string            `gorm:"column:bug_due_date_field"`
	TaskDueDateField  string            `gorm:"column:task_due_date_field"`
	StoryDueDateField string            `gorm:"column:story_due_date_field"`
}

func (upgradePreTeambitionScopeConfig) TableName() string {
	return "_tool_teambition_scope_configs"
}

// upgradePreTestmoScopeConfig lacks `connection_id` and `name`.
type upgradePreTestmoScopeConfig struct {
	archived.Model
	Entities              []string `gorm:"type:json;serializer:json"`
	AcceptanceTestPattern string   `gorm:"type:varchar(255)"`
	SmokeTestPattern      string   `gorm:"type:varchar(255)"`
	TeamPattern           string   `gorm:"type:varchar(255)"`
}

func (upgradePreTestmoScopeConfig) TableName() string { return "_tool_testmo_scope_configs" }

// upgradePreCopilotEnterpriseCredits and friends lack the seven credit
// breakdown columns, which 20260708 declared through an unexported embedded
// struct that GORM silently ignores.
type upgradePreCopilotEnterpriseCredits struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(191)"`
	Year         int    `gorm:"primaryKey"`
	Month        int    `gorm:"primaryKey"`
	Day          int    `gorm:"primaryKey"`
	Enterprise   string `gorm:"primaryKey;type:varchar(191)"`
	Model        string `gorm:"primaryKey;type:varchar(191)"`
	Product      string `gorm:"type:varchar(32)"`
	archived.NoPKModel
}

func (upgradePreCopilotEnterpriseCredits) TableName() string {
	return "_tool_copilot_enterprise_ai_credit_usage"
}

type upgradePreCopilotOrgCredits struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(191)"`
	Year         int    `gorm:"primaryKey"`
	Month        int    `gorm:"primaryKey"`
	Day          int    `gorm:"primaryKey"`
	Organization string `gorm:"primaryKey;type:varchar(191)"`
	Model        string `gorm:"primaryKey;type:varchar(191)"`
	Product      string `gorm:"type:varchar(32)"`
	archived.NoPKModel
}

func (upgradePreCopilotOrgCredits) TableName() string {
	return "_tool_copilot_org_ai_credit_usage"
}

type upgradePreCopilotUserCredits struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(191)"`
	Year         int    `gorm:"primaryKey"`
	Month        int    `gorm:"primaryKey"`
	Day          int    `gorm:"primaryKey"`
	User         string `gorm:"primaryKey;type:varchar(191)"`
	Model        string `gorm:"primaryKey;type:varchar(191)"`
	Product      string `gorm:"type:varchar(32)"`
	archived.NoPKModel
}

func (upgradePreCopilotUserCredits) TableName() string {
	return "_tool_copilot_user_ai_credit_usage"
}

var copilotCreditColumns = []string{
	"gross_quantity", "discount_quantity", "net_quantity", "price_per_unit",
	"gross_amount", "discount_amount", "net_amount",
}

// upgradeCase describes one repair migration and how to exercise it on a table
// that already contains data.
type upgradeCase struct {
	// plugin owning the migration script, used to look the script up by version
	// instead of duplicating it here (so the test breaks if the script is
	// removed or renumbered).
	plugin  plugin.PluginMeta
	version uint64
	// pre is the table as the buggy migration left it.
	pre dal.Tabler
	// seed rows inserted BEFORE the repair migration runs.
	seed []map[string]interface{}
	// wantColumns must exist after the repair.
	wantColumns []string
	// wantPrimaryKey asserts the table has a primary key afterwards. A plain
	// AutoMigrate cannot add one, which is invisible to a column-only check
	// (and silently accepted by PostgreSQL).
	wantPrimaryKey bool
	// autoIncColumn, if set, must be backfilled with distinct non-zero values
	// for the pre-existing rows, and a subsequent INSERT must still work.
	autoIncColumn string
}

// TestMigrationUpgradePathOnPopulatedTables complements
// TestMigrationSchemaMatchesModels: that guard proves the END STATE of a fresh
// migration run matches the models, but every table it inspects is empty, so it
// cannot exercise the upgrade path of a repair migration on a database that
// already holds rows — which is the only situation those migrations exist for.
//
// For each repair migration this test therefore
//  1. recreates the table exactly as the buggy migration left it,
//  2. inserts rows,
//  3. runs ONLY that repair script,
//  4. asserts the columns were added, the rows survived untouched, the primary
//     key exists and auto-increment ids were backfilled.
//
// Step 4 is what a column-presence check on an empty table cannot see: replacing
// the explicit AUTO_INCREMENT DDL in the teambition script with a plain
// AutoMigrate is accepted by PostgreSQL (it adds `bigserial` without a primary
// key), and only this test notices.
func TestMigrationUpgradePathOnPopulatedTables(t *testing.T) {
	db := e2ehelper.NewIsolatedMigrationDb(t, "upgrade_path")
	dalInstance := dalgorm.NewDalgorm(db)
	basicRes := runner.CreateBasicRes(config.GetConfig(), logruslog.Global, db)

	cases := map[string]upgradeCase{
		"jira sprint report raw data columns": {
			plugin:  jiraimpl.Jira{},
			version: 20260727000000,
			pre:     upgradePreJiraSprintReport{},
			seed: []map[string]interface{}{
				{"connection_id": 1, "board_id": 10, "sprint_id": 100, "issue_id": 1000, "issue_key": "TEST-1", "bucket": "committed", "done": false},
				{"connection_id": 1, "board_id": 10, "sprint_id": 100, "issue_id": 1001, "issue_key": "TEST-2", "bucket": "completed", "done": true},
			},
			wantColumns: []string{"_raw_data_params", "_raw_data_table", "_raw_data_id", "_raw_data_remark", "created_at", "updated_at"},
		},
		"taiga scope config type_mappings": {
			plugin:  taigaimpl.Taiga{},
			version: 20260727000001,
			pre:     upgradePreTaigaScopeConfig{},
			seed: []map[string]interface{}{
				{"id": 1, "connection_id": 1, "name": "taiga-cfg-a"},
				{"id": 2, "connection_id": 1, "name": "taiga-cfg-b"},
			},
			wantColumns:    []string{"type_mappings"},
			wantPrimaryKey: true,
		},
		"teambition scope config primary key": {
			plugin:  teambitionimpl.Teambition{},
			version: 20260727000001,
			pre:     upgradePreTeambitionScopeConfig{},
			seed: []map[string]interface{}{
				{"connection_id": 1, "name": "teambition-cfg-a"},
				{"connection_id": 1, "name": "teambition-cfg-b"},
			},
			wantColumns:    []string{"id", "created_at", "updated_at"},
			wantPrimaryKey: true,
			autoIncColumn:  "id",
		},
		"testmo scope config connection_id/name": {
			plugin:  testmoimpl.Testmo{},
			version: 20260727000001,
			pre:     upgradePreTestmoScopeConfig{},
			seed: []map[string]interface{}{
				// `name` carries a uniqueIndex in the repaired shape; both rows
				// are backfilled with NULL, which MySQL and PostgreSQL accept.
				{"id": 1, "acceptance_test_pattern": "a"},
				{"id": 2, "acceptance_test_pattern": "b"},
			},
			wantColumns:    []string{"connection_id", "name"},
			wantPrimaryKey: true,
		},
		"gh-copilot enterprise credit breakdown": {
			plugin:  copilotimpl.GhCopilot{},
			version: 20260731000000,
			pre:     upgradePreCopilotEnterpriseCredits{},
			seed: []map[string]interface{}{
				{"connection_id": 1, "scope_id": "ent-1", "year": 2026, "month": 8, "day": 1, "enterprise": "acme", "model": "gpt-4.1", "product": "copilot"},
			},
			wantColumns: copilotCreditColumns,
		},
		"gh-copilot org credit breakdown": {
			plugin:  copilotimpl.GhCopilot{},
			version: 20260731000000,
			pre:     upgradePreCopilotOrgCredits{},
			seed: []map[string]interface{}{
				{"connection_id": 1, "scope_id": "org-1", "year": 2026, "month": 8, "day": 1, "organization": "acme", "model": "gpt-4.1", "product": "copilot"},
			},
			wantColumns: copilotCreditColumns,
		},
		"gh-copilot user credit breakdown": {
			plugin:  copilotimpl.GhCopilot{},
			version: 20260731000000,
			pre:     upgradePreCopilotUserCredits{},
			seed: []map[string]interface{}{
				{"connection_id": 1, "scope_id": "user-1", "year": 2026, "month": 8, "day": 1, "user": "octocat", "model": "gpt-4.1", "product": "copilot"},
			},
			wantColumns: copilotCreditColumns,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			table := c.pre.TableName()
			script := findMigrationScript(t, c.plugin, c.version)

			// 1. table exactly as the buggy migration left it
			require.NoError(t, db.Migrator().DropTable(table))
			require.NoError(t, db.Table(table).AutoMigrate(c.pre))
			for _, column := range c.wantColumns {
				require.Falsef(t, db.Migrator().HasColumn(c.pre, column),
					"precondition failed: %q already has column %q, the pre-repair shape is wrong",
					table, column)
			}

			// 2. rows, so the migration has to survive real data
			for _, row := range c.seed {
				require.NoError(t, db.Table(table).Create(row).Error)
			}

			// 3. run ONLY the repair script
			require.NoErrorf(t, script.Up(basicRes),
				"migration %q failed on a populated %q", script.Name(), table)

			// 4a. the columns are there now
			actual, colErr := dal.GetColumnNames(dalInstance, dal.DefaultTabler{Name: table},
				func(dal.ColumnMeta) bool { return true })
			require.NoError(t, colErr)
			existing := make(map[string]struct{}, len(actual))
			for _, column := range actual {
				existing[column] = struct{}{}
			}
			for _, column := range c.wantColumns {
				_, ok := existing[column]
				assert.Truef(t, ok, "table %q is still missing column %q after %q",
					table, column, script.Name())
			}

			// 4b. no data was lost
			var rows int64
			require.NoError(t, db.Table(table).Count(&rows).Error)
			assert.EqualValuesf(t, len(c.seed), rows,
				"migration %q changed the row count of %q", script.Name(), table)

			// 4c. the primary key survived / was created
			if c.wantPrimaryKey {
				assert.Truef(t, hasPrimaryKey(t, db, table),
					"table %q has no primary key after %q — AutoMigrate cannot add one, "+
						"the script needs explicit DDL", table, script.Name())
			}

			// 4d. auto-increment ids were backfilled and the counter still works
			if c.autoIncColumn != "" {
				var ids []uint64
				require.NoError(t, db.Table(table).Pluck(c.autoIncColumn, &ids).Error)
				require.Len(t, ids, len(c.seed))
				seen := map[uint64]struct{}{}
				for _, id := range ids {
					assert.NotZerof(t, id, "pre-existing row was not assigned a %q", c.autoIncColumn)
					_, dup := seen[id]
					assert.Falsef(t, dup, "duplicate %q=%d after backfill", c.autoIncColumn, id)
					seen[id] = struct{}{}
				}
				require.NoErrorf(t, db.Table(table).Create(map[string]interface{}{
					"connection_id": 2, "name": "inserted-after-migration",
				}).Error, "INSERT after the migration failed, the sequence/counter is out of sync")
			}
		})
	}
}

// findMigrationScript looks the script up through the plugin's own
// MigrationScripts() so the test fails if it is removed or renumbered.
func findMigrationScript(t *testing.T, p plugin.PluginMeta, version uint64) plugin.MigrationScript {
	migratable, ok := p.(plugin.PluginMigration)
	require.Truef(t, ok, "plugin %s does not implement PluginMigration", p.Name())
	for _, script := range migratable.MigrationScripts() {
		if script.Version() == version {
			return script
		}
	}
	t.Fatalf("plugin %s has no migration script with version %d", p.Name(), version)
	return nil
}

// hasPrimaryKey works on MySQL and PostgreSQL alike.
func hasPrimaryKey(t *testing.T, db *gorm.DB, table string) bool {
	schemaFunc := "current_schema()"
	if db.Dialector.Name() == "mysql" {
		schemaFunc = "DATABASE()"
	}
	var count int64
	err := db.Raw(fmt.Sprintf(
		`SELECT COUNT(*) FROM information_schema.table_constraints
		 WHERE constraint_type = 'PRIMARY KEY' AND table_name = ? AND table_schema = %s`,
		schemaFunc), table).Scan(&count).Error
	require.NoError(t, err)
	return count > 0
}
