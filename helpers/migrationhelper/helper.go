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

// ReBuildTable method can be used when we need to change the table structure and reprocess all the data in the table.
func ReBuildTable(db *gorm.DB, old schema.Tabler, oldBak schema.Tabler, new schema.Tabler, callback_transform func(old schema.Tabler) schema.Tabler) (errs errors.Error) {
	var err error

	// param cheking
	errs = paramCheckingForReBuildTable(db, old, oldBak, new, callback_transform)
	if errs != nil {
		return errors.Default.Wrap(errs, "ReBuildTable param cheking error")
	}

	// rename the old to oldBak for cache old table
	err = db.Migrator().RenameTable(old, oldBak)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error no rename [%s] to [%s]", old.TableName(), oldBak.TableName()))
	}

	// rollback for rename back
	defer func() {
		if errs != nil {
			err = db.Migrator().RenameTable(oldBak, old)
			if err != nil {
				errs = errors.Default.Wrap(err, fmt.Sprintf("fail to rollback table [%s] , you must to rollback by yourself. %s", oldBak.TableName(), err.Error()))
			}
		}
	}()

	return ReBuildTableWithOutBak(db, oldBak, new, callback_transform)
}

// ReBuildTableWithOutBak method can be used when we need to change the table structure and reprocess all the data in the table.
// It request the old table and the new table with different table name.
func ReBuildTableWithOutBak(db *gorm.DB, old schema.Tabler, new schema.Tabler, callback_transform func(old schema.Tabler) schema.Tabler) (errs errors.Error) {
	var err error

	errs = paramCheckingForReBuildTableWithOutBak(db, old, new, callback_transform)
	if errs != nil {
		return errors.Default.Wrap(errs, "ReBuildTableWithOutBak param cheking error")
	}

	// create new commit_files table
	err = db.Migrator().AutoMigrate(new)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error on auto migrate [%s]", new.TableName()))
	}

	// rollback for create new table
	defer func() {
		if errs != nil {
			err = db.Migrator().DropTable(new)
			if err != nil {
				errs = errors.Default.Wrap(err, fmt.Sprintf("fail to rollback table [%s] , you must to rollback by yourself. %s", new.TableName(), err.Error()))
			}
		}
	}()

	// update old id to new id and write to the new table
	cursor, err := db.Model(old).Rows()
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error on select [%s]]", old.TableName()))
	}
	defer cursor.Close()

	// caculate and save the data to new table
	batch, err := helper.NewBatchSave(api.BasicRes, reflect.TypeOf(new), 200)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error getting batch from table [%s]", new.TableName()))
	}

	defer batch.Close()
	for cursor.Next() {
		err = db.ScanRows(cursor, old)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("error scan rows from table [%s]", old.TableName()))
		}

		cf := callback_transform(old)

		err = batch.Add(&cf)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("error on [%s] batch add", new.TableName()))
		}
	}

	// drop the old table
	err = db.Migrator().DropTable(old)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error no drop [%s]", old.TableName()))
	}

	return nil
}

// paramCheckingForReBuildTable check the params of ReBuildTable
func paramCheckingForReBuildTable(db *gorm.DB, old schema.Tabler, oldBak schema.Tabler, new schema.Tabler, callback_transform func(old schema.Tabler) schema.Tabler) (errs errors.Error) {
	errs = paramCheckingForReBuildTableShare(db, new, callback_transform)
	if errs != nil {
		return errs
	}

	if old == nil {
		return errors.Default.New("can not working with param old nil")
	}

	if oldBak == nil {
		return errors.Default.New("can not working with param oldBack nil")
	}

	if new.TableName() == oldBak.TableName() {
		return errors.Default.New(fmt.Sprintf("the oldBak and new can not use the same table name [%s][%s].",
			oldBak.TableName(), new.TableName()))
	}

	if old.TableName() == oldBak.TableName() {
		return errors.Default.New(fmt.Sprintf("the old and oldBak can not use the same table name [%s][%s].",
			old.TableName(), oldBak.TableName()))
	}

	return nil
}

// paramCheckingForReBuildTableWithOutBak check the params of ReBuildTableWithOutBak
func paramCheckingForReBuildTableWithOutBak(db *gorm.DB, old schema.Tabler, new schema.Tabler, callback_transform func(old schema.Tabler) schema.Tabler) (errs errors.Error) {
	errs = paramCheckingForReBuildTableShare(db, new, callback_transform)
	if errs != nil {
		return errs
	}

	if old == nil {
		return errors.Default.New("can not working with param old nil")
	}

	if old.TableName() == new.TableName() {
		return errors.Default.New(fmt.Sprintf("old and new can not use the same table name [%s][%s].",
			old.TableName(), new.TableName()))
	}

	return nil
}

// paramCheckingForReBuildTable check the Share part params of ReBuildTable and ReBuildTableWithOutBak
func paramCheckingForReBuildTableShare(db *gorm.DB, new schema.Tabler, callback_transform func(old schema.Tabler) schema.Tabler) (errs errors.Error) {
	if db == nil {
		return errors.Default.New("can not working with param db nil")
	}

	if new == nil {
		return errors.Default.New("can not working with param new nil")
	}

	if callback_transform == nil {
		return errors.Default.New("can not working with param callback_transform nil")
	}

	return nil
}
