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

package migrationhelper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
)

// AutoMigrateTables runs AutoMigrate for muliple tables
func AutoMigrateTables(basicRes core.BasicRes, dst ...interface{}) errors.Error {
	db := basicRes.GetDal()
	for _, entity := range dst {
		err := db.AutoMigrate(entity)
		if err != nil {
			return err
		}
		_ = db.All(entity)
	}
	return nil
}

// ChangeColumnsType change the type of specified columns for the table
func ChangeColumnsType[D any](
	basicRes core.BasicRes,
	script core.MigrationScript,
	tableName string,
	columns []string,
	update func(tmpColumnParams []interface{}) errors.Error,
) (err errors.Error) {
	db := basicRes.GetDal()
	tmpColumnsNames := make([]string, len(columns))
	for i, v := range columns {
		tmpColumnsNames[i] = fmt.Sprintf("%s_%s", v, hashScript(script))
		err = db.RenameColumn(tableName, v, tmpColumnsNames[i])
		if err != nil {
			return err
		}

		defer func(tmpColumnName string, ColumnsName string) {
			if err != nil {
				err1 := db.RenameColumn(tableName, tmpColumnName, ColumnsName)
				if err1 != nil {
					err = errors.Default.Wrap(err, fmt.Sprintf("RollBack by RenameColum failed.Relevant data needs to be repaired manually.%s", err1.Error()))
				}
			}
		}(tmpColumnsNames[i], v)
	}

	err = db.AutoMigrate(new(D), dal.From(tableName))
	if err != nil {
		return errors.Default.Wrap(err, "AutoMigrate for Add Colume Error")
	}

	defer func() {
		if err != nil {
			err1 := db.DropColumns(tableName, columns...)
			if err1 != nil {
				err = errors.Default.Wrap(err, fmt.Sprintf("RollBack by DropColume failed.Relevant data needs to be repaired manually.%s", err1.Error()))
			}
		}
	}()

	if update == nil {
		dalSet := make([]dal.DalSet, 0, len(columns))
		for i, v := range columns {
			dalSet = append(dalSet, dal.DalSet{
				ColumnName: v,
				Value:      dal.ClauseColumn{Name: tmpColumnsNames[i]},
			})
		}
		err = db.UpdateColumns(
			new(D),
			dalSet,
		)
	} else {
		params := make([]interface{}, 0, len(tmpColumnsNames))
		for _, v := range tmpColumnsNames {
			params = append(params, dal.ClauseColumn{Name: v})
		}
		err = update(params)
	}
	if err != nil {
		return err
	}

	err = db.DropColumns(tableName, tmpColumnsNames...)
	if err != nil {
		return err
	}

	return nil
}

// TransformColumns change the type of specified columns for the table and transform data one by one
func TransformColumns[S any, D any](
	basicRes core.BasicRes,
	script core.MigrationScript,
	tableName string,
	columns []string,
	transform func(src *S) (*D, errors.Error),
) (err errors.Error) {
	db := basicRes.GetDal()
	return ChangeColumnsType[D](
		basicRes,
		script,
		tableName,
		columns,
		func(tmpColumnParams []interface{}) errors.Error {
			// create selectStr for transform tmpColumnsNames
			params := make([]interface{}, 0, len(columns)*2)
			selectStr := " * "
			for i, v := range columns {
				selectStr += ",? as ?"
				params = append(params, tmpColumnParams[i])
				params = append(params, dal.ClauseColumn{Name: v})
			}

			cursor, err := db.Cursor(
				dal.Select(selectStr, params...),
				dal.From(dal.ClauseTable{Name: tableName}),
			)
			if err != nil {
				return errors.Default.Wrap(err, fmt.Sprintf("failed to load data from src table [%s]", tableName))
			}

			defer cursor.Close()
			batch, err := helper.NewBatchSave(basicRes, reflect.TypeOf(new(D)), 200, tableName)
			if err != nil {
				return errors.Default.Wrap(err, fmt.Sprintf("failed to instantiate BatchSave for table [%s]", tableName))
			}
			defer batch.Close()
			src := new(S)
			for cursor.Next() {
				err = db.Fetch(cursor, src)
				if err != nil {
					return errors.Default.Wrap(err, fmt.Sprintf("fail to load record from table [%s]", tableName))
				}

				dst, err := transform(src)

				if err != nil {
					return errors.Default.Wrap(err, fmt.Sprintf("failed to update row %v", src))
				}
				err = batch.Add(dst)
				if err != nil {
					return errors.Default.Wrap(err, fmt.Sprintf("push to BatchSave failed %v", dst))
				}
			}
			return nil
		},
	)
}

