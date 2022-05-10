package testhelper

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"encoding/csv"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/runner"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// DataFlowTester facilitate a universal data integration e2e-test to help you verifying records between each step
type DataFlowTester struct {
	Cfg    *viper.Viper
	Db     *gorm.DB
	T      *testing.T
	Name   string
	Plugin core.PluginMeta
	Log    core.Logger
}

type CsvFileIterator struct {
	file   *os.File
	reader *csv.Reader
	fields []string
	row    map[string]interface{}
}

func NewCsvFileIterator(csvPath string) *CsvFileIterator {
	// open csv file
	csvFile, err := os.Open(csvPath)
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(csvFile)
	// load field names
	fields, err := csvReader.Read()
	if err != nil {
		panic(err)
	}
	return &CsvFileIterator{
		file:   csvFile,
		reader: csvReader,
		fields: fields,
	}
}

func (ci *CsvFileIterator) Close() {
	err := ci.file.Close()
	if err != nil {
		panic(err)
	}
}

func (ci *CsvFileIterator) HasNext() bool {
	row, err := ci.reader.Read()
	if err == io.EOF {
		ci.row = nil
		return false
	}
	if err != nil {
		ci.row = nil
		panic(err)
	}
	// convert row tuple to map type, so gorm can insert data with it
	ci.row = make(map[string]interface{})
	for index, field := range ci.fields {
		ci.row[field] = row[index]
	}
	return true
}

func (ci *CsvFileIterator) Fetch() map[string]interface{} {
	return ci.row
}

// NewDataFlowTester create a *DataFlowTester to help developer test their subtasks data flow
func NewDataFlowTester(t *testing.T, pluginName string, pluginMeta core.PluginMeta) *DataFlowTester {
	core.RegisterPlugin(pluginName, pluginMeta)
	cfg := config.GetConfig()
	db, err := runner.NewGormDb(cfg, logger.Global)
	if err != nil {
		panic(err)
	}
	return &DataFlowTester{
		Cfg:    cfg,
		Db:     db,
		T:      t,
		Name:   pluginName,
		Plugin: pluginMeta,
		Log:    logger.Global,
	}
}

func (t *DataFlowTester) ImportCsv(csvRelPath string, tableName string) {
	csvIter := NewCsvFileIterator(csvRelPath)
	defer csvIter.Close()
	// create table if not exists
	err := t.Db.Table(tableName).AutoMigrate(&helper.RawData{})
	if err != nil {
		panic(err)
	}
	t.FlushTable(tableName)
	// load rows and insert into target table
	for csvIter.HasNext() {
		// make sure
		result := t.Db.Table(tableName).Create(csvIter.Fetch())
		if result.Error != nil {
			panic(result.Error)
		}
		assert.Equal(t.T, int64(1), result.RowsAffected)
	}
}

func (t *DataFlowTester) FlushTable(tableName string) {
	// flush target table
	err := t.Db.Exec(fmt.Sprintf("DELETE FROM %s", tableName)).Error
	if err != nil {
		panic(err)
	}
}

func (t *DataFlowTester) Subtask(subtaskMeta core.SubTaskMeta, taskData interface{}) {
	subtaskCtx := helper.NewStandaloneSubTaskContext(t.Cfg, t.Log, t.Db, context.Background(), t.Name, taskData)
	subtaskMeta.EntryPoint(subtaskCtx)
}

func (t *DataFlowTester) VerifyTable(tableName string, csvRelPath string, pkfields []string, targetfields []string) {
	csvIter := NewCsvFileIterator(csvRelPath)
	defer csvIter.Close()
	for csvIter.HasNext() {
		expected := csvIter.Fetch()
		pkvalues := make([]interface{}, 0, len(pkfields))
		for _, pkf := range pkfields {
			pkvalues = append(pkvalues, expected[pkf])
		}
		actual := make(map[string]interface{})
		where := ""
		for _, field := range pkfields {
			where += fmt.Sprintf(" %s = ?", field)
		}
		err := t.Db.Table(tableName).Where(where, pkvalues...).Find(actual).Error
		if err != nil {
			panic(err)
		}
		for _, field := range targetfields {
			actualValue := ""
			switch actual[field].(type) {
			// TODO: ensure testing database is in UTC timezone
			case time.Time:
				actualValue = actual[field].(time.Time).Format("2006-01-02 15:04:05.000000000")
			default:
				actualValue = fmt.Sprint(actual[field])
			}
			assert.Equal(t.T, expected[field], actualValue)
		}
	}
}
