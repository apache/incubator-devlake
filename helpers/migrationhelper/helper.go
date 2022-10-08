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
	"fmt"
	"reflect"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/gitlab/api"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// TransformRowsInPlace method can be used when we need to change the table structure and reprocess all the data in the table.
func TransformRowsInPlace(db *gorm.DB, src schema.Tabler, bak schema.Tabler, dst schema.Tabler, callback_transform func(src schema.Tabler) schema.Tabler) (errs errors.Error) {
	var err error

	// param cheking
	errs = paramCheckingForTransformRowsInPlace(db, src, bak, dst, callback_transform)
	if errs != nil {
		return errors.Default.Wrap(errs, "TransformRowsInPlace param cheking error")
	}

	// rename the src to bak for cache src table
	err = db.Migrator().RenameTable(src, bak)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error no rename [%s] to [%s]", src.TableName(), bak.TableName()))
	}

	// rollback for rename back
	defer func() {
		if errs != nil {
			err = db.Migrator().RenameTable(bak, src)
			if err != nil {
				errs = errors.Default.Wrap(err, fmt.Sprintf("fail to rollback table [%s] , you must to rollback by yourself. %s", bak.TableName(), err.Error()))
			}
		}
	}()

	return TransformRowsBetweenTables(db, bak, dst, callback_transform)
}

// TransformRowsBetweenTables method can be used when we need to change the table structure and reprocess all the data in the table.
// It request the src table and the dst table with different table name.
func TransformRowsBetweenTables(db *gorm.DB, src schema.Tabler, dst schema.Tabler, callback_transform func(src schema.Tabler) schema.Tabler) (errs errors.Error) {
	var err error

	errs = paramCheckingForTransformRowsBetweenTables(db, src, dst, callback_transform)
	if errs != nil {
		return errors.Default.Wrap(errs, "TransformRowsBetweenTables param cheking error")
	}

	// create new commit_files table
	err = db.Migrator().AutoMigrate(dst)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error on auto migrate [%s]", dst.TableName()))
	}

	// rollback for create new table
	defer func() {
		if errs != nil {
			err = db.Migrator().DropTable(dst)
			if err != nil {
				errs = errors.Default.Wrap(err, fmt.Sprintf("fail to rollback table [%s] , you must to rollback by yourself. %s", dst.TableName(), err.Error()))
			}
		}
	}()

	// update src id to dst id and write to the dst table
	cursor, err := db.Model(src).Rows()
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error on select [%s]]", src.TableName()))
	}
	defer cursor.Close()

	// caculate and save the data to new table
	batch, err := helper.NewBatchSave(api.BasicRes, reflect.TypeOf(dst), 200)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error getting batch from table [%s]", dst.TableName()))
	}

	defer batch.Close()
	for cursor.Next() {
		err = db.ScanRows(cursor, src)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("error scan rows from table [%s]", src.TableName()))
		}

		cf := callback_transform(src)

		err = batch.Add(&cf)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("error on [%s] batch add", dst.TableName()))
		}
	}

	// drop the src table
	err = db.Migrator().DropTable(src)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error no drop [%s]", src.TableName()))
	}

	return nil
}

// paramCheckingForTransformRowsInPlace check the params of TransformRowsInPlace
func paramCheckingForTransformRowsInPlace(db *gorm.DB, src schema.Tabler, bak schema.Tabler, dst schema.Tabler, callback_transform func(src schema.Tabler) schema.Tabler) (errs errors.Error) {
	errs = paramCheckingShare(db, dst, callback_transform)
	if errs != nil {
		return errs
	}

	if src == nil {
		return errors.Default.New("can not working with param src nil")
	}

	if bak == nil {
		return errors.Default.New("can not working with param bak nil")
	}

	if dst.TableName() == bak.TableName() {
		return errors.Default.New(fmt.Sprintf("the bak and dst can not use the same table name [%s][%s].",
			bak.TableName(), dst.TableName()))
	}

	if src.TableName() == bak.TableName() {
		return errors.Default.New(fmt.Sprintf("the src and bak can not use the same table name [%s][%s].",
			src.TableName(), bak.TableName()))
	}

	return nil
}

// paramCheckingForTransformRowsBetweenTables check the params of ReBuildTableWithOutBak
func paramCheckingForTransformRowsBetweenTables(db *gorm.DB, src schema.Tabler, dst schema.Tabler, callback_transform func(src schema.Tabler) schema.Tabler) (errs errors.Error) {
	errs = paramCheckingShare(db, dst, callback_transform)
	if errs != nil {
		return errs
	}

	if src == nil {
		return errors.Default.New("can not working with param src nil")
	}

	if src.TableName() == dst.TableName() {
		return errors.Default.New(fmt.Sprintf("src and dst can not use the same table name [%s][%s].",
			src.TableName(), dst.TableName()))
	}

	return nil
}

// paramCheckingShare check the Share part params of TransformRowsBetweenTables and TransformRowsInPlace
func paramCheckingShare(db *gorm.DB, dst schema.Tabler, callback_transform func(src schema.Tabler) schema.Tabler) (errs errors.Error) {
	if db == nil {
		return errors.Default.New("can not working with param db nil")
	}

	if dst == nil {
		return errors.Default.New("can not working with param dst nil")
	}

	if callback_transform == nil {
		return errors.Default.New("can not working with param callback_transform nil")
	}

	return nil
}
