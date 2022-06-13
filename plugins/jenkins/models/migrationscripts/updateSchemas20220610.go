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

	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/jenkins/models/migrationscripts/archived"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// JenkinsJobProps current used jenkins job props
type JenkinsJobProps struct {
	// collected fields
	ConnectionId uint64 `gorm:"primaryKey"`
	Name         string `gorm:"primaryKey;type:varchar(255)"`
	Class        string `gorm:"type:varchar(255)"`
	Color        string `gorm:"type:varchar(255)"`
	Base         string `gorm:"type:varchar(255)"`
}

// JenkinsJob db entity for jenkins job
type JenkinsJob20220610 struct {
	JenkinsJobProps
	common.NoPKModel
}

func (JenkinsJob20220610) TableName() string {
	return "_tool_jenkins_jobs_20220609"
}

// JenkinsJob db entity for jenkins job
type JenkinsJobNew struct {
	JenkinsJobProps
	common.NoPKModel
}

func (JenkinsJobNew) TableName() string {
	return "_tool_jenkins_jobs"
}

type UpdateSchemas20220610 struct{}

func (*UpdateSchemas20220610) Up(ctx context.Context, db *gorm.DB) error {
	cursor, err := db.Model(archived.JenkinsJob{}).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	// 1. create a temporary table to store unique records
	err = db.Migrator().CreateTable(JenkinsJob20220610{})
	if err != nil {
		return err
	}
	// 2. dedupe records and insert into the temporary table
	for cursor.Next() {
		//inputRow := archived.RefsIssuesDiffs{}
		inputRow := JenkinsJob20220610{}
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
	err = db.Migrator().DropTable(archived.JenkinsJob{})
	if err != nil {
		return err
	}
	// 4. rename the temporary table to the old table
	db.Migrator().RenameTable(JenkinsJob20220610{}, JenkinsJobNew{})

	return nil
}

func (*UpdateSchemas20220610) Version() uint64 {
	return 20220610154646
}

func (*UpdateSchemas20220610) Name() string {
	return "Add connectionId column to _tool_jenkins_jobs"
}
