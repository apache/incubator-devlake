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
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/magiconair/properties/assert"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// go test -gcflags=all=-l
func CreateTestDataConverter(t *testing.T) (*DataConverter, error) {
	var ctx context.Context
	ctx, Cancel = context.WithCancel(context.Background())
	return NewDataConverter(DataConverterArgs{
		RawDataSubTaskArgs: RawDataSubTaskArgs{
			Ctx: &DefaultSubTaskContext{
				defaultExecContext: newDefaultExecContext(GetConfigForTest("../../"), logger.NewDefaultLogger(logrus.New(), "Test", make(map[string]*logrus.Logger)), &gorm.DB{}, ctx, "Test", nil, nil),
			},
			Table: TestTable{}.TableName(),
			Params: &TestParam{
				Test: TestUrlParam,
			},
		},
		InputRowType: reflect.TypeOf(TestTable{}),
		Input:        &sql.Rows{},
		BatchSize:    TestBatchSize,
		Convert: func(row interface{}) ([]interface{}, error) {
			assert.Equal(t, row, TestTableData)
			results := make([]interface{}, 0, TestDataCount)
			for i := 0; i < TestDataCount; i++ {
				results = append(results, TestTableData)
			}
			return results, nil
		},
	})
}

func TestDataConvertorExecute(t *testing.T) {
	MockDB(t)
	defer UnMockDB()

	gt.Reset()
	gt = gomonkey.ApplyMethod(reflect.TypeOf(&gorm.DB{}), "Table", func(db *gorm.DB, name string, args ...interface{}) *gorm.DB {
		assert.Equal(t, name, "_raw_"+TestTableData.TableName())
		return db
	},
	)

	dataConvertor, _ := CreateTestDataConverter(t)

	datacount := TestDataCount
	gn := gomonkey.ApplyMethod(reflect.TypeOf(&sql.Rows{}), "Next", func(r *sql.Rows) bool {
		if datacount > 0 {
			datacount--
			return true
		} else {
			return false
		}
	})
	defer gn.Reset()

	gcl := gomonkey.ApplyMethod(reflect.TypeOf(&sql.Rows{}), "Close", func(r *sql.Rows) error {
		return nil
	})
	defer gcl.Reset()

	gs.Reset()

	scanrowTimes := 0
	gs = gomonkey.ApplyMethod(reflect.TypeOf(&gorm.DB{}), "ScanRows", func(db *gorm.DB, rows *sql.Rows, dest interface{}) error {
		scanrowTimes++
		*dest.(*TestTable) = *TestTableData
		return nil
	},
	)

	fortypeTimes := 0
	gf := gomonkey.ApplyMethod(reflect.TypeOf(&BatchSaveDivider{}), "ForType", func(d *BatchSaveDivider, rowType reflect.Type, log core.Logger) (*BatchSave, error) {
		fortypeTimes++
		assert.Equal(t, rowType, reflect.TypeOf(TestTableData))
		err := d.onNewBatchSave(rowType)
		assert.Equal(t, err, nil)

		return &BatchSave{}, nil
	})
	defer gf.Reset()

	addTimes := 0
	gsv := gomonkey.ApplyMethod(reflect.TypeOf(&BatchSave{}), "Add", func(c *BatchSave, slot interface{}) error {
		addTimes++
		assert.Equal(t, slot, TestTableData)
		return nil
	})
	defer gsv.Reset()

	// begin testing
	err := dataConvertor.Execute()
	assert.Equal(t, err, nil)
	assert.Equal(t, scanrowTimes, TestDataCount)
	assert.Equal(t, fortypeTimes, TestDataCount*TestDataCount)
	assert.Equal(t, addTimes, TestDataCount*TestDataCount)
}

func TestDataConvertorExecute_Cancel(t *testing.T) {
	MockDB(t)
	defer UnMockDB()

	gt.Reset()
	gt = gomonkey.ApplyMethod(reflect.TypeOf(&gorm.DB{}), "Table", func(db *gorm.DB, name string, args ...interface{}) *gorm.DB {
		assert.Equal(t, name, "_raw_"+TestTableData.TableName())
		return db
	},
	)

	dataConvertor, _ := CreateTestDataConverter(t)

	gn := gomonkey.ApplyMethod(reflect.TypeOf(&sql.Rows{}), "Next", func(r *sql.Rows) bool {
		// death loop for testing cancel
		return true
	})
	defer gn.Reset()

	gcl := gomonkey.ApplyMethod(reflect.TypeOf(&sql.Rows{}), "Close", func(r *sql.Rows) error {
		return nil
	})
	defer gcl.Reset()

	gs.Reset()

	scanrowTimes := 0
	gs = gomonkey.ApplyMethod(reflect.TypeOf(&gorm.DB{}), "ScanRows", func(db *gorm.DB, rows *sql.Rows, dest interface{}) error {
		scanrowTimes++
		*dest.(*TestTable) = *TestTableData
		return nil
	},
	)

	fortypeTimes := 0
	gf := gomonkey.ApplyMethod(reflect.TypeOf(&BatchSaveDivider{}), "ForType", func(d *BatchSaveDivider, rowType reflect.Type, log core.Logger) (*BatchSave, error) {
		fortypeTimes++
		assert.Equal(t, rowType, reflect.TypeOf(TestTableData))
		err := d.onNewBatchSave(rowType)
		assert.Equal(t, err, nil)

		return &BatchSave{}, nil
	})
	defer gf.Reset()

	addTimes := 0
	gsv := gomonkey.ApplyMethod(reflect.TypeOf(&BatchSave{}), "Add", func(c *BatchSave, slot interface{}) error {
		addTimes++
		assert.Equal(t, slot, TestTableData)
		return nil
	})
	defer gsv.Reset()

	go func() {
		time.Sleep(time.Duration(500) * time.Microsecond)
		Cancel()
	}()

	err := dataConvertor.Execute()
	assert.Equal(t, err, fmt.Errorf("context canceled"))
}
