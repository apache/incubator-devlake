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

package azuredevops

import (
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*AddRawDataForScope)(nil)

type azureDevopsGitRepositories20230714 struct {
	ConnectionId  uint64 `gorm:"primaryKey"`
	Id            string `gorm:"primaryKey"`
	RawDataTable  string `gorm:"column:_raw_data_table"`
	RawDataParams string `gorm:"column:_raw_data_params"`
}

func (azureDevopsGitRepositories20230714) TableName() string {
	return "_tool_azuredevops_gitrepositories"
}

type rawDataParams20230714 struct {
	ConnectionId uint64
	ScopeId      string
}

type AddRawDataForScope struct{}

func (script *AddRawDataForScope) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	return migrationhelper.CopyTableColumns(basicRes,
		azureDevopsGitRepositories20230714{}.TableName(),
		azureDevopsGitRepositories20230714{}.TableName(),
		func(src *azureDevopsGitRepositories20230714) (*azureDevopsGitRepositories20230714, errors.Error) {
			src.RawDataTable = "_raw_azuredevops_scopes"
			src.RawDataParams = string(errors.Must1(json.Marshal(&rawDataParams20230714{
				ConnectionId: src.ConnectionId,
				ScopeId:      src.Id,
			})))
			updateSet := []dal.DalSet{
				{ColumnName: "_raw_data_table", Value: src.RawDataTable},
				{ColumnName: "_raw_data_params", Value: src.RawDataParams},
			}
			where := dal.Where("id = ?", fmt.Sprintf("azuredevops:GitRepository:%v:%v", src.ConnectionId, src.Id))
			errors.Must(db.UpdateColumns("repos", updateSet, where))
			errors.Must(db.UpdateColumns("boards", updateSet, where))
			errors.Must(db.UpdateColumns("cicd_scopes", updateSet, where))
			errors.Must(db.UpdateColumns("cq_projects", updateSet, where))
			return src, nil
		})
}

func (*AddRawDataForScope) Version() uint64 {
	return 20230714000001
}

func (script *AddRawDataForScope) Name() string {
	return "populated _raw_data columns for azuredevops"
}
