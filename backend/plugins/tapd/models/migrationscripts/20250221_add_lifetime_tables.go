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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type TapdLifeTime struct {
	ConnectionId uint64     `gorm:"primaryKey"`
	Id           uint64     `gorm:"primaryKey;type:BIGINT NOT NULL;autoIncrement:false" json:"id,string"`
	WorkspaceId  uint64     `json:"workspace_id,string"`
	EntityType   string     `json:"entity_type" gorm:"type:varchar(255)"`
	EntityId     uint64     `json:"entity_id,string"`
	Status       string     `json:"status" gorm:"type:varchar(255)"`
	Owner        string     `json:"owner" gorm:"type:varchar(255)"`
	BeginDate    *time.Time `json:"begin_date"`
	EndDate      *time.Time `json:"end_date"`
	TimeCost     float64    `json:"time_cost"`
	Created      *time.Time `json:"created"`
	Operator     string     `json:"operator" gorm:"type:varchar(255)"`
	IsRepeated   int        `json:"is_repeated"`
	ChangeFrom   string     `json:"change_from" gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (TapdLifeTime) TableName() string {
	return "_tool_tapd_life_times"
}

var _ plugin.MigrationScript = (*addLifetimeTables)(nil)

type addLifetimeTables struct{}

func (*addLifetimeTables) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &TapdLifeTime{})
}

func (*addLifetimeTables) Version() uint64 {
	return 20250221000000
}

func (*addLifetimeTables) Name() string {
	return "add tapd lifetime tables"
}
