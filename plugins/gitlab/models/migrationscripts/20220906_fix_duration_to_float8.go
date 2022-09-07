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
	"gorm.io/gorm"
)

type fixDurationToFloat8 struct{}

type GitlabJob20220906 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	GitlabId     int    `gorm:"primaryKey"`

	Duration  float64 `gorm:"type:text"`
	Duration2 float64 `gorm:"type:float8"`
}

func (GitlabJob20220906) TableName() string {
	return "_tool_gitlab_jobs"
}

func (*fixDurationToFloat8) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().AddColumn(&GitlabJob20220906{}, `duration2`)
	if err != nil {
		return err
	}

	cursor, err := db.Model(&GitlabJob20220906{}).Select([]string{"connection_id", "gitlab_id", "duration"}).Rows()
	for cursor.Next() {
		job := GitlabJob20220906{}
		err = db.ScanRows(cursor, &job)
		if err != nil {
			return err
		}
		err = db.
			Model(job).
			Where(`connection_id=? AND gitlab_id=?`, job.ConnectionId, job.GitlabId).
			Update(`duration2`, job.Duration).
			Error
		if err != nil {
			return err
		}
	}

	err = db.Migrator().DropColumn(&GitlabJob20220906{}, `duration`)
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(&GitlabJob20220906{}, `duration2`, `duration`)
	if err != nil {
		return err
	}
	return nil
}

func (*fixDurationToFloat8) Version() uint64 {
	return 20220906000005
}

func (*fixDurationToFloat8) Name() string {
	return "UpdateSchemas for fixDurationToFloat8"
}
