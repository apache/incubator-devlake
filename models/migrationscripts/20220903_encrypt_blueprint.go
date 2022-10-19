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

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*encryptBlueprint)(nil)

type encryptBlueprint struct{}

type Blueprint20220903Before struct {
	Name           string          `json:"name" validate:"required"`
	Mode           string          `json:"mode" gorm:"varchar(20)" validate:"required,oneof=NORMAL ADVANCED"`
	Plan           json.RawMessage `json:"plan"`
	Enable         bool            `json:"enable"`
	CronConfig     string          `json:"cronConfig" format:"* * * * *" example:"0 0 * * 1"`
	IsManual       bool            `json:"isManual"`
	Settings       json.RawMessage `json:"settings" swaggertype:"array,string" example:"please check api: /blueprints/<PLUGIN_NAME>/blueprint-setting"`
	archived.Model `swaggerignore:"true"`
}

func (Blueprint20220903Before) TableName() string {
	return "_devlake_blueprints"
}

type Blueprint20220903After struct {
	/* unchanged part */
	Name           string `json:"name" validate:"required"`
	Mode           string `json:"mode" gorm:"varchar(20)" validate:"required,oneof=NORMAL ADVANCED"`
	Enable         bool   `json:"enable"`
	CronConfig     string `json:"cronConfig" format:"* * * * *" example:"0 0 * * 1"`
	IsManual       bool   `json:"isManual"`
	archived.Model `swaggerignore:"true"`
	/* changed part */
	Plan     string `json:"plan"`
	Settings string `json:"settings"`
}

func (Blueprint20220903After) TableName() string {
	return "_devlake_blueprints"
}

func (script *encryptBlueprint) Up(basicRes core.BasicRes) errors.Error {
	encKey := basicRes.GetConfig(core.EncodeKeyEnvStr)
	if encKey == "" {
		return errors.BadInput.New("invalid encKey")
	}

	return migrationhelper.TransformTable(
		basicRes,
		script,
		"_devlake_blueprints",
		func(s *Blueprint20220903Before) (*Blueprint20220903After, errors.Error) {
			encryptedPlan, err := core.Encrypt(encKey, string(s.Plan))
			if err != nil {
				return nil, err
			}
			encryptedSettings, err := core.Encrypt(encKey, string(s.Settings))
			if err != nil {
				return nil, err
			}

			dst := &Blueprint20220903After{
				Name:       s.Name,
				Mode:       s.Mode,
				Enable:     s.Enable,
				CronConfig: s.CronConfig,
				IsManual:   s.IsManual,
				Model:      archived.Model{ID: s.ID},
				Plan:       encryptedPlan,
				Settings:   encryptedSettings,
			}
			return dst, nil
		},
	)
}

func (*encryptBlueprint) Version() uint64 {
	return 20220904142321
}

func (*encryptBlueprint) Name() string {
	return "encrypt Blueprint"
}
