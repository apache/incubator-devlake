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

package dalgorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Dalgorm FIXME ...
type Dalgorm struct {
	db *gorm.DB
}

func transformParams(params []interface{}) []interface{} {
	tp := make([]interface{}, 0, len(params))

	for _, v := range params {
		switch p := v.(type) {
		case dal.ClauseColumn:
			tp = append(tp, clause.Column{
				Table: p.Table,
				Name:  p.Name,
				Alias: p.Alias,
				Raw:   p.Raw,
			})
		case dal.ClauseTable:
			tp = append(tp, clause.Table{
				Name:  p.Name,
				Alias: p.Alias,
				Raw:   p.Raw,
			})
		default:
			tp = append(tp, p)
		}
	}

	return tp
}

func buildTx(tx *gorm.DB, clauses []dal.Clause) *gorm.DB {
	for _, c := range clauses {
		t := c.Type
		d := c.Data
		switch t {
		case dal.JoinClause:
			tx = tx.Joins(d.(dal.DalClause).Expr, transformParams(d.(dal.DalClause).Params)...)
		case dal.WhereClause:
			tx = tx.Where(d.(dal.DalClause).Expr, transformParams(d.(dal.DalClause).Params)...)
		case dal.OrderbyClause:
			tx = tx.Order(d.(string))
		case dal.LimitClause:
			tx = tx.Limit(d.(int))
		case dal.OffsetClause:
			tx = tx.Offset(d.(int))
		case dal.FromClause:
			switch dd := d.(type) {
			case string:
				tx = tx.Table(dd)
			case dal.DalClause:
				tx = tx.Table(dd.Expr, transformParams(dd.Params)...)
			case dal.ClauseTable:
				tx = tx.Table(" ? ", clause.Table{
					Name:  dd.Name,
					Alias: dd.Alias,
					Raw:   dd.Raw,
				})
			default:
				tx = tx.Model(d)
			}
		case dal.SelectClause:
			tx = tx.Select(d.(dal.DalClause).Expr, transformParams(d.(dal.DalClause).Params)...)
		case dal.GroupbyClause:
			tx = tx.Group(d.(string))
		case dal.HavingClause:
			tx = tx.Having(d.(dal.DalClause).Expr, transformParams(d.(dal.DalClause).Params)...)
		}
	}
	return tx
}

var _ dal.Dal = (*Dalgorm)(nil)

// Exec executes raw sql query
func (d *Dalgorm) Exec(query string, params ...interface{}) errors.Error {
	return errors.Convert(d.db.Exec(query, transformParams(params)...).Error)
}

// AutoMigrate runs auto migration for given models
func (d *Dalgorm) AutoMigrate(entity interface{}, clauses ...dal.Clause) errors.Error {
	err := errors.Convert(buildTx(d.db, clauses).AutoMigrate(entity))
	if err == nil {
		// fix pg cache plan error
		_ = d.First(entity, clauses...)
	}
	return err
}

// Cursor returns a database cursor, cursor is especially useful when handling big amount of rows of data
func (d *Dalgorm) Cursor(clauses ...dal.Clause) (dal.Rows, errors.Error) {
	return errors.Convert01(buildTx(d.db, clauses).Rows())
}

// CursorTx FIXME ...
func (d *Dalgorm) CursorTx(clauses ...dal.Clause) *gorm.DB {
	return buildTx(d.db, clauses)
}

// Fetch loads row data from `cursor` into `dst`
func (d *Dalgorm) Fetch(cursor dal.Rows, dst interface{}) errors.Error {
	if rows, ok := cursor.(*sql.Rows); ok {
		return errors.Convert(d.db.ScanRows(rows, dst))
	} else {
		return errors.Default.New(fmt.Sprintf("can not support type %s to be a dal.Rows interface", reflect.TypeOf(cursor).String()))
	}
}

// All loads matched rows from database to `dst`, USE IT WITH COUTIOUS!!
func (d *Dalgorm) All(dst interface{}, clauses ...dal.Clause) errors.Error {
	return errors.Convert(buildTx(d.db, clauses).Find(dst).Error)
}

