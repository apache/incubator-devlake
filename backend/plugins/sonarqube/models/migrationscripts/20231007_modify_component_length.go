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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*modifyComponentLength)(nil)

type modifyComponentLength struct{}

type sonarqubeHotspot20231007 struct {
	Component string `gorm:"index;type:varchar(500)"`
}

func (sonarqubeHotspot20231007) TableName() string {
	return "_tool_sonarqube_hotspots"
}

type sonarqubeIssueCodeBlock20231007 struct {
	Component string `gorm:"index;type:varchar(500)"`
}

func (sonarqubeIssueCodeBlock20231007) TableName() string {
	return "_tool_sonarqube_issue_code_blocks"
}

type sonarqubeIssue20231007 struct {
	Component string `gorm:"index;type:varchar(500)"`
}

func (sonarqubeIssue20231007) TableName() string {
	return "_tool_sonarqube_issues"
}

func (script *modifyComponentLength) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := migrationhelper.ChangeColumnsType[sonarqubeHotspot20231007](
		basicRes,
		script,
		sonarqubeHotspot20231007{}.TableName(),
		[]string{"component"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&sonarqubeHotspot20231007{},
				"component",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? != '' ", tmpColumnParams...),
			)
		},
	)
	if err != nil {
		return err
	}

	err = migrationhelper.ChangeColumnsType[sonarqubeIssueCodeBlock20231007](
		basicRes,
		script,
		sonarqubeIssueCodeBlock20231007{}.TableName(),
		[]string{"component"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&sonarqubeIssueCodeBlock20231007{},
				"component",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? != '' ", tmpColumnParams...),
			)
		},
	)
	if err != nil {
		return err
	}

	err = migrationhelper.ChangeColumnsType[sonarqubeIssue20231007](
		basicRes,
		script,
		sonarqubeIssue20231007{}.TableName(),
		[]string{"component"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&sonarqubeIssue20231007{},
				"component",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? != '' ", tmpColumnParams...),
			)
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (*modifyComponentLength) Version() uint64 {
	return 20231007145127
}

func (*modifyComponentLength) Name() string {
	return "modify component type to varchar(500)"
}
