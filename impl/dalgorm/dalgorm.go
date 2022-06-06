package dalgorm

import (
	"database/sql"

	"github.com/apache/incubator-devlake/plugins/core/dal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Dalgorm struct {
	db *gorm.DB
}

// To accommodate gorm
type stubTable struct {
	name string
}

func (s *stubTable) TableName() string {
	return s.name
}

func buildTx(tx *gorm.DB, clauses []interface{}) *gorm.DB {
	for _, clause := range clauses {
		switch clause := clause.(type) {
		case *dal.JoinClause:
			tx = tx.Joins(clause.Expr, clause.Params...)
		case *dal.WhereClause:
			tx = tx.Where(clause.Expr, clause.Params...)
		case dal.OrderbyClause:
			tx = tx.Order(string(clause))
		case dal.LimitClause:
			tx = tx.Limit(int(clause))
		case dal.OffsetClause:
			tx = tx.Offset(int(clause))
		case dal.TableClause:
			tx = tx.Table(string(clause))
		}
	}
	return tx
}

var _ dal.Dal = (*Dalgorm)(nil)

// Exec executes raw sql query
func (d *Dalgorm) Exec(query string, params ...interface{}) error {
	return d.db.Exec(query, params...).Error
}

// CreateTable creates a table with gorm definition from `entity`
func (d *Dalgorm) AutoMigrate(entity interface{}, clauses ...interface{}) error {
	return buildTx(d.db, clauses).AutoMigrate(entity)
}

// Cursor returns a database cursor, cursor is especially useful when handling big amount of rows of data
func (d *Dalgorm) Cursor(clauses ...interface{}) (dal.Cursor, error) {
	return buildTx(d.db, clauses).Rows()
}

// Fetch loads row data from `cursor` into `dst`
func (d *Dalgorm) Fetch(cursor dal.Cursor, dst interface{}) error {
	return d.db.ScanRows(cursor.(*sql.Rows), dst)
}

// All loads matched rows from database to `dst`, USE IT WITH COUTIOUS!!
func (d *Dalgorm) All(dst interface{}, clauses ...interface{}) error {
	return buildTx(d.db, clauses).Find(dst).Error
}

// First loads first matched row from database to `dst`, error will be returned if no records were found
func (d *Dalgorm) First(dst interface{}, clauses ...interface{}) error {
	return buildTx(d.db, clauses).First(dst).Error
}

// Create insert record to database
func (d *Dalgorm) Create(entity interface{}, clauses ...interface{}) error {
	return buildTx(d.db, clauses).Create(entity).Error
}

// Update updates record
func (d *Dalgorm) Update(entity interface{}, clauses ...interface{}) error {
	return buildTx(d.db, clauses).Save(entity).Error
}

// CreateOrUpdate tries to create the record, or fallback to update all if failed
func (d *Dalgorm) CreateOrUpdate(entity interface{}, clauses ...interface{}) error {
	return buildTx(d.db, clauses).Clauses(clause.OnConflict{UpdateAll: true}).Create(entity).Error
}

// Delete records from database
func (d *Dalgorm) Delete(entity interface{}, clauses ...interface{}) error {
	return buildTx(d.db, clauses).Delete(entity).Error
}

func NewDalgorm(db *gorm.DB) *Dalgorm {
	return &Dalgorm{db}
}