// First loads the first matched row from database to `dst`, error will be returned if no records were found
func (d *Dalgorm) First(dst interface{}, clauses ...dal.Clause) errors.Error {
	err := buildTx(d.db, clauses).First(dst).Error
	if err == gorm.ErrRecordNotFound {
		err = errors.ErrRecordNotFound
	}
	return errors.Convert(err)
}

// Count total records
func (d *Dalgorm) Count(clauses ...dal.Clause) (int64, errors.Error) {
	var count int64
	err := buildTx(d.db, clauses).Count(&count).Error
	return errors.Convert01(count, err)
}

// Pluck used to query single column
func (d *Dalgorm) Pluck(column string, dest interface{}, clauses ...dal.Clause) errors.Error {
	return errors.Convert(buildTx(d.db, clauses).Pluck(column, dest).Error)
}

// Create insert record to database
func (d *Dalgorm) Create(entity interface{}, clauses ...dal.Clause) errors.Error {
	return errors.Convert(buildTx(d.db, clauses).Create(entity).Error)
}

// Update updates record
func (d *Dalgorm) Update(entity interface{}, clauses ...dal.Clause) errors.Error {
	return errors.Convert(buildTx(d.db, clauses).Save(entity).Error)
}

// CreateOrUpdate tries to create the record, or fallback to update all if failed
func (d *Dalgorm) CreateOrUpdate(entity interface{}, clauses ...dal.Clause) errors.Error {
	return errors.Convert(buildTx(d.db, clauses).Clauses(clause.OnConflict{UpdateAll: true}).Create(entity).Error)
}

// CreateIfNotExist tries to create the record if not exist
func (d *Dalgorm) CreateIfNotExist(entity interface{}, clauses ...dal.Clause) errors.Error {
	return errors.Convert(buildTx(d.db, clauses).Clauses(clause.OnConflict{DoNothing: true}).Create(entity).Error)
}

// Delete records from database
func (d *Dalgorm) Delete(entity interface{}, clauses ...dal.Clause) errors.Error {
	return errors.Convert(buildTx(d.db, clauses).Delete(entity).Error)
}

// UpdateColumn allows you to update mulitple records
func (d *Dalgorm) UpdateColumn(entity interface{}, columnName string, value interface{}, clauses ...dal.Clause) errors.Error {
	if expr, ok := value.(dal.DalClause); ok {
		value = gorm.Expr(expr.Expr, transformParams(expr.Params)...)
	}
	return errors.Convert(buildTx(d.db, clauses).Model(entity).Update(columnName, value).Error)
}

// UpdateColumns allows you to update multiple columns of mulitple records
func (d *Dalgorm) UpdateColumns(entity interface{}, set []dal.DalSet, clauses ...dal.Clause) errors.Error {
	updatesSet := make(map[string]interface{})

	for _, s := range set {
		if expr, ok := s.Value.(dal.DalClause); ok {
			s.Value = gorm.Expr(expr.Expr, transformParams(expr.Params)...)
		}
		updatesSet[s.ColumnName] = s.Value
	}

	return errors.Convert(buildTx(d.db, clauses).Model(entity).Updates(updatesSet).Error)
}

// UpdateAllColumn updated all Columns of entity
func (d *Dalgorm) UpdateAllColumn(entity interface{}, clauses ...dal.Clause) errors.Error {
	return errors.Convert(buildTx(d.db, clauses).UpdateColumns(entity).Error)
}

// GetColumns FIXME ...
func (d *Dalgorm) GetColumns(dst dal.Tabler, filter func(columnMeta dal.ColumnMeta) bool) (cms []dal.ColumnMeta, _ errors.Error) {
	columnTypes, err := d.db.Migrator().ColumnTypes(dst.TableName())
	if err != nil {
		return nil, errors.Convert(err)
	}
	for _, columnType := range columnTypes {
		if filter == nil {
			cms = append(cms, columnType)
		} else if filter(columnType) {
			cms = append(cms, columnType)
		}
	}
	return errors.Convert01(cms, nil)
}

