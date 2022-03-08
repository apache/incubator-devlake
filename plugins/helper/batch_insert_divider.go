package helper

import (
	"reflect"

	"gorm.io/gorm"
)

type OnNewBatchInsert func(rowType reflect.Type) error

// Holds a map of BatchInsert, return `*BatchInsert` for a specific records, so caller can do batch operation for it
type BatchInsertDivider struct {
	db               *gorm.DB
	batches          map[reflect.Type]*BatchInsert
	batchSize        int
	onNewBatchInsert OnNewBatchInsert
}

// Return a new BatchInsertDivider instance
func NewBatchInsertDivider(db *gorm.DB, batchSize int) *BatchInsertDivider {
	return &BatchInsertDivider{
		db:        db,
		batches:   make(map[reflect.Type]*BatchInsert),
		batchSize: batchSize,
	}
}

func (d *BatchInsertDivider) OnNewBatchInsert(cb OnNewBatchInsert) {
	d.onNewBatchInsert = cb
}

// return *BatchInsert for specified type
func (d *BatchInsertDivider) ForType(rowType reflect.Type) (*BatchInsert, error) {
	// get the cache for the specific type
	batch := d.batches[rowType]
	var err error
	// create one if not exists
	if batch == nil {
		batch, err = NewBatchInsert(d.db, rowType, d.batchSize)
		if err != nil {
			return nil, err
		}
		if d.onNewBatchInsert != nil {
			err = d.onNewBatchInsert(rowType)
			if err != nil {
				return nil, err
			}
		}
		d.batches[rowType] = batch
	}
	return batch, nil
}

// close all batches so all rest records get saved into db as well
func (d *BatchInsertDivider) Close() error {
	for _, batch := range d.batches {
		err := batch.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
