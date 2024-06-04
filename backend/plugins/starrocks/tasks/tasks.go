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

package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"github.com/apache/incubator-devlake/plugins/starrocks/utils"

	"github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Table struct {
	name string
}

func (t *Table) TableName() string {
	return t.name
}

type DataConfigParams struct {
	Ctx           plugin.SubTaskContext
	Config        *StarRocksConfig
	SrcDb         dal.Dal
	DestDb        dal.Dal
	SrcTableName  string
	DestTableName string
}

func ExportData(c plugin.SubTaskContext) errors.Error {
	logger := c.GetLogger()
	config := c.GetData().(*StarRocksConfig)

	// 1. Get db instance
	var db dal.Dal
	if config.SourceDsn != "" && config.SourceType != "" {
		o, err := getDbInstance(c)
		if err != nil {
			return errors.Convert(err)
		}
		db = dalgorm.NewDalgorm(o)
		sqlDB, err := o.DB()
		if err != nil {
			return errors.Convert(err)
		}
		defer sqlDB.Close()
	} else {
		db = c.GetDal()
	}

	// 2. Filter out the tables to export
	starrocksTables, err := getExportingTables(c, db)
	if err != nil {
		return errors.Convert(err)
	}
	// 3. copy devlake data to starrocks
	sr, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Password, config.Host, config.Port, config.Database)))
	if err != nil {
		return errors.Convert(err)
	}
	starrocksDb := dalgorm.NewDalgorm(sr)
	sqlStarrocksDB, err := sr.DB()
	if err != nil {
		return errors.Convert(err)
	}
	defer sqlStarrocksDB.Close()

	for _, table := range starrocksTables {
		select {
		case <-c.GetContext().Done():
			return errors.Convert(c.GetContext().Err())
		default:
		}

		dc := DataConfigParams{
			Ctx:           c,
			Config:        config,
			SrcDb:         db,
			DestDb:        starrocksDb,
			SrcTableName:  table,
			DestTableName: strings.TrimLeft(table, "_"),
		}
		columnMap, orderBy, skip, err := createTmpTableInStarrocks(&dc)
		if skip {
			logger.Info(fmt.Sprintf("table %s is up to date, so skip it", table))
			continue
		}
		if err != nil {
			logger.Error(err, "create table %s in starrocks error", table)
			return errors.Convert(err)
		}
		err = copyDataToDst(&dc, columnMap, orderBy)
		if err != nil {
			return errors.Convert(err)
		}
	}
	return nil
}

