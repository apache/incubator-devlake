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

	"gorm.io/datatypes"
)

const BLUEPRINT_MODE_NORMAL = "NORMAL"
const BLUEPRINT_MODE_ADVANCED = "ADVANCED"

type Blueprint struct {
	Name       string         `json:"name" validate:"required"`
	Mode       string         `json:"mode" gorm:"varchar(20)" validate:"required,oneof=NORMAL ADVANCED"`
	Plan       datatypes.JSON `json:"plan"`
	Enable     bool           `json:"enable"`
	CronConfig string         `json:"cronConfig"`
	IsManual   bool           `json:"isManual"`
	Settings   datatypes.JSON `json:"settings"`
	Model
}

func (Blueprint) TableName() string {
	return "_devlake_blueprints"
}

type BlueprintSettings struct {
	Version     string          `json:"version" validate:"required,semver,oneof=1.0.0"`
	Connections json.RawMessage `json:"connections" validate:"required"`
}
