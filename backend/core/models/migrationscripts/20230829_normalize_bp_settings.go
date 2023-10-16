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
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*normalizeBpSettings)(nil)

type normalizeBpSettings struct{}

type blueprint20230829 struct {
	ID         uint64 `gorm:"primaryKey"`
	Settings   string
	BeforePlan json.RawMessage
	AfterPlan  json.RawMessage
	TimeAfter  *time.Time
}

func (blueprint20230829) TableName() string {
	return "_devlake_blueprints"
}

type blueprintSettings20230829 struct {
	BeforePlan  json.RawMessage       `json:"before_plan"`
	AfterPlan   json.RawMessage       `json:"after_plan"`
	TimeAfter   *archived.Iso8601Time `json:"timeAfter"`
	Connections []struct {
		Plugin       string `json:"plugin"`
		ConnectionID uint64 `json:"connectionId"`
		Scopes       []struct {
			ID string `json:"id"`
		} `json:"scopes"`
	} `json:"connections"`
}

func (script *normalizeBpSettings) Up(basicRes context.BasicRes) errors.Error {
	err := migrationhelper.AutoMigrateTables(
		basicRes,
		new(blueprint20230829),
		new(archived.BlueprintConnection),
		new(archived.BlueprintScope),
	)
	if err != nil {
		return err
	}
	encKey := basicRes.GetConfig(plugin.EncodeKeyEnvStr)
	if encKey == "" {
		return errors.BadInput.New("invalid encKey")
	}
	db := basicRes.GetDal()
	bp := &blueprint20230829{}
	cursor := errors.Must1(db.Cursor(dal.From("_devlake_blueprints"), dal.Where("mode = ?", "NORMAL")))
	defer cursor.Close()

	for cursor.Next() {
		// load row
		errors.Must(db.Fetch(cursor, bp))
		// decrypt and unmarshal settings
		settingsJson := errors.Must1(plugin.Decrypt(encKey, bp.Settings))
		if settingsJson == "" {
			continue
		}
		settings := &blueprintSettings20230829{}
		errors.Must(json.Unmarshal([]byte(settingsJson), settings))
		// update bp fields
		bp.BeforePlan = settings.BeforePlan
		bp.AfterPlan = settings.AfterPlan
		bp.TimeAfter = archived.Iso8601TimeToTime(settings.TimeAfter)
		errors.Must(db.Update(bp))
		// create bp connections and scopes records
		for _, conn := range settings.Connections {
			errors.Must(db.CreateOrUpdate(&archived.BlueprintConnection{
				BlueprintId:  bp.ID,
				PluginName:   conn.Plugin,
				ConnectionId: conn.ConnectionID,
			}))
			for _, scope := range conn.Scopes {
				errors.Must(db.CreateOrUpdate(&archived.BlueprintScope{
					BlueprintId:  bp.ID,
					PluginName:   conn.Plugin,
					ConnectionId: conn.ConnectionID,
					ScopeId:      scope.ID,
				}))
			}
		}
	}
	// drop settings column
	errors.Must(db.DropColumns("_devlake_blueprints", "settings"))
	return nil
}

func (*normalizeBpSettings) Version() uint64 {
	return 20230829145125
}

func (*normalizeBpSettings) Name() string {
	return "normalize bp settings to multiple tables"
}
