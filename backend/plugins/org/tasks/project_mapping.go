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

package tasks

import (
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/plugin"
)

var SetProjectMappingMeta = plugin.SubTaskMeta{
	Name:             "setProjectMapping",
	EntryPoint:       SetProjectMapping,
	EnabledByDefault: true,
	Description:      "set project mapping",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

// SetProjectMapping binds projects and scopes
func SetProjectMapping(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*TaskData)
	var err errors.Error

	for _, mapping := range data.Options.ProjectMappings {
		err = db.Delete(&crossdomain.ProjectMapping{}, dal.Where("project_name = ?", mapping.ProjectName))
		if err != nil {
			return err
		}
		var projectMappings []crossdomain.ProjectMapping
		for _, scope := range mapping.Scopes {
			projectMappings = append(projectMappings, crossdomain.ProjectMapping{
				ProjectName: mapping.ProjectName,
				Table:       scope.Table,
				RowId:       scope.RowID,
				NoPKModel: common.NoPKModel{
					RawDataOrigin: common.RawDataOrigin{
						// set the RawDataParams equals to projectName. In the case of importing from CSV file, records would be deleted in terms of this field
						RawDataParams: mapping.ProjectName,
					},
				},
			})
		}
		if len(projectMappings) > 0 {
			err = db.CreateOrUpdate(projectMappings)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
