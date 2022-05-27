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
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type Changelog20220527 struct {
	domainlayer.DomainEntity

	// collected fields
	IssueId     string `gorm:"index;type:varchar(255)"`
	AuthorId    string `gorm:"type:varchar(255)"`
	AuthorName  string `gorm:"type:varchar(255)"`
	FieldId     string `gorm:"type:varchar(255)"`
	FieldName   string `gorm:"type:varchar(255)"`
	FromValue   string
	ToValue     string
	CreatedDate time.Time
}

func (Changelog20220527) TableName() string {
	return "changelogs"
}

type updateSchemas20220527 struct{}

func (*updateSchemas20220527) Up(ctx context.Context, db *gorm.DB) error {

	err := db.Migrator().RenameColumn(archived.Changelog{}, "from", "from_value")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.Changelog{}, "to", "to_value")
	if err != nil {
		return err
	}
	return nil
}

func (*updateSchemas20220527) Version() uint64 {
	return 20220527000005
}

func (*updateSchemas20220527) Name() string {
	return "update `from` and `to` columns to `from_value` and `to_value` at changelogs"
}
