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

type Job20220721 struct {
	Name string `gorm:"type:varchar(255)"`
	Type string `gorm:"type:varchar(255)"`
	archived.DomainEntity
}

func (Job20220721) TableName() string {
	return "jobs"
}

type addTypeField struct{}

func (*addTypeField) Up(ctx context.Context, db *gorm.DB) error {

	err := db.Migrator().AddColumn(Job20220721{}, "type")
	if err != nil {
		return err
	}

	return nil
}

func (*addTypeField) Version() uint64 {
	return 20220721000005
}

func (*addTypeField) Name() string {
	return "add column `type` at jobs"
}
