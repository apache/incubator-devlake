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
)

type modifyFileMetaTable struct{}

func (*modifyFileMetaTable) Name() string {
	return "Modify QDevS3FileMeta table to allow NULL processed_time"
}

func (*modifyFileMetaTable) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	// Target table and column
	tableName := "_tool_q_dev_s3_file_meta"
	columnName := "processed_time"

	// If column doesn't exist, no migration needed, idempotent
	if !db.HasColumn(tableName, columnName) {
		return nil
	}

	// Read column metadata to check if already nullable, return idempotently if already nullable
	var processedTimeNullable bool
	{
		cols, err := db.GetColumns(dal.DefaultTabler{Name: tableName}, func(cm dal.ColumnMeta) bool {
			return cm.Name() == columnName
		})
		if err != nil {
			return errors.Default.Wrap(err, "failed to load column metadata for _tool_q_dev_s3_file_meta.processed_time")
		}
		if len(cols) == 0 {
			// If column is not visible in metadata, treat as no processing needed
			return nil
		}
		if nullable, ok := cols[0].Nullable(); ok {
			processedTimeNullable = nullable
		}
	}
	if processedTimeNullable {
		return nil
	}

	// Execute compatible SQL by dialect
	switch db.Dialect() {
	case "postgres":
		// PostgreSQL makes column nullable via DROP NOT NULL, without changing data type
		if err := db.Exec(
			"ALTER TABLE ? ALTER COLUMN ? DROP NOT NULL",
			dal.ClauseTable{Name: tableName},
			dal.ClauseColumn{Name: columnName},
		); err != nil {
			return errors.Default.Wrap(err, "failed to drop NOT NULL on processed_time for postgres")
		}
		return nil
	case "mysql":
		// MySQL requires MODIFY COLUMN with original type specification, preserve original type as much as possible
		cols, err := db.GetColumns(dal.DefaultTabler{Name: tableName}, func(cm dal.ColumnMeta) bool {
			return cm.Name() == columnName
		})
		if err != nil {
			return errors.Default.Wrap(err, "failed to load column metadata for mysql type preservation")
		}
		columnTypeSql := "DATETIME"
		if len(cols) > 0 {
			if ct, ok := cols[0].ColumnType(); ok && ct != "" {
				columnTypeSql = ct
			} else if dbt := cols[0].DatabaseTypeName(); dbt != "" {
				// DatabaseTypeName may return DATETIME, TIMESTAMP etc
				columnTypeSql = dbt
			}
		}
		alterSql := "ALTER TABLE ? MODIFY COLUMN ? " + columnTypeSql + " NULL"
		if err := db.Exec(
			alterSql,
			dal.ClauseTable{Name: tableName},
			dal.ClauseColumn{Name: columnName},
		); err != nil {
			return errors.Default.Wrap(err, "failed to modify processed_time to NULL for mysql")
		}
		return nil
	default:
		// Other dialects are not forced to migrate for now, return idempotently
		return nil
	}
}

func (*modifyFileMetaTable) Version() uint64 {
	return 20250320
}
