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

package runner

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormmigrator "gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

// primaryKeyIndexName is the fixed name MySQL gives to the primary key index
const primaryKeyIndexName = "PRIMARY"

// wrapMysqlDialector decorates the GORM MySQL dialector so that AutoMigrate keeps
// working for DevLake's migration scripts.
//
// Since gorm.io/driver/mysql v1.6.0, `AddColumn` appends `ADD PRIMARY KEY (<col>)`
// whenever the added field carries the `primaryKey` tag. DevLake migration scripts
// regularly add such a column to a table that already owns a primary key, which
// makes MySQL fail with `Error 1068: Multiple primary key defined`. Primary key
// changes are always performed explicitly by the migration scripts (see
// migrationhelper.ChangePrimaryKeyColumnsType), so the column is added without the
// implicit clause, matching the behaviour of the previous driver release. The only
// exception is an `AUTO_INCREMENT` column, which MySQL requires to be a key.
func wrapMysqlDialector(dialector gorm.Dialector) gorm.Dialector {
	if mysqlDialector, ok := dialector.(*mysql.Dialector); ok {
		return safeMysqlDialector{Dialector: *mysqlDialector}
	}
	return dialector
}

type safeMysqlDialector struct {
	mysql.Dialector
}

// Migrator returns the MySQL migrator with the AddColumn workaround applied
func (d safeMysqlDialector) Migrator(db *gorm.DB) gorm.Migrator {
	return safeMysqlMigrator{Migrator: d.Dialector.Migrator(db), db: db}
}

type safeMysqlMigrator struct {
	gorm.Migrator
	db *gorm.DB
}

// mysqlMigratorInternals exposes the helpers of the embedded GORM migrator that
// are needed to rebuild the `ALTER TABLE ... ADD ...` statement
type mysqlMigratorInternals interface {
	RunWithValue(value interface{}, fc func(*gorm.Statement) error) error
	FullDataTypeOf(field *schema.Field) clause.Expr
	CurrentTable(stmt *gorm.Statement) interface{}
}

// AddColumn adds a column without implicitly promoting it to the primary key of
// the table, see wrapMysqlDialector for the rationale
func (m safeMysqlMigrator) AddColumn(value interface{}, name string) error {
	base, ok := m.Migrator.(mysqlMigratorInternals)
	if !ok {
		return m.Migrator.AddColumn(value, name)
	}
	return base.RunWithValue(value, func(stmt *gorm.Statement) error {
		field := stmt.Schema.LookUpField(name)
		if field == nil {
			return fmt.Errorf("failed to look up field with name: %s", name)
		}
		if field.IgnoreMigration {
			return nil
		}
		fieldType := base.FullDataTypeOf(field)
		columnName := clause.Column{Name: field.DBName}
		values := []interface{}{base.CurrentTable(stmt), columnName, fieldType}
		alterSql := "ALTER TABLE ? ADD ? ?"
		// MySQL requires an AUTO_INCREMENT column to be a key, every other primary
		// key is managed explicitly by the migration scripts
		if strings.Contains(strings.ToLower(fieldType.SQL), "auto_increment") &&
			!m.Migrator.HasIndex(value, primaryKeyIndexName) {
			alterSql += ", ADD PRIMARY KEY (?)"
			values = append(values, columnName)
		}
		return m.db.Exec(alterSql, values...).Error
	})
}

// BuildIndexOptions forwards to the embedded migrator, GORM requires every
// migrator to implement migrator.BuildIndexOptionsInterface while creating tables
func (m safeMysqlMigrator) BuildIndexOptions(opts []schema.IndexOption, stmt *gorm.Statement) []interface{} {
	if builder, ok := m.Migrator.(gormmigrator.BuildIndexOptionsInterface); ok {
		return builder.BuildIndexOptions(opts, stmt)
	}
	return nil
}
