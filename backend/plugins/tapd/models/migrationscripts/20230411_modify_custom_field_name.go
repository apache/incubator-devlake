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
	"fmt"
	"strings"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"gorm.io/gorm/clause"
)

type modifyCustomFieldName struct{}

type columnRename struct {
	Old string
	New string
}

func (*modifyCustomFieldName) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	// Define all column renames for each table
	tableRenames := map[string][]columnRename{
		"_tool_tapd_bugs": {
			{Old: "custom_field6", New: "custom_field_6"},
			{Old: "custom_field7", New: "custom_field_7"},
			{Old: "custom_field8", New: "custom_field_8"},
		},
		"_tool_tapd_stories": {
			{Old: "custom_field6", New: "custom_field_six"},
			{Old: "custom_field7", New: "custom_field_seven"},
			{Old: "custom_field8", New: "custom_field_eight"},
		},
		"_tool_tapd_tasks": {
			{Old: "custom_field6", New: "custom_field_six"},
			{Old: "custom_field7", New: "custom_field_seven"},
			{Old: "custom_field8", New: "custom_field_eight"},
		},
	}

	// Add custom_field_9 to custom_field_50 for all tables
	for i := 9; i <= 50; i++ {
		oldCol := fmt.Sprintf("custom_field%d", i)
		newCol := fmt.Sprintf("custom_field_%d", i)
		for _, table := range []string{"_tool_tapd_bugs", "_tool_tapd_stories", "_tool_tapd_tasks"} {
			tableRenames[table] = append(tableRenames[table], columnRename{Old: oldCol, New: newCol})
		}
	}

	// Execute batch rename for each table
	for tableName, renames := range tableRenames {
		if err := batchRenameColumns(db, tableName, renames); err != nil {
			return err
		}
	}

	return nil
}

func (*modifyCustomFieldName) Version() uint64 {
	return 20230411000004
}

func (*modifyCustomFieldName) Name() string {
	return "modify tapd custom field name"
}

// batchRenameColumns renames multiple columns in a single ALTER TABLE statement
func batchRenameColumns(db dal.Dal, table string, renames []columnRename) errors.Error {
	if len(renames) == 0 {
		return nil
	}

	// Get all existing column names once to avoid repeated information_schema queries
	existingColumns := getExistingColumns(db, table)

	// Filter out renames where old column doesn't exist or new column already exists
	validRenames := filterValidRenamesCached(renames, existingColumns)
	if len(validRenames) == 0 {
		return nil
	}

	var sql string
	var dialect = db.Dialect()

	if dialect == "postgres" {
		// PostgreSQL requires separate ALTER TABLE statements for each RENAME COLUMN
		for _, rename := range validRenames {
			sql = fmt.Sprintf(`ALTER TABLE "%s" RENAME COLUMN "%s" TO "%s"`, table, rename.Old, rename.New)
			if err := db.Exec(sql); err != nil {
				return err
			}
		}
		// Clear PostgreSQL cached plan after all renames
		_ = db.Exec("SELECT * FROM ? LIMIT 1", clause.Table{Name: table})
	} else {
		// MySQL: ALTER TABLE t CHANGE COLUMN a new_name TEXT, CHANGE COLUMN c new_name2 TEXT
		clauses := make([]string, 0, len(validRenames))
		for _, rename := range validRenames {
			clauses = append(clauses, fmt.Sprintf("CHANGE COLUMN `%s` `%s` %s", rename.Old, rename.New, dal.Text.String()))
		}
		sql = fmt.Sprintf("ALTER TABLE `%s` %s", table, strings.Join(clauses, ", "))
		if err := db.Exec(sql); err != nil {
			return err
		}
	}

	return nil
}

// getExistingColumns fetches all column names for a table in a single query
func getExistingColumns(db dal.Dal, table string) map[string]bool {
	columns := make(map[string]bool)
	columnMetas, err := db.GetColumns(&dal.DefaultTabler{Name: table}, nil)
	if err != nil {
		return columns
	}
	for _, col := range columnMetas {
		columns[col.Name()] = true
	}
	return columns
}

// filterValidRenamesCached checks which renames are needed using pre-fetched column map
func filterValidRenamesCached(renames []columnRename, existingColumns map[string]bool) []columnRename {
	valid := make([]columnRename, 0, len(renames))
	for _, rename := range renames {
		oldExists := existingColumns[rename.Old]
		newExists := existingColumns[rename.New]
		if oldExists && !newExists {
			valid = append(valid, rename)
		}
	}
	return valid
}
