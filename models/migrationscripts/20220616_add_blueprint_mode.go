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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
)

var _ core.MigrationScript = (*addBlueprintMode)(nil)

type blueprint20220616 struct {
	Mode     string `json:"mode" gorm:"varchar(20)" validate:"required,oneof=NORMAL ADVANCED"`
	IsManual bool   `json:"isManual"`
}

func (blueprint20220616) TableName() string {
	return "_devlake_blueprints"
}

type addBlueprintMode struct{}

func (*addBlueprintMode) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.AutoMigrate(&blueprint20220616{})
	if err != nil {
		return err
	}
	err = db.UpdateColumn(&blueprint20220616{}, "mode", "ADVANCED", dal.Where("mode is null"))
	if err != nil {
		return err
	}
	err = db.UpdateColumn(&blueprint20220616{}, "is_manual", false, dal.Where("is_manual is null"))
	if err != nil {
		return err
	}
	return nil
}

func (*addBlueprintMode) Version() uint64 {
	return 20220616110537
}

func (*addBlueprintMode) Name() string {
	return "add mode field to blueprint"
}
