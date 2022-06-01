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

package dal

// Dal aims to facilitate an isolation of Database Access Layer by defining a set of operations should a
// Database Access Layer provide
// This is inroduced by the fact that mocking *gorm.DB is hard, and `gomonkey` is not working on macOS
type Dal interface {
	// Exec executes raw sql query
	Exec(query string, params ...interface{}) error
	// CreateTable creates a table with gorm definition from `entity`
	AutoMigrate(entity interface{}, clauses ...interface{}) error
	// Cursor returns a database cursor, cursor is especially useful when handling big amount of rows of data
	Cursor(clauses ...interface{}) (Cursor, error)
	// Fetch loads row data from `cursor` into `dst`
	Fetch(cursor Cursor, dst interface{}) error
	// All loads matched rows from database to `dst`, USE IT WITH COUTIOUS!!
	All(dst interface{}, clauses ...interface{}) error
	// First loads first matched row from database to `dst`, error will be returned if no records were found
	First(dst interface{}, clauses ...interface{}) error
	// Create insert record to database
	Create(entity interface{}, clauses ...interface{}) error
	// Update updates record
	Update(entity interface{}, clauses ...interface{}) error
	// CreateOrUpdate tries to create the record, or fallback to update all if failed
	CreateOrUpdate(entity interface{}, clauses ...interface{}) error
	// Delete records from database
	Delete(entity interface{}, clauses ...interface{}) error
}

// Cursor represents a database cursor
type Cursor interface {
	Close() error
	Next() bool
}

type dalClause struct {
	Expr   string
	Params []interface{}
}

// JoinClause represents a SQL `JOIN` clause
type JoinClause dalClause

// Join creates a new JoinClause
func Join(clause string, params ...interface{}) *JoinClause {
	return &JoinClause{clause, params}
}

// WhereClause represents a SQL `WHERE` clause
type WhereClause dalClause

// Where creates a new WhereClause
func Where(clause string, params ...interface{}) *WhereClause {
	return &WhereClause{clause, params}
}

// LimitClause represents a SQL `LIMIT` clause
type LimitClause int

// Limit creates a new LimitClause
func Limit(limit int) LimitClause {
	return LimitClause(limit)
}

// OffsetClause represents a SQL `OFFSET` clause
type OffsetClause int

// Offset creates a new OffsetClause
func Offset(offset int) OffsetClause {
	return OffsetClause(offset)
}

// TableClause represents a SQL `OFFSET` clause
type TableClause string

// Table creates a new TableClause
func Table(table string) TableClause {
	return TableClause(table)
}

// OrderbyClause represents a SQL `ORDER BY` clause
type OrderbyClause string

// Orderby creates a new Orderby
func Orderby(expr string) OrderbyClause {
	return OrderbyClause(expr)
}
