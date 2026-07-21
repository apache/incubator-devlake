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
	"encoding/json"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

var ExtractListMeta = plugin.SubTaskMeta{
	Name:             "Extract Lists",
	EntryPoint:       ExtractLists,
	EnabledByDefault: true,
	Description:      "Extract raw folder lists into _tool_clickup_lists, classifying sprint lists",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

var _ plugin.SubTaskEntryPoint = ExtractLists

// ClickUpApiList is the subset of a ClickUp list JSON the extractor reads.
type ClickUpApiList struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Archived bool   `json:"archived"`
	Folder   *struct {
		Id string `json:"id"`
	} `json:"folder"`
	Space *struct {
		Id string `json:"id"`
	} `json:"space"`
}

func ExtractLists(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ClickUpTaskData)
	detector, err := newSprintDetector(data.ScopeConfig)
	if err != nil {
		return err
	}
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickUpApiParams{
				ConnectionId: data.Options.ConnectionId,
				FolderId:     data.Options.FolderId,
			},
			Table: RAW_LIST_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			apiList := &ClickUpApiList{}
			if err := errors.Convert(json.Unmarshal(row.Data, apiList)); err != nil {
				return nil, err
			}
			if apiList.Id == "" {
				return nil, nil
			}
			folderId := data.Options.FolderId
			if apiList.Folder != nil && apiList.Folder.Id != "" {
				folderId = apiList.Folder.Id
			}
			list := &models.ClickUpList{
				ConnectionId: data.Options.ConnectionId,
				ListId:       apiList.Id,
				FolderId:     folderId,
				Name:         apiList.Name,
				Archived:     apiList.Archived,
			}
			if apiList.Space != nil {
				list.SpaceId = apiList.Space.Id
			}
			if sprint := detector.detect(apiList.Name); sprint != nil {
				list.IsSprint = true
				list.SprintName = sprint.name
				list.StartDate = sprint.start
				list.EndDate = sprint.end
			}
			return []interface{}{list}, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
