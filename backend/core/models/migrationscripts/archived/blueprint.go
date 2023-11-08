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

package archived

import (
	"encoding/json"
)

type Blueprint struct {
	Name       string
	Tasks      json.RawMessage `gorm:"type:json"`
	Enable     bool
	CronConfig string
	Model
}

func (Blueprint) TableName() string {
	return "_devlake_blueprints"
}

type BlueprintConnection struct {
	BlueprintId  uint64 `gorm:"primaryKey"`
	PluginName   string `gorm:"primaryKey;type:varchar(255)"`
	ConnectionId uint64 `gorm:"primaryKey"`
}

func (BlueprintConnection) TableName() string {
	return "_devlake_blueprint_connections"
}

type BlueprintScope struct {
	BlueprintId  uint64 `gorm:"primaryKey"`
	PluginName   string `gorm:"primaryKey;type:varchar(255)"`
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(255)"`
}

func (BlueprintScope) TableName() string {
	return "_devlake_blueprint_scopes"
}
