package helper

import (
	"reflect"

	"gorm.io/gorm"
)

type OnNewBatchSave func(rowType reflect.Type) error

// Holds a map of BatchInsert, return `*BatchInsert` for a specific records, so caller can do batch operation for it
type BatchSaveDivider struct {
	db               *gorm.DB
	batches          map[reflect.Type]*BatchSave
	batchSize        int
	onNewBatchInsert OnNewBatchSave
}

// Return a new BatchInsertDivider instance
func NewBatchSaveDivider(db *gorm.DB, batchSize int) *BatchSaveDivider {
	return &BatchSaveDivider{
		db:        db,
		batches:   make(map[reflect.Type]*BatchSave),
		batchSize: batchSize,
	}
}

func (d *BatchSaveDivider) OnNewBatchSave(cb OnNewBatchSave) {
	d.onNewBatchInsert = cb
}

// return *BatchInsert for specified type
func (d *BatchSaveDivider) ForType(rowType reflect.Type) (*BatchSave, error) {
	// get the cache for the specific type
	batch := d.batches[rowType]
	var err error
	// create one if not exists
	if batch == nil {
		batch, err = NewBatchSave(d.db, rowType, d.batchSize)
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
func (d *BatchSaveDivider) Close() error {
	for _, batch := range d.batches {
		err := batch.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
