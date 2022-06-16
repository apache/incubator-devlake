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
package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"io/ioutil"
	"net/http"
	"strings"
)

func LoadData(c core.SubTaskContext) error {
	config := c.GetData().(*StarRocksConfig)
	db := c.GetDal()
	tables := config.Tables
	var err error
	if len(tables) == 0 {
		tables, err = db.AllTables()
		if err != nil {
			return err
		}
	}
	starrocks, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Password, config.Host, config.Port, config.Database))
	if err != nil {
		return err
	}

	for _, table := range tables {
		err = loadData(starrocks, c, table, db, config)
		if err != nil {
			c.GetLogger().Error("load data %s error: %s", table, err)
		}
	}
	return nil
}
func loadData(starrocks *sql.DB, c core.SubTaskContext, table string, db dal.Dal, config *StarRocksConfig) error {
	var data []map[string]interface{}
	// select data from db
	rows, err := db.Raw(fmt.Sprintf("select * from %s", table))
	if err != nil {
		return err
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	for rows.Next() {
		row := make(map[string]interface{})
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		err = rows.Scan(columnPointers...)
		if err != nil {
			return err
		}
		for i, colName := range cols {
			row[colName] = columns[i]
		}
		data = append(data, row)
	}
	if len(data) == 0 {
		c.GetLogger().Warn("table %s is empty, so skip", table)
		return nil
	}
	starrocksTable := strings.TrimLeft(table, "_")
	// create tmp table in starrocks
	_, err = starrocks.Exec(fmt.Sprintf("create table %s_tmp like %s", starrocksTable, starrocksTable))
	if err != nil {
		return err
	}
	// insert data to tmp table
	url := fmt.Sprintf("http://%s:%d/api/%s/%s_tmp/_stream_load", config.Host, config.BePort, config.Database, starrocksTable)
	headers := map[string]string{
		"format":            "json",
		"strip_outer_array": "true",
		"Expect":            "100-continue",
		"ignore_json_size":  "true",
	}
	// marshal User to json
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(config.User, config.Password)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result map[string]interface{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		c.GetLogger().Error("%s %s", resp.StatusCode, b)
	}
	if result["Status"] != "Success" {
		c.GetLogger().Error("load %s failed: %s", table, b)
	} else {
		// drop old table and rename tmp table to old table
		_, err = starrocks.Exec(fmt.Sprintf("drop table if exists %s;alter table %s_tmp rename %s", starrocksTable, starrocksTable, starrocksTable))
		if err != nil {
			return err
		}
		c.GetLogger().Info("load %s to starrocks success", table)
	}
	return err
}

var (
	LoadDataTaskMeta = core.SubTaskMeta{
		Name:             "LoadData",
		EntryPoint:       LoadData,
		EnabledByDefault: true,
		Description:      "Load data to StarRocks",
	}
)
