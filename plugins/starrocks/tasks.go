package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strings"
)

func getAllTables(database string, db *gorm.DB) ([]string, error) {
	var tableSql string
	if db.Dialector.Name() == "mysql" {
		tableSql = fmt.Sprintf("select table_name from information_schema.tables where table_schema = '%s'", database)
	} else {
		tableSql = "select table_name from information_schema.tables where table_schema = 'public' and table_name not like '_devlake%'"
	}
	var tables []string
	err := db.Raw(tableSql).Scan(&tables).Error
	if err != nil {
		return nil, err
	}
	return tables, nil
}
func LoadData(c core.SubTaskContext) error {
	config := c.GetData().(*StarRocksConfig)
	db := c.GetDb()
	tables := config.Tables
	var err error
	if len(tables) == 0 {
		tables, err = getAllTables(config.Database, db)
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
func loadData(starrocks *sql.DB, c core.SubTaskContext, table string, db *gorm.DB, config *StarRocksConfig) error {
	var data []map[string]interface{}
	// select data from db
	err := db.Raw(fmt.Sprintf("select * from %s", table)).Scan(&data).Error
	if err != nil {
		return err
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
