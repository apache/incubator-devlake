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
	"time"

	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/jenkins/models/migrationscripts/archived"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JenkinsBuild20220609 struct {
	common.NoPKModel

	// collected fields
	ConnectionId      uint64    `gorm:"primaryKey"`
	JobName           string    `gorm:"primaryKey;type:varchar(255)"`
	Duration          float64   // build time
	DisplayName       string    `gorm:"type:varchar(255)"` // "#7"
	EstimatedDuration float64   // EstimatedDuration
	Number            int64     `gorm:"primaryKey"`
	Result            string    // Result
	Timestamp         int64     // start time
	StartTime         time.Time // convered by timestamp
	CommitSha         string    `gorm:"type:varchar(255)"`
}

func (JenkinsBuild20220609) TableName() string {
	return "_tool_jenkins_builds_20220609"
}

type JenkinsBuildNew struct {
	common.NoPKModel

	// collected fields
	ConnectionId      uint64    `gorm:"primaryKey"`
	JobName           string    `gorm:"primaryKey;type:varchar(255)"`
	Duration          float64   // build time
	DisplayName       string    `gorm:"type:varchar(255)"` // "#7"
	EstimatedDuration float64   // EstimatedDuration
	Number            int64     `gorm:"primaryKey"`
	Result            string    // Result
	Timestamp         int64     // start time
	StartTime         time.Time // convered by timestamp
	CommitSha         string    `gorm:"type:varchar(255)"`
}

func (JenkinsBuildNew) TableName() string {
	return "_tool_jenkins_builds"
}

type UpdateSchemas20220609 struct{}

func (*UpdateSchemas20220609) Up(ctx context.Context, db *gorm.DB) error {
	cursor, err := db.Model(archived.JenkinsBuild{}).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	// 1. create a temporary table to store unique records
	err = db.Migrator().CreateTable(JenkinsBuild20220609{})
	if err != nil {
		return err
	}
	// 2. dedupe records and insert into the temporary table
	for cursor.Next() {
		//inputRow := archived.RefsIssuesDiffs{}
		inputRow := JenkinsBuild20220609{}
		err := db.ScanRows(cursor, &inputRow)
		if err != nil {
			return err
		}
		err = db.Clauses(clause.OnConflict{UpdateAll: true}).Create(inputRow).Error
		if err != nil {
			return err
		}
	}
	// 3. drop old table
	err = db.Migrator().DropTable(archived.JenkinsBuild{})
	if err != nil {
		return err
	}
	// 4. rename the temporary table to the old table
	err = db.Migrator().RenameTable(JenkinsBuild20220609{}, JenkinsBuildNew{})
	if err != nil {
		return err
	}
	return nil
}

func (*UpdateSchemas20220609) Version() uint64 {
	return 20220609154646
}

func (*UpdateSchemas20220609) Name() string {
	return "Add connectionId column to _tool_jenkins_builds"
}
