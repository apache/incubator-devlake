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
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*createQaTables)(nil)

type createQaTables struct{}

func (*createQaTables) Version() uint64 {
	return 20250421104500 // YYYYMMDDHHMMSS format
}

func (*createQaTables) Name() string {
	return "create QA tables"
}

// QaProject represents a QA project in the domain layer
type QaProject struct {
	archived.DomainEntityExtended
	Name string `gorm:"type:varchar(255);comment:Project name"`
}

// QaApi represents a QA API in the domain layer
type QaApi struct {
	archived.DomainEntityExtended
	Name        string    `gorm:"type:varchar(255);comment:API name"`
	Path        string    `gorm:"type:varchar(255);comment:API path"`
	Method      string    `gorm:"type:varchar(255);comment:API method"`
	CreateTime  time.Time `gorm:"comment:API creation time"`
	CreatorId   string    `gorm:"type:varchar(255);comment:Creator ID"`
	QaProjectId string    `gorm:"type:varchar(255);index;comment:Project ID"`
}

// QaTestCase represents a QA test case in the domain layer
type QaTestCase struct {
	archived.DomainEntityExtended
	Name        string    `gorm:"type:varchar(255);comment:Test case name"`
	CreateTime  time.Time `gorm:"comment:Test case creation time"`
	CreatorId   string    `gorm:"type:varchar(255);comment:Creator ID"`
	Type        string    `gorm:"type:varchar(255);comment:Test case type | functional | api"`                // enum in image, using string
	QaApiId     string    `gorm:"type:varchar(255);comment:Valid only when type = api, represents qa_api_id"` // nullable in image, using string
	QaProjectId string    `gorm:"type:varchar(255);index;comment:Project ID"`
}

// QaTestCaseExecution represents a QA test case execution in the domain layer
type QaTestCaseExecution struct {
	archived.DomainEntityExtended
	QaProjectId  string    `gorm:"type:varchar(255);index;comment:Project ID"`
	QaTestCaseId string    `gorm:"type:varchar(255);index;comment:Test case ID"`
	CreateTime   time.Time `gorm:"comment:Test (plan) creation time"`
	StartTime    time.Time `gorm:"comment:Test start time"`
	FinishTime   time.Time `gorm:"comment:Test finish time"`
	CreatorId    string    `gorm:"type:varchar(255);comment:Executor ID"`
	Status       string    `gorm:"type:varchar(255);comment:Test execution status | PENDING | IN_PROGRESS | SUCCESS | FAILED"` // enum, using string
}

func (*createQaTables) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	if err := db.AutoMigrate(&QaProject{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&QaApi{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&QaTestCase{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&QaTestCaseExecution{}); err != nil {
		return err
	}

	return nil
}
