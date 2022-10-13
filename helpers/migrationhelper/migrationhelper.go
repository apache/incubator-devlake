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
	}
	return nil
}

// DropColumns drops multiple columns for specified table
func DropColumns(basicRes core.BasicRes, tableName string, columnNames ...string) errors.Error {
	db := basicRes.GetDal()
	for _, columnName := range columnNames {
		err := db.DropColumn(tableName, columnName)
		if err != nil {
			return err
		}
	}
	return nil
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

func hashScript(script core.MigrationScript) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(fmt.Sprintf("%s:%v", script.Name(), script.Version())))
	if err != nil {
		panic(err)
	}
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}
