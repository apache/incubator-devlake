package helper

import (
	"database/sql"
	"reflect"

	"gorm.io/gorm"
)

type Iterator interface {
	HasNext() bool
	Fetch() (interface{}, error)
	Close() error
}

type CursorIterator struct {
	db       *gorm.DB
	cursor   *sql.Rows
	elemType reflect.Type
}

func NewCursorIterator(db *gorm.DB, cursor *sql.Rows, elemType reflect.Type) (*CursorIterator, error) {
	return &CursorIterator{
		db:       db,
		cursor:   cursor,
		elemType: elemType,
	}, nil
}

func (c *CursorIterator) HasNext() bool {
	return c.cursor.Next()
}

func (c *CursorIterator) Fetch() (interface{}, error) {
	elem := reflect.New(c.elemType).Interface()
	err := c.db.ScanRows(c.cursor, elem)
	if err != nil {
		return nil, err
	}
	return elem, nil
}

func (c *CursorIterator) Close() error {
	return c.cursor.Close()
}

var _ Iterator = (*CursorIterator)(nil)
