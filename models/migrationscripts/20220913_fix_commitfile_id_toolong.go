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
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*fixCommitFileIdTooLong)(nil)

type fixCommitFileIdTooLong struct{}

type commitFile20220913Before struct {
	archived.DomainEntity
	CommitSha string `gorm:"type:varchar(40)"`
	FilePath  string `gorm:"type:varchar(255)"` // target field
}

type commitFile20220913After struct {
	archived.DomainEntity
	CommitSha string `gorm:"type:varchar(40)"`
	FilePath  string `gorm:"type:text"` // target field
}

type commitFileComponent20220913 struct {
	archived.NoPKModel
	CommitFileId  string `gorm:"primaryKey;type:varchar(255)"`
	ComponentName string `gorm:"type:varchar(255)"`
}

func (script *fixCommitFileIdTooLong) Up(basicRes core.BasicRes) errors.Error {
	// To recalculate the primary key values for the `commit_files` since
	// we used the `FilePath` as part of the primary key which would exceed
	// the maximum length of the column in some cases.
	// The purpose of this script is to replace the `FilePath` with its sha1
	// migrate main table
	err := migrationhelper.TransformTable(
		basicRes,
		script,
		"commit_files",
		func(s *commitFile20220913Before) (*commitFile20220913After, errors.Error) {
			// copy data
			dst := commitFile20220913After(*s)
			// generate new id with hashed file path to avoid length problem
			shaFilePath := sha256.New()
			shaFilePath.Write([]byte(dst.FilePath))
			dst.Id = dst.CommitSha + ":" + hex.EncodeToString(shaFilePath.Sum(nil))
			return &dst, nil
		},
	)
	if err != nil {
		return err
	}
	// migrate related table
	return migrationhelper.TransformTable(
		basicRes,
		script,
		"commit_files",
		func(s *commitFileComponent20220913) (*commitFileComponent20220913, errors.Error) {
			// copy data
			dst := commitFileComponent20220913(*s)
			// generate new id with hashed file path to avoid length problem
			ids := strings.Split(dst.CommitFileId, ":")

			commitSha := ""
			filePath := ""

			if len(ids) > 0 {
				commitSha = ids[0]
				if len(ids) > 1 {
					for i := 1; i < len(ids); i++ {
						if i > 1 {
							filePath += ":"
						}
						filePath += ids[i]
					}
				}
			}

			// With some long path,the varchar(255) was not enough both ID and file_path
			// So we use the hash to compress the path in ID and add length of file_path.
			shaFilePath := sha256.New()
			shaFilePath.Write([]byte(filePath))
			dst.CommitFileId = commitSha + ":" + hex.EncodeToString(shaFilePath.Sum(nil))
			return &dst, nil
		},
	)
}

func (*fixCommitFileIdTooLong) Version() uint64 {
	return 20220913165805
}

func (*fixCommitFileIdTooLong) Name() string {
	return "add length of commit_file file_path"
}
