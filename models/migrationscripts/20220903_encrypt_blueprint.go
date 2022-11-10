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
	"encoding/json"
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/gorm"
)

type Blueprint0903Before struct {
	NewPlan     string `json:"plan"`
	NewSettings string `json:"settings"`
}

func (Blueprint0903Before) TableName() string {
	return "_devlake_blueprints"
}

type Blueprint0903Temp struct {
	NewPlan     string `json:"plan"`
	NewSettings string `json:"settings"`
	Name        string `json:"name" validate:"required"`
	Mode        string `json:"mode" gorm:"varchar(20)" validate:"required,oneof=NORMAL ADVANCED"`
	Plan        json.RawMessage
	Enable      bool `json:"enable"`
	//please check this https://crontab.guru/ for detail
	CronConfig     string `json:"cronConfig" format:"* * * * *" example:"0 0 * * 1"`
	IsManual       bool   `json:"isManual"`
	Settings       json.RawMessage
	archived.Model `swaggerignore:"true"`
}

func (Blueprint0903Temp) TableName() string {
	return "_devlake_blueprints"
}

type Blueprint0903TempAfter struct {
	Plan        string
	Settings    string
	NewPlan     string
	NewSettings string
}

func (Blueprint0903TempAfter) TableName() string {
	return "_devlake_blueprints"
}

type encryptBLueprint struct{}

func (*encryptBLueprint) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.AutoMigrate(&Blueprint0903Before{})
	if err != nil {
		return errors.Convert(err)
	}
	var result *gorm.DB
	var blueprintList []Blueprint0903Temp
	result = db.Find(&blueprintList)
	if result.Error != nil {
		return errors.Convert(result.Error)
	}

	// Encrypt all blueprints.plan&settings which had been stored before v0.14
	for _, v := range blueprintList {
		c := config.GetConfig()
		encKey := c.GetString(core.EncodeKeyEnvStr)
		if encKey == "" {
			return errors.BadInput.New("invalid encKey")
		}
		v.NewPlan, err = core.Encrypt(encKey, string(v.Plan))
		if err != nil {
			return errors.Convert(err)
		}
		v.NewSettings, err = core.Encrypt(encKey, string(v.Settings))
		if err != nil {
			return errors.Convert(err)
		}
		err = errors.Convert(db.Save(&v).Error)
		if err != nil {
			return errors.Convert(err)
		}
	}
	err = db.Migrator().DropColumn(&Blueprint0903Temp{}, "plan")
	if err != nil {
		return errors.Convert(err)
	}

	err = db.Migrator().DropColumn(&Blueprint0903Temp{}, "settings")
	if err != nil {
		return errors.Convert(err)
	}
	err = db.Migrator().RenameColumn(&Blueprint0903TempAfter{}, "new_plan", "plan")
	if err != nil {
		return errors.Convert(err)
	}
	err = db.Migrator().RenameColumn(&Blueprint0903TempAfter{}, "new_settings", "settings")
	if err != nil {
		return errors.Convert(err)
	}
	_ = db.Find(&blueprintList)
	return nil
}

func (*encryptBLueprint) Version() uint64 {
	return 20220904142321
}

func (*encryptBLueprint) Name() string {
	return "encrypt Blueprint"
}