// create temp table for dealing with some complex logic
func createTmpTableInStarrocks(dc *DataConfigParams) (map[string]string, string, bool, error) {
	logger := dc.Ctx.GetLogger()
	config := dc.Config
	db := dc.SrcDb
	starrocksDb := dc.DestDb
	table := dc.SrcTableName
	starrocksTable := dc.DestTableName
	starrocksTmpTable := fmt.Sprintf("%s_tmp", starrocksTable)

	columnMetas, err := db.GetColumns(&Table{name: table}, nil)
	updateColumn := config.UpdateColumn
	columnMap := make(map[string]string)
	if err != nil {
		if strings.Contains(err.Error(), "cached plan must not change result type") {
			logger.Warn(err, "skip err: cached plan must not change result type")
			columnMetas, err = db.GetColumns(&Table{name: table}, nil)
			if err != nil {
				return nil, "", false, err
			}
		} else {
			return nil, "", false, err
		}
	}

	var pks, orders, columns []string
	var separator, firstcm, firstcmName string
	if db.Dialect() == "postgres" {
		separator = "\""
	} else if db.Dialect() == "mysql" {
		separator = "`"
	} else {
		return nil, "", false, errors.NotFound.New(fmt.Sprintf("unsupported dialect %s", db.Dialect()))
	}
	for _, cm := range columnMetas {
		name := cm.Name()
		if name == updateColumn {
			// check update column to detect skip or not
			var updatedFrom time.Time
			err = db.All(&updatedFrom, dal.Select(updateColumn), dal.From(table), dal.Limit(1), dal.Orderby(fmt.Sprintf("%s desc", updateColumn)))
			if err != nil {
				return nil, "", false, err
			}

			var updatedTo time.Time
			err = starrocksDb.All(&updatedTo, dal.Select(updateColumn), dal.From(starrocksTable), dal.Limit(1), dal.Orderby(fmt.Sprintf("%s desc", updateColumn)))
			if err != nil {
				if !strings.Contains(err.Error(), "Unknown table") {
					return nil, "", false, err
				}
			} else {
				if updatedFrom.Equal(updatedTo) {
					sourceCount, err := db.Count(dal.From(table))
					if err != nil {
						return nil, "", false, err
					}
					starrocksCount, err := starrocksDb.Count(dal.From(starrocksTable))
					if err != nil {
						return nil, "", false, err
					}
					// When updated time is equal but record count is different,
					// need to execute the following process, not returning here
					if sourceCount == starrocksCount {
						return nil, "", true, nil
					}
				}
			}
		}

		columnDatatype, ok := cm.ColumnType()
		if !ok {
			return columnMap, "", false, errors.Default.New(fmt.Sprintf("Get [%s] ColumeType Failed", name))
		}
		dataType := utils.GetStarRocksDataType(columnDatatype)
		columnMap[name] = dataType
		column := fmt.Sprintf("`%s` %s", name, dataType)
		columns = append(columns, column)
		isPrimaryKey, ok := cm.PrimaryKey()
		if isPrimaryKey && ok {
			pks = append(pks, fmt.Sprintf("`%s`", name))
			orders = append(orders, fmt.Sprintf("%s%s%s", separator, name, separator))
		}
		if firstcm == "" {
			firstcm = fmt.Sprintf("`%s`", name)
			firstcmName = fmt.Sprintf("%s%s%s", separator, name, separator)
		}
	}

	if len(pks) == 0 {
		pks = append(pks, firstcm)
	}
	orderBy := strings.Join(orders, ", ")
	if config.OrderBy != nil {
		if v, ok := config.OrderBy[table]; ok {
			orderBy = v
		}
	}
	if orderBy == "" {
		orderBy = firstcmName
	}
	replicationNum := os.Getenv("STARROCKS_REPLICAS_NUM")
	if replicationNum == "" {
		replicationNum = "1"
	}
	extra := fmt.Sprintf(`engine=olap distributed by hash(%s) properties("replication_num" = "%s")`, strings.Join(pks, ", "), replicationNum)
	if config.Extra != nil {
		if v, ok := config.Extra[table]; ok {
			extra = v
		}
	}
	tableSql := fmt.Sprintf("DROP TABLE IF EXISTS %s; CREATE TABLE IF NOT EXISTS `%s` ( %s ) %s", starrocksTmpTable, starrocksTmpTable, strings.Join(columns, ","), extra)
	logger.Debug(tableSql)
	err = starrocksDb.Exec(tableSql)
	return columnMap, orderBy, false, err
}

// put data to final dst database
func copyDataToDst(dc *DataConfigParams, columnMap map[string]string, orderBy string) error {
	c := dc.Ctx
	logger := dc.Ctx.GetLogger()
	config := dc.Config
	db := dc.SrcDb
	starrocksDb := dc.DestDb
	table := dc.SrcTableName
	starrocksTable := dc.DestTableName
	starrocksTmpTable := fmt.Sprintf("%s_tmp", starrocksTable)

	var offset int
	var err error
	var rows dal.Rows
	rows, err = db.Cursor(
		dal.From(table),
		dal.Orderby(orderBy),
	)
	if err != nil {
		if strings.Contains(err.Error(), "cached plan must not change result type") {
			logger.Warn(err, "skip err: cached plan must not change result type")
			rows, err = db.Cursor(
				dal.From(table),
				dal.Orderby(orderBy),
			)
			if err != nil {
				return err
			}
		} else {
			return err
		}

	}
	defer rows.Close()

	var data []map[string]interface{}
	cols, err := (rows).Columns()
	if err != nil {
		return err
	}

	var batchCount int
	for rows.Next() {
		select {
		case <-c.GetContext().Done():
			return c.GetContext().Err()
		default:
		}
		row := make(map[string]interface{})
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			dataType := columnMap[cols[i]]
			if strings.HasPrefix(dataType, "array") {
				var arr []string
				columns[i] = &arr
				columnPointers[i] = pq.Array(&arr)
			} else {
				columnPointers[i] = &columns[i]
			}
		}
		err = rows.Scan(columnPointers...)
		if err != nil {
			return err
		}
		for i, colName := range cols {
			row[colName] = columns[i]
		}
		data = append(data, row)
		batchCount += 1
		if batchCount == config.BatchSize {
			err = putBatchData(c, starrocksTmpTable, table, data, config, offset)
			if err != nil {
				return err
			}
			batchCount = 0
			data = nil
		}
	}
	if batchCount != 0 {
		err = putBatchData(c, starrocksTmpTable, table, data, config, offset)
		if err != nil {
			return err
		}
	}

	// drop old table
	err = starrocksDb.Exec("DROP TABLE IF EXISTS ?", clause.Table{Name: starrocksTable})
	if err != nil {
		return err
	}
	// rename tmp table to old table
	err = starrocksDb.Exec("ALTER TABLE ? RENAME ?", clause.Table{Name: starrocksTmpTable}, clause.Table{Name: starrocksTable})
	if err != nil {
		return err
	}

	// check data count
	sourceCount, err := db.Count(dal.From(table))
	if err != nil {
		return err
	}
	starrocksCount, err := starrocksDb.Count(dal.From(starrocksTable))
	if err != nil {
		return err
	}
	if sourceCount != starrocksCount {
		logger.Warn(nil, "source count %d not equal to starrocks count %d", sourceCount, starrocksCount)
	}
	logger.Info("load %s to starrocks success", table)
	return nil
}

