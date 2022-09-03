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
	"github.com/apache/incubator-devlake/errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// model blueprint
type blueprintNormalMode_Blueprint struct {
	Settings datatypes.JSON `json:"settings"`
}

func (blueprintNormalMode_Blueprint) TableName() string {
	return "_devlake_blueprints"
}

// model pipeline
type blueprintNormalMode_Pipeline struct {
}

func (blueprintNormalMode_Pipeline) TableName() string {
	return "_devlake_pipelines"
}

// migration script
type renameTasksToPlan struct{}

func (*renameTasksToPlan) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.Migrator().AddColumn(blueprintNormalMode_Blueprint{}, "settings")
	if err != nil {
		return errors.Convert(err)
	}
	err = db.Migrator().RenameColumn(&blueprintNormalMode_Blueprint{}, "tasks", "plan")
	if err != nil {
		return errors.Convert(err)
	}
	err = db.Migrator().RenameColumn(&blueprintNormalMode_Pipeline{}, "tasks", "plan")
	if err != nil {
		return errors.Convert(err)
	}
	return nil
}

func (*renameTasksToPlan) Version() uint64 {
	return 20220622110537
}

func (*renameTasksToPlan) Name() string {
	return "blueprint normal mode support"
}
