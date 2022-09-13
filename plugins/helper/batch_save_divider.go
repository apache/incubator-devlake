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
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"reflect"

	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
)

// BatchDivider is base struct of BatchSaveDivider&BatchUpdateDivider
type BatchDivider struct {
	basicRes  core.BasicRes
	db        dal.Dal
	batchSize int
	table     string
	params    string
}

// BatchSaveDivider creates and caches BatchSave, this is helpful when dealing with massive amount of data records
// with arbitrary types.
type BatchSaveDivider struct {
	*BatchDivider
	log     core.Logger
	batches map[reflect.Type]*BatchSave
}

// BatchUpdateDivider creates and caches BatchUpdate, this is helpful when dealing with massive amount of data records
// with arbitrary types.
type BatchUpdateDivider struct {
	*BatchDivider
	log     core.Logger
	batches map[reflect.Type]*BatchUpdate
}

// NewBatchSaveDivider create a new BatchInsertDivider instance
func NewBatchSaveDivider(basicRes core.BasicRes, batchSize int, table string, params string) *BatchSaveDivider {
	log := basicRes.GetLogger().Nested("batch divider")
	batchDivider := &BatchDivider{
		basicRes:  basicRes,
		db:        basicRes.GetDal(),
		batchSize: batchSize,
		table:     table,
		params:    params,
	}
	batchSaveDivider := &BatchSaveDivider{
		log:     log,
		batches: make(map[reflect.Type]*BatchSave),
	}
	batchSaveDivider.BatchDivider = batchDivider
	return batchSaveDivider
}

// NewBatchUpdateDivider create a new BatchInsertDivider instance
func NewBatchUpdateDivider(basicRes core.BasicRes, batchSize int, table string, params string) *BatchUpdateDivider {
	log := basicRes.GetLogger().Nested("batch update divider")
	batchDivider := &BatchDivider{
		basicRes:  basicRes,
		db:        basicRes.GetDal(),
		batchSize: batchSize,
		table:     table,
		params:    params,
	}
	batchUpdateDivider := &BatchUpdateDivider{
		log:     log,
		batches: make(map[reflect.Type]*BatchUpdate),
	}
	batchUpdateDivider.BatchDivider = batchDivider
	return batchUpdateDivider
}

// ForType returns a `BatchSave` instance for specific type
func (d *BatchSaveDivider) ForType(rowType reflect.Type) (*BatchSave, error) {
	// get the cache for the specific type
	batch := d.batches[rowType]
	var err error
	// create one if not exists
	if batch == nil {
		batch, err = NewBatchSave(d.basicRes, rowType, d.batchSize)
		if err != nil {
			return nil, err
		}
		d.batches[rowType] = batch
		// delete outdated records if rowType was not PartialUpdate
		rowElemType := rowType.Elem()
		d.log.Debug("missing BatchSave for type %s", rowElemType.Name())
		row := reflect.New(rowElemType).Interface()
		// check if rowType had RawDataOrigin embeded
		field, hasField := rowElemType.FieldByName("RawDataOrigin")
		if !hasField || field.Type != reflect.TypeOf(common.RawDataOrigin{}) {
			return nil, errors.Default.New(fmt.Sprintf("type %s must have RawDataOrigin embeded", rowElemType.Name()))
		}
		// all good, delete outdated records before we insertion
		d.log.Debug("deleting outdate records for %s", rowElemType.Name())
		err = d.db.Delete(
			row,
			dal.Where("_raw_data_table = ? AND _raw_data_params = ?", d.table, d.params),
		)
		if err != nil {
			return nil, err
		}
	}
	return batch, nil
}

// Close all batches so the rest records get saved into db
func (d *BatchSaveDivider) Close() error {
	for _, batch := range d.batches {
		err := batch.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// ForType returns a `BatchUpdate` instance for specific type
func (d *BatchUpdateDivider) ForType(rowType reflect.Type) (*BatchUpdate, error) {
	// get the cache for the specific type
	batch := d.batches[rowType]
	var err error
	// create one if not exists
	if batch == nil {
		batch, err = NewBatchUpdate(d.basicRes, rowType, d.batchSize)
		if err != nil {
			return nil, err
		}
		d.batches[rowType] = batch
	}
	return batch, nil
}

// Close all batches so the rest records get saved into db
func (d *BatchUpdateDivider) Close() error {
	for _, batch := range d.batches {
		err := batch.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
