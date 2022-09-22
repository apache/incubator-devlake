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
func ReBuildTable(db gorm.DB, old schema.Tabler, oldBak schema.Tabler, new schema.Tabler, callback_transform func(old schema.Tabler) schema.Tabler) (errs errors.Error) {
	var err error

	// rename the commit_file_bak to cache old table
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
	cursor, err := db.Model(oldBak).Rows()
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error on select [%s]]", oldBak.TableName()))
	}
	defer cursor.Close()

	// caculate and save the data to new table
	batch, err := helper.NewBatchSave(api.BasicRes, reflect.TypeOf(new), 200)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error getting batch from table [%s]", new.TableName()))
	}

	defer batch.Close()
	for cursor.Next() {
		err = db.ScanRows(cursor, oldBak)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("error scan rows from table [%s]", oldBak.TableName()))
		}

		cf := callback_transform(oldBak)

		err = batch.Add(&cf)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("error on [%s] batch add", new.TableName()))
		}
	}

	// drop the old table
	err = db.Migrator().DropTable(oldBak)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error no drop [%s]", oldBak.TableName()))
	}

	return nil
}

// ReBuildTableWithOutBak method can be used when we need to change the table structure and reprocess all the data in the table.
// It request the old table and the new table with different table name.
func ReBuildTableWithOutBak(db gorm.DB, old schema.Tabler, new schema.Tabler, callback_transform func(old schema.Tabler) schema.Tabler) (errs errors.Error) {
	var err error

	if old.TableName() == new.TableName() {
		return errors.Default.New(fmt.Sprintf("you can not use the ReBuildTableWithOutBak with old table and new table have same name [%s][%s].please set a bak table and use ReBuildTable",
			old.TableName(), new.TableName()))
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

	return nil
}
