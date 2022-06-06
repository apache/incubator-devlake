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
	"github.com/apache/incubator-devlake/logger"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// go test -gcflags=all=-l

var TestBatchSize int = 100

func CreateTestBatchSaveDivider() *BatchSaveDivider {
	return NewBatchSaveDivider(&gorm.DB{}, TestBatchSize)
}

func TestBatchSaveDivider(t *testing.T) {
	MockDB(t)
	defer UnMockDB()
	batchSaveDivider := CreateTestBatchSaveDivider()
	initTimes := 0

	batchSaveDivider.OnNewBatchSave(func(rowType reflect.Type) error {
		initTimes++
		return nil
	})

	var err error
	var b1 *BatchSave
	var b2 *BatchSave
	var b3 *BatchSave

	// test if it saved and only saved once for one Type
	b1, err = batchSaveDivider.ForType(reflect.TypeOf(TestTableData), logger.Global)
	assert.Equal(t, initTimes, 1)
	assert.Equal(t, err, nil)
	b2, err = batchSaveDivider.ForType(reflect.TypeOf(&TestTable2{}), logger.Global)
	assert.Equal(t, initTimes, 2)
	assert.Equal(t, err, nil)
	b3, err = batchSaveDivider.ForType(reflect.TypeOf(TestTableData), logger.Global)
	assert.Equal(t, initTimes, 2)
	assert.Equal(t, err, nil)

	assert.NotEqual(t, b1, b2)
	assert.Equal(t, b1, b3)
}
