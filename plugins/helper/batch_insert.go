package helper

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// Insert data by batch can increase database performance drastically, this class aim to make batch-insertion easier,
// It takes care the database operation for specified `slotType`, records got saved into database whenever cache hits
// The `size` limit, remember to call the `Close` method to save the last batch
type BatchInsert struct {
	slotType reflect.Type
	// slots can not be []interface{}, because gorm wouldn't take it
	// I'm guessing the reason is the type information lost when converted to interface{}
	slots               reflect.Value
	db                  *gorm.DB
	current             int
	size                int
}

func NewBatchInsert(db *gorm.DB, slotType reflect.Type, size int) (*BatchInsert, error) {
	if slotType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("slotType must be a pointer")
	}

	return &BatchInsert{
		slotType: slotType,
		slots:    reflect.MakeSlice(reflect.SliceOf(slotType), size, size),
		db:       db,
		size:     size,
	}, nil
}

func (c *BatchInsert) Add(slot interface{}) error {
	// type checking
	if reflect.TypeOf(slot) != c.slotType {
		return fmt.Errorf("sub cache type mismatched")
	}
	if reflect.ValueOf(slot).Kind() != reflect.Ptr {
		return fmt.Errorf("slot is not a pointer")
	}

	// push into slot
	c.slots.Index(c.current).Set(reflect.ValueOf(slot))
	c.current++
	// flush out into database if max outed
	if c.current == c.size {
		return c.Flush()
	}
	return nil
}

func (c *BatchInsert) Flush() error {
	err := c.db.CreateInBatches(c.slots.Slice(0, c.current).Interface(), c.current).Error
	c.current = 0
	if err != nil {
		return err
	}
	return nil
}

func (c *BatchInsert) Close() error {
	if c.current > 0 {
		return c.Flush()
	}
	return nil
}
