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
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type Issue20220524 struct {
	archived.DomainEntity
	CreatorName string `gorm:"type:varchar(255)"`
}

func (Issue20220524) TableName() string {
	return "issues"
}

type updateSchemas20220524 struct{}

func (*updateSchemas20220524) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().AddColumn(&Issue20220524{}, "creator_name")
	if err != nil {
		return err
	}
	return nil
}

func (*updateSchemas20220524) Version() uint64 {
	return 20220524000005
}

func (*updateSchemas20220524) Name() string {
	return "Add creator_name column to Issue"
}
