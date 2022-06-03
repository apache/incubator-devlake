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

package e2ehelper

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/helpers/pluginhelper"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// DataFlowTester provides a universal data integrity validation facility to help `Plugin` verifying records between
// each step.
//
// How it works:
//
//   1. Flush specified `table` and import data from a `csv` file by `ImportCsv` method
//   2. Execute specified `subtask` by `Subtask` method
//   3. Verify actual data from specified table against expected data from another `csv` file
//   4. Repeat step 2 and 3
//
// Recommended Usage:

//   1. Create a folder under your plugin root folder. i.e. `plugins/gitlab/e2e/ to host all your e2e-tests`
//   2. Create a folder named `tables` to hold all data in `csv` format
//   3. Create e2e test-cases to cover all possible data-flow routes
//
// Example code:
//
//   See [Gitlab Project Data Flow Test](plugins/gitlab/e2e/project_test.go) for detail
//
// DataFlowTester use `N`
type DataFlowTester struct {
	Cfg    *viper.Viper
	Db     *gorm.DB
	T      *testing.T
	Name   string
	Plugin core.PluginMeta
	Log    core.Logger
}

// NewDataFlowTester create a *DataFlowTester to help developer test their subtasks data flow
func NewDataFlowTester(t *testing.T, pluginName string, pluginMeta core.PluginMeta) *DataFlowTester {
	err := core.RegisterPlugin(pluginName, pluginMeta)
	if err != nil {
		panic(err)
	}
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

// ImportCsv imports records from specified csv file into target table, note that existing data would be deleted first.
func (t *DataFlowTester) ImportCsv(csvRelPath string, tableName string) {
	csvIter := pluginhelper.NewCsvFileIterator(csvRelPath)
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

// FlushTable deletes all records from specified table
func (t *DataFlowTester) FlushTable(tableName string) {
	// flush target table
	err := t.Db.Exec(fmt.Sprintf("DELETE FROM %s", tableName)).Error
	if err != nil {
		panic(err)
	}
}

// Subtask executes specified subtasks
func (t *DataFlowTester) Subtask(subtaskMeta core.SubTaskMeta, taskData interface{}) {
	subtaskCtx := helper.NewStandaloneSubTaskContext(t.Cfg, t.Log, t.Db, context.Background(), t.Name, taskData)
	err := subtaskMeta.EntryPoint(subtaskCtx)
	if err != nil {
		panic(err)
	}
}

// VerifyTable reads rows from csv file and compare with records from database one by one. You must specified the
// Primary Key Fields with `pkfields` so DataFlowTester could select the exact record from database, as well as which
// fields to compare with by specifying `targetfields` parameter.
func (t *DataFlowTester) VerifyTable(tableName string, csvRelPath string, pkfields []string, targetfields []string) {
	csvIter := pluginhelper.NewCsvFileIterator(csvRelPath)
	defer csvIter.Close()

	var expectedTotal int64
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
		expectedTotal++
	}

	var actualTotal int64
	err := t.Db.Table(tableName).Count(&actualTotal).Error
	if err != nil {
		panic(err)
	}
	assert.Equal(t.T, expectedTotal, actualTotal)
}