// AddColumn add one column for the table
func (d *Dalgorm) AddColumn(table, columnName, columnType string) errors.Error {
	// work around the error `cached plan must not change result type` for postgres
	// wrap in func(){} to make the linter happy
	defer func() {
		_ = d.Exec("SELECT * FROM ? LIMIT 1", clause.Table{Name: table})
	}()
	return d.Exec("ALTER TABLE ? ADD ? ?", clause.Table{Name: table}, clause.Column{Name: columnName}, clause.Expr{SQL: columnType})
}

// DropColumns drop one column from the table
func (d *Dalgorm) DropColumns(table string, columnNames ...string) errors.Error {
	// work around the error `cached plan must not change result type` for postgres
	// wrap in func(){} to make the linter happy
	defer func() {
		_ = d.Exec("SELECT * FROM ? LIMIT 1", clause.Table{Name: table})
	}()
	for _, columnName := range columnNames {
		err := d.Exec("ALTER TABLE ? DROP COLUMN ?", clause.Table{Name: table}, clause.Column{Name: columnName})
		// err := d.db.Migrator().DropColumn(table, columnName)
		if err != nil {
			return errors.Convert(err)
		}
	}
	return nil
}

// GetPrimaryKeyFields get the PrimaryKey from `gorm` tag
func (d *Dalgorm) GetPrimaryKeyFields(t reflect.Type) []reflect.StructField {
	return utils.WalkFields(t, func(field *reflect.StructField) bool {
		return strings.Contains(strings.ToLower(field.Tag.Get("gorm")), "primarykey")
	})
}

// RenameColumn renames column name for specified table
func (d *Dalgorm) RenameColumn(table, oldColumnName, newColumnName string) errors.Error {
	// work around the error `cached plan must not change result type` for postgres
	// wrap in func(){} to make the linter happy
	defer func() {
		_ = d.Exec("SELECT * FROM ? LIMIT 1", clause.Table{Name: table})
	}()
	return d.Exec(
		"ALTER TABLE ? RENAME COLUMN ? TO ?",
		clause.Table{Name: table},
		clause.Column{Name: oldColumnName},
		clause.Column{Name: newColumnName},
	)
}

// AllTables returns all tables in the database
func (d *Dalgorm) AllTables() ([]string, errors.Error) {
	var tableSql string
	if d.db.Dialector.Name() == "mysql" {
		tableSql = "show tables"
	} else {
		tableSql = "select table_name from information_schema.tables where table_schema = 'public' and table_name not like '_devlake%'"
	}
	var tables []string
	err := d.db.Raw(tableSql).Scan(&tables).Error
	if err != nil {
		return nil, errors.Convert(err)
	}
	var filteredTables []string
	for _, table := range tables {
		if !strings.HasPrefix(table, "_devlake") {
			filteredTables = append(filteredTables, table)
		}
	}
	return filteredTables, nil
}

// DropTables drop multiple tables by Model Pointer or Table Name
func (d *Dalgorm) DropTables(dst ...interface{}) errors.Error {
	return errors.Convert(d.db.Migrator().DropTable(dst...))
}

// RenameTable renames table name
func (d *Dalgorm) RenameTable(oldName, newName string) errors.Error {
	return errors.Convert(d.db.Migrator().RenameTable(oldName, newName))
}

// DropIndexes drops indexes for specified table
func (d *Dalgorm) DropIndexes(table string, indexNames ...string) errors.Error {
	for _, indexName := range indexNames {
		err := d.db.Migrator().DropIndex(table, indexName)
		if err != nil {
			return errors.Convert(err)
		}
	}
	return nil
}

// Dialect returns the dialect of the database
func (d *Dalgorm) Dialect() string {
	return d.db.Dialector.Name()
}

// NewDalgorm FIXME ...
func NewDalgorm(db *gorm.DB) *Dalgorm {
	return &Dalgorm{db}
}
