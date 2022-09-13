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

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
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

type commitfileAddLength struct{}

func (*commitfileAddLength) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.Migrator().AlterColumn(&CommitFileAddLength{}, "file_path")
	if err != nil {
		return err
	}

	// update old id to new id
	cursor, err := db.Model(&CommitFileAddLength{}).Select([]string{"commit_sha", "file_path"}).Rows()
	if err != nil {
		return err
	}

	for cursor.Next() {
		cf := CommitFileAddLength{}
		err = db.ScanRows(cursor, &cf)
		if err != nil {
			return err
		}

		ShaFilePath := sha256.New()
		ShaFilePath.Write([]byte(cf.FilePath))

		err = db.Model(cf).
			Where(`commit_sha=? AND file_path=?`, cf.CommitSha, cf.FilePath).
			Update(`id`, cf.CommitSha+hex.EncodeToString(ShaFilePath.Sum(nil))).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (*commitfileAddLength) Version() uint64 {
	return 20220913165805
}

func (*commitfileAddLength) Name() string {
	return "add length of commit_file file_path"
}
