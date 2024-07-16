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
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

var _ plugin.SubTaskEntryPoint = ExtractEpics

var ExtractEpicsMeta = plugin.SubTaskMeta{
	Name:             "extractEpics",
	EntryPoint:       ExtractEpics,
	EnabledByDefault: true,
	Description:      "extract Jira epics from all boards",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET, plugin.DOMAIN_TYPE_CROSS},
}

func ExtractEpics(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("extract external epic Issues, connection_id=%d, board_id=%d", connectionId, boardId)
	mappings, err := getTypeMappings(data, db)
	if err != nil {
		return err
	}
	userFieldMap, err := getUserFieldMap(db, connectionId, logger)
	if err != nil {
		return err
	}
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_EPIC_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			return extractIssues(data, mappings, row, userFieldMap)
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}

func getIssueFieldMap(db dal.Dal, connectionId uint64, logger log.Logger) (map[string]models.JiraIssueField, errors.Error) {
	var allIssueFields []models.JiraIssueField
	if err := db.All(&allIssueFields, dal.Where("connection_id = ?", connectionId)); err != nil {
		return nil, err
	}
	issueFieldMap := make(map[string]models.JiraIssueField)
	for _, v := range allIssueFields {
		if _, ok := issueFieldMap[v.Name]; ok {
			logger.Warn(nil, "filed name %s is duplicated", v.Name)
			if v.SchemaType == "user" {
				issueFieldMap[v.Name] = v
			}
		} else {
			issueFieldMap[v.Name] = v
		}
	}
	return issueFieldMap, nil
}

func getUserFieldMap(db dal.Dal, connectionId uint64, logger log.Logger) (map[string]struct{}, errors.Error) {
	userFieldMap := make(map[string]struct{})
	issueFieldMap, err := getIssueFieldMap(db, connectionId, logger)
	if err != nil {
		return nil, err
	}
	for filedName, issueField := range issueFieldMap {
		if issueField.SchemaType == "user" {
			userFieldMap[filedName] = struct{}{}
		}
	}
	return userFieldMap, nil
}
