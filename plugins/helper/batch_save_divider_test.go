package helper

import (
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
	b1, err = batchSaveDivider.ForType(reflect.TypeOf(TestTableData))
	assert.Equal(t, initTimes, 1)
	assert.Equal(t, err, nil)
	b2, err = batchSaveDivider.ForType(reflect.TypeOf(&TestTable2{}))
	assert.Equal(t, initTimes, 2)
	assert.Equal(t, err, nil)
	b3, err = batchSaveDivider.ForType(reflect.TypeOf(TestTableData))
	assert.Equal(t, initTimes, 2)
	assert.Equal(t, err, nil)

	assert.NotEqual(t, b1, b2)
	assert.Equal(t, b1, b3)
}
