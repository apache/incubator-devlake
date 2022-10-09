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

func hashScript(script core.MigrationScript) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(fmt.Sprintf("%s:%v", script.Name(), script.Version())))
	if err != nil {
		panic(err)
	}
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}

// TransformRowsInPlace can be used when we need to change the table structure and reprocess all the data in the table.
func TransformTable[S dal.Tabler, T dal.Tabler, D dal.Tabler](
	basicRes core.BasicRes,
	script core.MigrationScript,
	src S,
	dst D,
	transform func(S) (D, errors.Error),
) (err errors.Error) {

	db := basicRes.GetDal()

	srcTableName := src.TableName()
	tmpTableName := fmt.Sprintf("%s_%s", srcTableName, hashScript(script))
	// rename the src to tmp in case of failure
	err = db.RenameTable(srcTableName, tmpTableName)
	if err != nil {
		return errors.Default.Wrap(
			err,
			fmt.Sprintf("failed to rename rename src table [%s] to [%s]", srcTableName, tmpTableName),
		)
	}

	// rollback for error
	defer func() {
		if err != nil {
			err = db.RenameTable(tmpTableName, srcTableName)
			if err != nil {
				msg := fmt.Sprintf(
					"fail to rollback table [%s] to [%s], you may have to do it manually",
					tmpTableName,
					srcTableName,
				)
				err = errors.Default.Wrap(err, msg)
			}
		}
	}()

	err = MoveTable(basicRes, src, dst, transform, srcTableName, "")
	return
}

// TransformRows can be used when we need to change the table structure and reprocess all the data in the table.
// It request the src table and the dst table with different table name.
func MoveTable[S dal.Tabler, D dal.Tabler](
	basicRes core.BasicRes,
	src S,
	dst D,
	transform func(S) (D, errors.Error),
	tableNames ...string,
) (err errors.Error) {

	db := basicRes.GetDal()
	srcTableName := src.TableName()
	dstTableName := dst.TableName()

	// overwrite src/dst table names optionally
	if len(tableNames) > 0 {
		srcTableName = tableNames[0]
	}
	if len(tableNames) > 1 {
		dstTableName = tableNames[1]
	}

	if srcTableName == dstTableName {
		err = errors.Default.New(fmt.Sprintf("src and dst are the same table %s", srcTableName))
		return
	}

	// create new table
	err = db.AutoMigrate(dst)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error on auto migrate [%s]", dst.TableName()))
	}

	// drop newly created table if any error
	defer func() {
		if err != nil {
			err = db.DropTables(dstTableName)
			if err != nil {
				msg := fmt.Sprintf(
					"fail to drop table [%s], you may have to do it manually",
					dstTableName,
				)
				err = errors.Default.Wrap(err, msg)
			}
		}
	}()

	// update src id to dst id and write to the dst table
	cursor, err := db.Cursor(
		dal.From(srcTableName),
	)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("failed to load data from src table [%s]", srcTableName))
	}
	defer cursor.Close()

	// caculate and save the data to new table
	batch, err := helper.NewBatchSave(basicRes, reflect.TypeOf(dst), 200, dstTableName)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("failed to instantiate BatchSave for table [%s]", dstTableName))
	}

	defer batch.Close()
	for cursor.Next() {
		err = db.Fetch(cursor, src)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("fail to load record from table [%s]", srcTableName))
		}

		cf, err := transform(src)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("failed to transform row %v", src))
		}

		err = batch.Add(cf)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("push to BatchSave failed %v", cf))
		}
	}

	// drop the src table
	err = db.DropTables(srcTableName)
	if err != nil {
		err = errors.Default.Wrap(err, fmt.Sprintf("fail to drop src table %s", srcTableName))
	}
	return
}
