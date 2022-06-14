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
	"reflect"

	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
)

// BatchSaveDivider creates and caches BatchSave, this is helpful when dealing with massive amount of data records
// with arbitrary types.
type BatchSaveDivider struct {
	basicRes  core.BasicRes
	log       core.Logger
	db        dal.Dal
	batches   map[reflect.Type]*BatchSave
	batchSize int
	table     string
	params    string
}

// NewBatchSaveDivider create a new BatchInsertDivider instance
func NewBatchSaveDivider(basicRes core.BasicRes, batchSize int, table string, params string) *BatchSaveDivider {
	log := basicRes.GetLogger().Nested("batch divider")
	return &BatchSaveDivider{
		basicRes:  basicRes,
		log:       log,
		db:        basicRes.GetDal(),
		batches:   make(map[reflect.Type]*BatchSave),
		batchSize: batchSize,
		table:     table,
		params:    params,
	}
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
			return nil, fmt.Errorf("type %s must have RawDataOrigin embeded", rowElemType.Name())
		}
		// all good, delete outdated records before we insertion
		d.log.Debug("deleting outdate records for %s", rowElemType.Name())
		d.db.Delete(
			row,
			dal.Where("_raw_data_table = ? AND _raw_data_params = ?", d.table, d.params),
		)
	}
	return batch, nil
}

func (d *BatchSaveDivider) flushBatch() {

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
