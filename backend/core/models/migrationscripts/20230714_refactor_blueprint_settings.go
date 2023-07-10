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
	"encoding/json"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"gorm.io/datatypes"
	"time"
)

var _ plugin.MigrationScript = (*refactorBlueprintSettings)(nil)

type oldBlueprint struct {
	archived.Model
	Settings string `json:"settings" gorm:"serializer:encdec"`
}

func (oldBlueprint) TableName() string {
	return "_devlake_blueprints"
}

type BlueprintSettings struct {
	archived.Model
	BlueprintId uint64
	TimeAfter   *time.Time
	BeforePlan  datatypes.JSON
	AfterPlan   datatypes.JSON
}

func (BlueprintSettings) TableName() string {
	return "_devlake_blueprint_settings"
}

type BlueprintConnection struct {
	archived.Model
	BlueprintId  uint64
	ConnectionId uint64
	SettingId    uint64
}

func (BlueprintConnection) TableName() string {
	return "_devlake_blueprint_connections"
}

type BlueprintScope struct {
	archived.Model
	BlueprintId           uint64
	ScopeId               string
	BlueprintConnectionId uint64 `gorm:"column:blueprint_connection_id;type:varchar(255)"`
	Name                  string
}

func (BlueprintScope) TableName() string {
	return "_devlake_blueprint_scopes"
}

type oldSettings struct {
	Version     string
	TimeAfter   *time.Time
	BeforePlan  json.RawMessage `json:"before_plan"`
	AfterPlan   json.RawMessage `json:"after_plan"`
	Connections []*struct {
		Plugin       string `json:"plugin" validate:"required"`
		ConnectionId uint64 `json:"connectionId" validate:"required"`
		Scopes       []*struct {
			Id       string   `json:"id"`
			Name     string   `json:"name"`
			Entities []string `json:"entities"`
		} `json:"scopes" validate:"required"`
	}
}

type refactorBlueprintSettings struct {
}

func (script *refactorBlueprintSettings) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	var oldBps []*oldBlueprint
	err := db.All(&oldBps)
	if err != nil {
		return err
	}
	for _, oldBp := range oldBps {
		old := oldSettings{}
		err = errors.Convert(json.Unmarshal([]byte(oldBp.Settings), &old))
		if err != nil {
			return err
		}
		settings := BlueprintSettings{
			BlueprintId: oldBp.ID,
			TimeAfter:   old.TimeAfter,
			BeforePlan:  datatypes.JSON(old.BeforePlan),
			AfterPlan:   datatypes.JSON(old.AfterPlan),
		}
		err = db.AutoMigrate(&settings)
		if err != nil {
			return err
		}
		err = db.CreateOrUpdate(&settings)
		if err != nil {
			return err
		}
		for _, connection := range old.Connections {
			bpConnection := BlueprintConnection{
				BlueprintId:  oldBp.ID,
				ConnectionId: connection.ConnectionId,
				SettingId:    settings.ID,
			}
			err = db.AutoMigrate(&bpConnection)
			if err != nil {
				return err
			}
			err = db.CreateOrUpdate(&bpConnection)
			if err != nil {
				return err
			}
			for _, scope := range connection.Scopes {
				bpScope := BlueprintScope{
					BlueprintId:           oldBp.ID,
					ScopeId:               scope.Id,
					BlueprintConnectionId: bpConnection.ConnectionId,
					Name:                  scope.Name,
				}
				err = db.AutoMigrate(&bpScope)
				if err != nil {
					return err
				}
				err = db.CreateOrUpdate(&bpScope)
				if err != nil {
					return err
				}
			}
		}
	}
	err = db.DropColumns(oldBlueprint{}.TableName(), "settings")
	return err
}

func (*refactorBlueprintSettings) Version() uint64 {
	return 20230714000001
}

func (*refactorBlueprintSettings) Name() string {
	return "refactor and normalize blueprint settings"
}