// put batch size data to database
func putBatchData(c plugin.SubTaskContext, starrocksTmpTable, table string, data []map[string]interface{}, config *StarRocksConfig, offset int) error {
	logger := c.GetLogger()
	// insert data to tmp table
	loadURL := fmt.Sprintf("http://%s:%d/api/%s/%s/_stream_load", config.BeHost, config.BePort, config.Database, starrocksTmpTable)
	headers := map[string]string{
		"format":            "json",
		"strip_outer_array": "true",
		"Expect":            "100-continue",
		"ignore_json_size":  "true",
		"Connection":        "close",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(http.MethodPut, loadURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.SetBasicAuth(config.User, config.Password)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var b []byte

	if resp.StatusCode == 307 {
		var location *url.URL
		location, err = resp.Location()
		if err != nil {
			return err
		}
		req, err = http.NewRequest(http.MethodPut, location.String(), bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.SetBasicAuth(config.User, config.Password)
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		respRetry, err := client.Do(req)
		if err != nil {
			return err
		}
		defer respRetry.Body.Close()
		b, err = io.ReadAll(respRetry.Body)
		if err != nil {
			return err
		}
	} else {
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	var result map[string]interface{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error(nil, "[%s]: %s", resp.StatusCode, string(b))
	}
	if result["Status"] != "Success" {
		logger.Error(nil, "load %s failed: %s", table, string(b))
	} else {
		logger.Debug("load %s success: %s, limit: %d, offset: %d", table, b, config.BatchSize, offset)
	}
	return nil
}

// get db instance
func getDbInstance(c plugin.SubTaskContext) (o *gorm.DB, err error) {
	config := c.GetData().(*StarRocksConfig)
	switch config.SourceType {
	case "mysql":
		o, err = gorm.Open(mysql.Open(config.SourceDsn))
		if err != nil {
			return nil, err
		}
	case "postgres":
		o, err = gorm.Open(postgres.Open(config.SourceDsn))
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.NotFound.New(fmt.Sprintf("unsupported source type %s", config.SourceType))
	}

	return o, nil
}

// get imported tables
func getExportingTables(c plugin.SubTaskContext, db dal.Dal) (starrocksTables []string, err error) {
	config := c.GetData().(*StarRocksConfig)
	if config.DomainLayer != "" {
		starrocksTables = utils.GetTablesByDomainLayer(config.DomainLayer)
		if starrocksTables == nil {
			return nil, errors.NotFound.New(fmt.Sprintf("no table found by domain layer: %s", config.DomainLayer))
		}
	} else {
		tables := config.Tables
		allTables, err := db.AllTables()
		if err != nil {
			return nil, err
		}
		if len(tables) == 0 {
			starrocksTables = allTables
		} else {
			for _, table := range allTables {
				for _, r := range tables {
					var ok bool
					ok, err := regexp.Match(r, []byte(table))
					if err != nil {
						return nil, err
					}
					if ok {
						starrocksTables = append(starrocksTables, table)
					}
				}
			}
		}
	}
	return starrocksTables, nil
}

var ExportDataTaskMeta = plugin.SubTaskMeta{
	Name:             "ExportData",
	EntryPoint:       ExportData,
	EnabledByDefault: true,
	Description:      "Load data to StarRocks",
}
