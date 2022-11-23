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
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*renameTasksToPlan)(nil)

type blueprint20220622 struct {
	Settings json.RawMessage `json:"settings"`
}

func (blueprint20220622) TableName() string {
	return "_devlake_blueprints"
}

type pipeline20220622 struct {
}

func (pipeline20220622) TableName() string {
	return "_devlake_pipelines"
}

// migration script
type renameTasksToPlan struct{}

func (*renameTasksToPlan) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.AutoMigrate(blueprint20220622{})
	if err != nil {
		return err
	}
	err = db.RenameColumn(blueprint20220622{}.TableName(), "tasks", "plan")
	if err != nil {
		return err
	}
	err = db.RenameColumn(pipeline20220622{}.TableName(), "tasks", "plan")
	if err != nil {
		return err
	}
	return nil
}

func (*renameTasksToPlan) Version() uint64 {
	return 20220622110537
}

func (*renameTasksToPlan) Name() string {
	return "blueprint normal mode support"
}
