package testhelper

import (
	"fmt"
	"io"
	"os"
	"testing"

	"encoding/csv"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/runner"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// DataFlowTester helps you write subtasks e2e tests with ease
type DataFlowTester struct {
	Db     *gorm.DB
	T      *testing.T
	Plugin core.PluginMeta
}

// NewDataFlowTester create a *DataFlowTester to help developer test their subtasks data flow
func NewDataFlowTester(t *testing.T, plugin core.PluginMeta) *DataFlowTester {
	cfg := config.GetConfig()
	db, err := runner.NewGormDb(cfg, logger.Global)
	if err != nil {
		panic(err)
	}
	return &DataFlowTester{
		Db:     db,
		T:      t,
		Plugin: plugin,
	}
}

func (t *DataFlowTester) ImportCsv(csvRelPath string, tableName string) {
	// open csv file
	csvFile, err := os.Open(csvRelPath)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)
	// load field names
	fields, err := csvReader.Read()
	if err != nil {
		panic(err)
	}
	// create table if not exists
	err = t.Db.Table(tableName).AutoMigrate(&helper.RawData{})
	if err != nil {
		panic(err)
	}
	// flush target table
	err = t.Db.Exec(fmt.Sprintf("DELETE FROM %s", tableName)).Error
	if err != nil {
		panic(err)
	}
	// load rows and insert into target table
	data := make(map[string]interface{})
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		// convert row tuple to map type, so gorm can insert data with it
		for index, field := range fields {
			data[field] = row[index]
		}
		// make sure
		result := t.Db.Table(tableName).Create(data)
		if result.Error != nil {
			panic(result.Error)
		}
		assert.Equal(t.T, int64(1), result.RowsAffected)
	}
}

func (t *DataFlowTester) Subtask(subtaskMeta core.SubTaskMeta) {
}

func (T *DataFlowTester) VerifyTable(tableName string, csvRelPaht string) {
}
