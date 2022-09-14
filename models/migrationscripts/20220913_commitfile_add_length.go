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

package migrationscripts

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"reflect"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/gitlab/api"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
)

type CommitFileAddLength struct {
	archived.DomainEntity
	CommitSha string `gorm:"type:varchar(40)"`
	FilePath  string `gorm:"type:text"`
	Additions int
	Deletions int
}

func (CommitFileAddLength) TableName() string {
	return "commit_files"
}

type CommitFileAddLengthBak struct {
	archived.DomainEntity
	CommitSha string `gorm:"type:varchar(40)"`
	FilePath  string `gorm:"type:varchar(255)"`
	Additions int
	Deletions int
}

func (CommitFileAddLengthBak) TableName() string {
	return "commit_files_bak"
}

type addCommitFilePathLength struct{}

func (*addCommitFilePathLength) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.Migrator().RenameTable(&CommitFile{}, &CommitFileAddLengthBak{})
	if err != nil {
		return errors.Default.Wrap(err, "error no rename commit_file to commit_files_bak")
	}

	err = db.Migrator().AutoMigrate(&CommitFileAddLength{})
	if err != nil {
		return errors.Default.Wrap(err, "error on auto migrate commit_file")
	}

	// update old id to new id and write to the new table
	cursor, err := db.Model(&CommitFileAddLengthBak{}).Rows()
	if err != nil {
		return errors.Default.Wrap(err, "error on select CommitFileAddLength")
	}

	batch, err := helper.NewBatchSave(api.BasicRes, reflect.TypeOf(&CommitFileAddLength{}), 100)
	if err != nil {
		return errors.Default.Wrap(err, "error getting batch from table")
	}

	defer batch.Close()
	for cursor.Next() {
		cfb := CommitFileAddLengthBak{}
		err = db.ScanRows(cursor, &cfb)
		if err != nil {
			return errors.Default.Wrap(err, "error scan rows from table")
		}

		cf := CommitFileAddLength(cfb)

		ShaFilePath := sha256.New()
		ShaFilePath.Write([]byte(cf.FilePath))
		cf.Id = cf.CommitSha + hex.EncodeToString(ShaFilePath.Sum(nil))

		err = batch.Add(&cf)
		if err != nil {
			return errors.Default.Wrap(err, "error on batch add")
		}
	}

	// drop the old table
	err = db.Migrator().DropTable(&CommitFileAddLengthBak{})
	if err != nil {
		return errors.Default.Wrap(err, "error no drop commit_files_bak")
	}

	return nil
}

func (*addCommitFilePathLength) Version() uint64 {
	return 20220913165805
}

func (*addCommitFilePathLength) Name() string {
	return "add length of commit_file file_path"
}
