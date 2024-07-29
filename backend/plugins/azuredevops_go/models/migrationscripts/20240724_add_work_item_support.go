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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"time"
)

type createWorkItemTable struct{}

type createWorkItemAzuredevopsRepo struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	WorkItemID   int    `gorm:"primaryKey"`
	Title        string `gorm:"type:varchar(255)"`
	Type         string `gorm:"type:varchar(255)"`
	State        string `gorm:"type:varchar(255)"`
	CreatedDate  *time.Time
	ResolvedDate *time.Time
	ChangedDate  *time.Time
	CreatorName  string `gorm:"type:varchar(255)"`
	CreatorId    string `gorm:"type:varchar(255)"`
	AssigneeName string `gorm:"type:varchar(255)"`
	Area         string `gorm:"type:varchar(255)"`
	Url          string `gorm:"type:varchar(255)"`
	Severity     string `gorm:"type:varchar(255)"`
	Priority     string `gorm:"type:varchar(255)"`
	StoryPoint   float64
}

func (createWorkItemAzuredevopsRepo) TableName() string {
	return "_tool_azuredevops_go_workitem"
}

func (*createWorkItemTable) Up(baseRes context.BasicRes) errors.Error {
	err := migrationhelper.AutoMigrateTables(
		baseRes,
		&createWorkItemAzuredevopsRepo{},
	)
	return err

}

func (*createWorkItemTable) Version() uint64 {
	return 20240724100000
}

func (*createWorkItemTable) Name() string {
	return "add new table _tool_azuredevops_go_workitem in order to support Azure Work Items"
}
