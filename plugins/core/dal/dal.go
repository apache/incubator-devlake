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

import (
	"database/sql"
)

type Clause struct {
	Type string
	Data interface{}
}

// Dal aims to facilitate an isolation of Database Access Layer by defining a set of operations should a
// Database Access Layer provide
// This is inroduced by the fact that mocking *gorm.DB is hard, and `gomonkey` is not working on macOS
type Dal interface {
	// Raw executes raw sql query with sql.Rows and error return
	Raw(query string, params ...interface{}) (*sql.Rows, error)
	// Exec executes raw sql query
	Exec(query string, params ...interface{}) error
	// CreateTable creates a table with gorm definition from `entity`R
	AutoMigrate(entity interface{}, clauses ...Clause) error
	// Cursor returns a database cursor, cursor is especially useful when handling big amount of rows of data
	Cursor(clauses ...Clause) (*sql.Rows, error)
	// Fetch loads row data from `cursor` into `dst`
	Fetch(cursor *sql.Rows, dst interface{}) error
	// All loads matched rows from database to `dst`, USE IT WITH COUTIOUS!!
	All(dst interface{}, clauses ...Clause) error
	// First loads first matched row from database to `dst`, error will be returned if no records were found
	First(dst interface{}, clauses ...Clause) error
	// All loads matched rows from database to `dst`, USE IT WITH COUTIOUS!!
	Count(clauses ...Clause) (int64, error)
	// Pluck used to query single column
	Pluck(column string, dest interface{}, clauses ...Clause) error
	// Create insert record to database
	Create(entity interface{}, clauses ...Clause) error
	// Update updates record
	Update(entity interface{}, clauses ...Clause) error
	// CreateOrUpdate tries to create the record, or fallback to update all if failed
	CreateOrUpdate(entity interface{}, clauses ...Clause) error
	// CreateIfNotExist tries to create the record if not exist
	CreateIfNotExist(entity interface{}, clauses ...Clause) error
	// Delete records from database
	Delete(entity interface{}, clauses ...Clause) error
}

type DalClause struct {
	Expr   string
	Params []interface{}
}

const JoinClause string = "Join"

// Join creates a new JoinClause
func Join(clause string, params ...interface{}) Clause {
	return Clause{Type: JoinClause, Data: DalClause{clause, params}}
}

const WhereClause string = "Where"

// Where creates a new WhereClause
func Where(clause string, params ...interface{}) Clause {
	return Clause{Type: WhereClause, Data: DalClause{clause, params}}
}

const LimitClause string = "Limit"

// Limit creates a new LimitClause
func Limit(limit int) Clause {
	return Clause{Type: LimitClause, Data: limit}
}

const OffsetClause string = "Offset"

// Offset creates a new OffsetClause
func Offset(offset int) Clause {
	return Clause{Type: OffsetClause, Data: offset}
}

const FromClause string = "From"

// From creates a new TableClause
func From(table interface{}) Clause {
	return Clause{Type: FromClause, Data: table}
}

const SelectClause string = "Select"

// Select creates a new TableClause
func Select(fields string) Clause {
	return Clause{Type: SelectClause, Data: fields}
}

const OrderbyClause string = "OrderBy"

// Orderby creates a new Orderby
func Orderby(expr string) Clause {
	return Clause{Type: OrderbyClause, Data: expr}
}
