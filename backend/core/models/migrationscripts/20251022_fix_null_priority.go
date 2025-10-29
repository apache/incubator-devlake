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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*fixNullPriority)(nil)

type fixNullPriority struct{}

type blueprint20251022 struct {
	Priority int `json:"priority"`
}

func (blueprint20251022) TableName() string {
	return "_devlake_blueprints"
}

type pipeline20251022 struct {
	Priority int `json:"priority"`
}

func (pipeline20251022) TableName() string {
	return "_devlake_pipelines"
}

func (*fixNullPriority) Up(basicRes context.BasicRes) errors.Error {
	// Set default value 0 for NULL priority values in existing rows
	db := basicRes.GetDal()
	err := db.UpdateColumn(&blueprint20251022{}, "priority", 0, dal.Where("priority IS NULL"))
	if err != nil {
		return err
	}
	err = db.UpdateColumn(&pipeline20251022{}, "priority", 0, dal.Where("priority IS NULL"))
	if err != nil {
		return err
	}
	return nil
}

func (*fixNullPriority) Version() uint64 {
	return 20251022195645
}

func (*fixNullPriority) Name() string {
	return "fix null priority values in blueprints and pipelines"
}
