package helper

import (
	"database/sql"
	"reflect"
	"time"

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

type DateIterator struct {
	startTime time.Time
	endTime   time.Time
	Days      int
	Current   int
}

type DatePair struct {
	PairStartTime time.Time
	PairEndTime   time.Time
}

func (c *DateIterator) HasNext() bool {
	return c.Current < c.Days
}

func (c *DateIterator) Fetch() (interface{}, error) {
	c.Current++
	return &DatePair{
		PairStartTime: c.startTime.AddDate(0, 0, c.Current),
		PairEndTime:   c.endTime.AddDate(0, 0, c.Current),
	}, nil

}

func (c *DateIterator) Close() error {
	return nil
}

func NewDateIterator(days int) (*DateIterator, error) {
	endTime := time.Now().Truncate(24 * time.Hour)
	return &DateIterator{
		startTime: endTime.AddDate(0, 0, -days-1),
		endTime:   endTime.AddDate(0, 0, -days),
		Days:      days,
		Current:   0,
	}, nil
}
