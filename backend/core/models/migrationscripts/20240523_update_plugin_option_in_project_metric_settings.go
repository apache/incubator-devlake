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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addCommitShaToCicdRelease)(nil)

type projectMetricSettings20240523 struct {
	PluginOption json.RawMessage `gorm:"type:json"`
}

func (projectMetricSettings20240523) TableName() string {
	return "project_metric_settings"
}

type updatePluginOptionInProjectMetricSetting struct{}

func (u *updatePluginOptionInProjectMetricSetting) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := migrationhelper.ChangeColumnsType[projectMetricSettings20240523](
		basicRes,
		u,
		projectMetricSettings20240523{}.TableName(),
		[]string{"plugin_option"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&projectMetricSettings20240523{},
				"plugin_option",
				dal.DalClause{Expr: " ? ", Params: []interface{}{"{}"}},
				dal.Where("? is not null", tmpColumnParams...),
			)
		},
	); err != nil {
		return err

	}
	return nil
}

func (*updatePluginOptionInProjectMetricSetting) Version() uint64 {
	return 20240523194205
}

func (*updatePluginOptionInProjectMetricSetting) Name() string {
	return "update plugin_option's type in project_metric_settings"
}
