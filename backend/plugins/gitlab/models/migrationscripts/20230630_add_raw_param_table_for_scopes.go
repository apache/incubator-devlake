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
	"fmt"
	"strconv"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addRawParamTableForScope)(nil)

type scope20230630 struct {
	ConnectionId  uint64 `gorm:"primaryKey"`
	GitlabId      string `gorm:"primaryKey"`
	RawDataTable  string `gorm:"column:_raw_data_table"`
	RawDataParams string `gorm:"column:_raw_data_params"`
}

func (scope20230630) TableName() string {
	return "_tool_gitlab_projects"
}

type params20230630 struct {
	ConnectionId uint64
	ProjectId    int
}

type addRawParamTableForScope struct{}

func (script *addRawParamTableForScope) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	return migrationhelper.CopyTableColumns(basicRes,
		scope20230630{}.TableName(),
		scope20230630{}.TableName(),
		func(src *scope20230630) (*scope20230630, errors.Error) {
			src.RawDataTable = "_raw_gitlab_scopes"
			src.RawDataParams = string(errors.Must1(json.Marshal(&params20230630{
				ConnectionId: src.ConnectionId,
				ProjectId:    int(errors.Must1(strconv.ParseInt(src.GitlabId, 10, 64))),
			})))
			updateSet := []dal.DalSet{
				{ColumnName: "_raw_data_table", Value: src.RawDataTable},
				{ColumnName: "_raw_data_params", Value: src.RawDataParams},
			}
			where := dal.Where("id = ?", fmt.Sprintf("gitlab:GitlabProject:%v:%v", src.ConnectionId, src.GitlabId))
			errors.Must(db.UpdateColumns("repos", updateSet, where))
			errors.Must(db.UpdateColumns("boards", updateSet, where))
			errors.Must(db.UpdateColumns("cicd_scopes", updateSet, where))
			errors.Must(db.UpdateColumns("cq_projects", updateSet, where))
			return src, nil
		})
}

func (*addRawParamTableForScope) Version() uint64 {
	return 20230630000001
}

func (script *addRawParamTableForScope) Name() string {
	return "populated _raw_data columns for gitlab projects"
}
