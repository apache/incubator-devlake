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
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

func ConvertWorkspace(taskCtx plugin.SubTaskContext) errors.Error {
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*TapdTaskData)
	logger.Info("convert workspace:%d", data.Options.WorkspaceId)
	var workspace models.TapdWorkspace
	err := db.First(&workspace, dal.Where("connection_id = ? AND id = ?", data.Options.ConnectionId, data.Options.WorkspaceId))
	if err != nil {
		return err
	}
	board := &ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: getWorkspaceIdGen().Generate(workspace.ConnectionId, workspace.Id),
		},
		Name: workspace.Name,
		Url:  fmt.Sprintf("%s/%d", "https://tapd.cn", workspace.Id),
	}

	return db.CreateOrUpdate(board)
}

var ConvertWorkspaceMeta = plugin.SubTaskMeta{
	Name:             "convertWorkspace",
	EntryPoint:       ConvertWorkspace,
	EnabledByDefault: true,
	Description:      "convert Tapd workspace",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
