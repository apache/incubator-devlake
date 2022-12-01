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

var _ core.MigrationScript = (*addProjectTables)(nil)

type addEnableToProjectMetric struct{}

type ProjectMetric20221201 struct {
	Enable bool `gorm:"type:boolean"`
}

func (ProjectMetric20221201) TableName() string {
	return "project_metrics"
}

func (script *addEnableToProjectMetric) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.AutoMigrate(&ProjectMetric20221201{})
	if err != nil {
		return err
	}
	err = db.UpdateColumn(
		&ProjectMetric20221201{},
		`enable`,
		true,
		dal.Where("enable is null"),
	)
	if err != nil {
		return err
	}
	return err
}

func (*addEnableToProjectMetric) Version() uint64 {
	return 20221201190341
}

func (*addEnableToProjectMetric) Name() string {
	return "add enable to project metric"
}
