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
	"gorm.io/gorm"
)

type JiraStatus struct {
	common.NoPKModel
	ConnectionId   uint64 `gorm:"primaryKey"`
	ID             string `gorm:"primaryKey"`
	Name           string
	Self           string
	StatusCategory string
}

func (JiraStatus) TableName() string {
	return "_tool_jira_statuses"
}

type UpdateSchemas20220614 struct{}

func (*UpdateSchemas20220614) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&JiraStatus{},
	)
}

func (*UpdateSchemas20220614) Version() uint64 {
	return 20220614112900
}

func (*UpdateSchemas20220614) Name() string {
	return "add jira status"
}
