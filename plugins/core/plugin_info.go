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

package core

import (
	"fmt"
	"reflect"
	"sync"

	"gorm.io/gorm/schema"
)

type TableInfo struct {
	name     string
	describe string
	table    schema.Tabler
	mutex    *sync.Mutex
}
type TableInfoCallBack func(pluginName string, name string, describe *string, table schema.Tabler) (err error)

type TablesInfo struct {
	name      string
	tableInfo map[string]*TableInfo
	mutex     *sync.Mutex
}
type TablesInfoCallBack func(pluginName string, info *TablesInfo) (err error)

type TotalInfo struct {
	tablesInfo map[string]*TablesInfo
}

var totalInfo TotalInfo = TotalInfo{tablesInfo: make(map[string]*TablesInfo)}

type PluginInfo interface {
	RegisterTablesInfo(ts *TablesInfo) error
	//ReleaseModelsInfo() error
}

func GetTotalInfo() *TotalInfo {
	return &totalInfo
}

func RegisterPluginInfo(pi PluginInfo) error {
	piType := reflect.ValueOf(pi)
	if piType.Kind() == reflect.Ptr {
		piType = piType.Elem()
	}
	ts := &TablesInfo{
		name:      piType.Type().Name(),
		tableInfo: make(map[string]*TableInfo),
		mutex:     &sync.Mutex{},
	}
	// regist model info
	if ts.name == "" {
		return fmt.Errorf("can not use empty string on name to RegisterModelsInfo() for plugin")
	}

	if ts, ok := totalInfo.tablesInfo[ts.name]; ok {
		return fmt.Errorf("name [%s] has been used for RegisterModelsInfo()", ts.name)
	}

	totalInfo.tablesInfo[ts.name] = ts
	err := pi.RegisterTablesInfo(ts)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TablesInfo) Add(table schema.Tabler, describe ...string) error {
	modelType := reflect.ValueOf(table)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	desc := "[" + table.TableName() + "]"
	for _, d := range describe {
		desc += d
	}

	m := &TableInfo{
		name:     modelType.Type().Name(),
		describe: desc,
		table:    table,
		mutex:    &sync.Mutex{},
	}
	if _, ok := ts.tableInfo[m.name]; ok {
		return fmt.Errorf("the table name [%s] has been used in plugin [%s]", m.name, ts.name)
	}
	ts.tableInfo[m.name] = m
	return nil
}

func (ts *TablesInfo) Adds(tables ...schema.Tabler) error {
	for _, table := range tables {
		err := ts.Add(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ts *TableInfo) Name() string {
	return ts.name
}

func (ts *TablesInfo) Visit(name string, handle TableInfoCallBack) error {
	if m, ok := ts.tableInfo[name]; ok {
		m.mutex.Lock()
		err := handle(ts.name, name, &m.describe, m.table)
		m.mutex.Unlock()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("can not find out the table with name [%s] in plugin [%s]", name, ts.name)
	}
	return nil
}

func (ts *TablesInfo) VisitWithOutLock(name string, handle TableInfoCallBack) error {
	if m, ok := ts.tableInfo[name]; ok {
		err := handle(ts.name, name, &m.describe, m.table)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("can not find out the table with name [%s] in plugin [%s]", name, ts.name)
	}
	return nil
}

func (ts *TablesInfo) Traversal(handle TableInfoCallBack) error {
	for name, m := range ts.tableInfo {
		m.mutex.Lock()
		err := handle(ts.name, name, &m.describe, m.table)
		m.mutex.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ts *TablesInfo) TraversalWithOutLock(handle TableInfoCallBack) error {
	for name, m := range ts.tableInfo {
		err := handle(ts.name, name, &m.describe, m.table)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TotalInfo) VisitTable(name string, handle TablesInfoCallBack) error {
	if ts, ok := t.tablesInfo[name]; ok {
		ts.mutex.Lock()
		err := handle(name, ts)
		ts.mutex.Unlock()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("can not find out the plugin with name [%s] for table", name)
	}
	return nil
}

func (t *TotalInfo) VisitTableWithOutLock(name string, handle TablesInfoCallBack) error {
	if ts, ok := t.tablesInfo[name]; ok {
		err := handle(name, ts)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("can not find out the plugin with name [%s] for table", name)
	}
	return nil
}

func (t *TotalInfo) TraversalTable(handle TablesInfoCallBack) error {
	for name, ts := range t.tablesInfo {
		ts.mutex.Lock()
		err := handle(name, ts)
		ts.mutex.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TotalInfo) TraversalTableWithOutLock(handle TablesInfoCallBack) error {
	for name, ts := range t.tablesInfo {
		err := handle(name, ts)
		if err != nil {
			return err
		}
	}
	return nil
}
