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

	"github.com/apache/incubator-devlake/models/domainlayer"
	"gorm.io/gorm"
)

type Changelog struct {
	domainlayer.DomainEntity

	// collected fields
	IssueId      string `gorm:"index;type:varchar(255)"`
	AuthorId     string `gorm:"type:varchar(255)"`
	AuthorName   string `gorm:"type:varchar(255)"`
	FieldId      string `gorm:"type:varchar(255)"`
	FieldName    string `gorm:"type:varchar(255)"`
	FromValue    string
	StandardFrom string
	StandardTo   string
	ToValue      string
	CreatedDate  time.Time
}

type updateSchemas2022061402 struct{}

func (*updateSchemas2022061402) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(&Changelog{})
}

func (*updateSchemas2022061402) Version() uint64 {
	return 20220614091600
}

func (*updateSchemas2022061402) Name() string {
	return "update table: changelogs, add standard_from, standard_to"
}
