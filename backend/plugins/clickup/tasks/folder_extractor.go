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

var _ plugin.SubTaskEntryPoint = ExtractFolder

func ExtractFolder(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_FOLDER_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,

		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			folder := Folder{}
			err := json.Unmarshal(resData.Data, &folder)
			if err != nil {
				panic(err)
			}
			extractedModels := make([]interface{}, 0)
			extractedModels = append(extractedModels, &models.ClickUpFolder{
				Id:           folder.Id,
				ConnectionId: data.Options.ConnectionId,
				SpaceId:      folder.Space.Id,
				Name:         folder.Name,
			})
			for _, list := range folder.Lists {
				extractedModels = append(extractedModels, &models.ClickUpList{
					Id:           list.Id,
					ConnectionId: data.Options.ConnectionId,
					SpaceId:      folder.Space.Id,
					Name:         list.Name,
					StartDate:    parseDate(list.StartDate),
					DueDate:      parseDate(list.DueDate),
				})
			}
			return extractedModels, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractFolderMeta = plugin.SubTaskMeta{
	Name:             "ExtractFolder",
	EntryPoint:       ExtractFolder,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table clickup_issue",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
