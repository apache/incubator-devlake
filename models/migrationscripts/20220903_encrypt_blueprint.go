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
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/gorm"
)

type BlueprintEncryption0904 struct {
	archived.Model
	Plan        string
	RawPlan     string
	Settings    string
	RawSettings string
}

func (BlueprintEncryption0904) TableName() string {
	return "_devlake_blueprints"
}

type encryptBLueprint struct{}

func (*encryptBLueprint) Up(ctx context.Context, db *gorm.DB) errors.Error {
	c := config.GetConfig()
	encKey := c.GetString(core.EncodeKeyEnvStr)
	if encKey == "" {
		return errors.BadInput.New("invalid encKey")
	}

	blueprint := &BlueprintEncryption0904{}
	err := db.Migrator().RenameColumn(blueprint, "plan", "raw_plan")
	if err != nil {
		return errors.Convert(err)
	}
	err = db.Migrator().RenameColumn(blueprint, "settings", "raw_settings")
	if err != nil {
		return errors.Convert(err)
	}
	err = db.AutoMigrate(blueprint)
	if err != nil {
		return errors.Convert(err)
	}

	// Encrypt all blueprints.plan which had been stored before v0.14
	cursor, err := db.Model(blueprint).Rows()
	if err != nil {
		return errors.Convert(err)
	}
	defer cursor.Close()

	for cursor.Next() {
		err = db.ScanRows(cursor, blueprint)
		if err != nil {
			return errors.Convert(err)
		}
		blueprint.Plan, err = core.Encrypt(encKey, blueprint.RawPlan)
		if err != nil {
			return errors.Convert(err)
		}
		blueprint.Settings, err = core.Encrypt(encKey, blueprint.RawSettings)
		if err != nil {
			return errors.Convert(err)
		}
		err = errors.Convert(db.Save(blueprint).Error)
		if err != nil {
			return errors.Convert(err)
		}
	}

	err = db.Migrator().DropColumn(blueprint, "raw_plan")
	if err != nil {
		return errors.Convert(err)
	}
	err = db.Migrator().DropColumn(blueprint, "raw_settings")
	if err != nil {
		return errors.Convert(err)
	}
	_ = db.First(blueprint)
	return nil
}

func (*encryptBLueprint) Version() uint64 {
	return 20220904142321
}

func (*encryptBLueprint) Name() string {
	return "encrypt Blueprint"
}
