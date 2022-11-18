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

var _ core.MigrationScript = (*addProjectToBluePrint)(nil)

type addProjectToBluePrint struct{}

type blueprint20221117 struct {
	ProjectName string `json:"project_name" gorm:"type:varchar(255)"`
}

func (blueprint20221117) TableName() string {
	return "_devlake_blueprints"
}

func (script *addProjectToBluePrint) Up(basicRes core.BasicRes) errors.Error {
	// Add column `project_name`  with default value "" and false to `_devlake_blueprints`
	db := basicRes.GetDal()
	err := db.AutoMigrate(&blueprint20221117{})
	if err != nil {
		return err
	}
	err = db.UpdateColumn(&blueprint20221117{}, "project_name", "", dal.Where("project_name is null"))
	if err != nil {
		return err
	}
	return nil
}

func (*addProjectToBluePrint) Version() uint64 {
	return 20221117184342
}

func (*addProjectToBluePrint) Name() string {
	return "add project to blueprint"
}