// TransformTable can be used when we need to change the table structure and reprocess all the data in the table.
func TransformTable[S any, D any](
	basicRes core.BasicRes,
	script core.MigrationScript,
	tableName string,
	transform func(*S) (*D, errors.Error),
) (err errors.Error) {
	db := basicRes.GetDal()
	tmpTableName := fmt.Sprintf("%s_%s", tableName, hashScript(script))

	// rename the src to tmp in case of failure
	err = db.RenameTable(tableName, tmpTableName)
	if err != nil {
		return errors.Default.Wrap(
			err,
			fmt.Sprintf("failed to rename rename src table [%s] to [%s]", tableName, tmpTableName),
		)
	}
	// rollback for error
	defer func() {
		if err != nil {
			err = db.RenameTable(tmpTableName, tableName)
			if err != nil {
				msg := fmt.Sprintf(
					"fail to rollback table [%s] to [%s], you may have to do it manually",
					tmpTableName,
					tableName,
				)
				err = errors.Default.Wrap(err, msg)
			}
		}
	}()

	// create new table with the same name
	err = db.AutoMigrate(new(D), dal.From(tableName))
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error on auto migrate [%s]", tableName))
	}
	// rollback for error
	defer func() {
		if err != nil {
			err = db.DropTables(tableName)
			if err != nil {
				msg := fmt.Sprintf(
					"fail to drop table [%s], you may have to do it manually",
					tableName,
				)
				err = errors.Default.Wrap(err, msg)
			}
		}
	}()

	// transform data from temp table to new table
	cursor, err := db.Cursor(dal.From(tmpTableName))
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("failed to load data from src table [%s]", tmpTableName))
	}
	defer cursor.Close()
	batch, err := helper.NewBatchSave(basicRes, reflect.TypeOf(new(D)), 200, tableName)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("failed to instantiate BatchSave for table [%s]", tableName))
	}
	defer batch.Close()
	src := new(S)
	for cursor.Next() {
		err = db.Fetch(cursor, src)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("fail to load record from table [%s]", tmpTableName))
		}
		dst, err := transform(src)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("failed to transform row %v", src))
		}
		err = batch.Add(dst)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("push to BatchSave failed %v", dst))
		}
	}

	// drop the temp table: we can safely ignore the error because it doesn't matter and there will be nothing we
	// can do in terms of rollback or anything in that nature.
	_ = db.DropTables(tmpTableName)

	return err
}

// CopyTableColumn can copy data from src table to dst table
func CopyTableColumn(
	basicRes core.BasicRes,
	srcTable core.Tabler,
	srcColumn []string,
	dstTable core.Tabler,
	dstColumn []string,
) (err errors.Error) {
	if len(srcColumn) != len(dstColumn) {
		return errors.Default.Wrap(err, "the srcColumn length should be same to the dstColumn length")
	}

	db := basicRes.GetDal()
	cursor, err := db.Cursor(dal.From(srcTable))
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("failed to load data from src table [%s]", srcTable.TableName()))
	}
	defer cursor.Close()
	batch, err := helper.NewBatchSave(basicRes, reflect.TypeOf(dstTable), 200, dstTable.TableName())
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("failed to instantiate BatchSave for table [%s]", dstTable.TableName()))
	}
	defer batch.Close()

	for cursor.Next() {
		err = db.Fetch(cursor, srcTable)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("fail to load record from table [%s]", srcTable.TableName()))
		}

		srcVal := reflect.ValueOf(srcTable).Elem()
		dstVal := reflect.ValueOf(dstTable).Elem()
		for i := 0; i < len(srcColumn); i++ {
			dstVal.FieldByName(dstColumn[i]).Set(srcVal.FieldByName(srcColumn[i]))
		}

		err := db.CreateOrUpdate(dstTable)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("push to CreateOrUpdate failed %v", dstTable.TableName()))
		}
	}

	return err
}

func hashScript(script core.MigrationScript) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(fmt.Sprintf("%s:%v", script.Name(), script.Version())))
	if err != nil {
		panic(err)
	}
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}
