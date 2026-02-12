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

type fixColumnNames struct{}

type columnRename struct {
	Old string
	New string
}

func (*fixColumnNames) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	// Rename columns for _tool_circleci_pipelines
	renames := []columnRename{
		{Old: "created_date", New: "created_at"},
		{Old: "updated_date", New: "updated_at"},
	}

	return batchRenameColumns(db, "_tool_circleci_pipelines", renames)
}

func (*fixColumnNames) Version() uint64 {
	return 20260208000001
}

func (*fixColumnNames) Name() string {
	return "fix circleci column names"
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
			clauses = append(clauses, fmt.Sprintf("CHANGE COLUMN `%s` `%s` %s", rename.Old, rename.New, "DATETIME"))
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
