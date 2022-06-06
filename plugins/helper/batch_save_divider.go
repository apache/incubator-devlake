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

package helper

import (
	"github.com/apache/incubator-devlake/plugins/core"
	"reflect"

	"gorm.io/gorm"
)

type OnNewBatchSave func(rowType reflect.Type) error

// Holds a map of BatchInsert, return `*BatchInsert` for a specific records, so caller can do batch operation for it
type BatchSaveDivider struct {
	db             *gorm.DB
	batches        map[reflect.Type]*BatchSave
	batchSize      int
	onNewBatchSave OnNewBatchSave
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
	d.onNewBatchSave = cb
}

// return *BatchSave for specified type
func (d *BatchSaveDivider) ForType(rowType reflect.Type, log core.Logger) (*BatchSave, error) {
	// get the cache for the specific type
	batch := d.batches[rowType]
	var err error
	// create one if not exists
	if batch == nil {
		batch, err = NewBatchSave(d.db, log, rowType, d.batchSize)
		if err != nil {
			return nil, err
		}
		if d.onNewBatchSave != nil {
			err = d.onNewBatchSave(rowType)
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
