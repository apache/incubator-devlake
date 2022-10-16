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

var _ core.MigrationScript = (*modifyLeadTimeMinutes)(nil)

type modifyLeadTimeMinutes struct{}

type Issues20220929 struct {
	LeadTimeMinutes int64
}

func (Issues20220929) TableName() string {
	return "issues"
}

func (*modifyLeadTimeMinutes) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	bakColumnName := "lead_time_minutes_20220929"
	err := db.RenameColumn("issues", "lead_time_minutes", bakColumnName)
	defer func() {
		if err != nil {
			_ = db.RenameColumn("issues", bakColumnName, "lead_time_minutes")
		}
	}()
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Issues20220929{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = db.DropColumns("issues", "lead_time_minutes")
		}
	}()
	err = db.UpdateColumn(
		&Issues20220929{},
		"lead_time_minutes",
		dal.DalClause{Expr: bakColumnName},
		dal.Where("lead_time_minutes != 0"),
	)
	if err != nil {
		return err
	}
	err = db.DropColumns("issues", bakColumnName)
	if err != nil {
		return err
	}
	_ = db.First(&Issues20220929{})
	return nil
}

func (*modifyLeadTimeMinutes) Version() uint64 {
	return 20220929145125
}

func (*modifyLeadTimeMinutes) Name() string {
	return "modify lead_time_minutes"
}
