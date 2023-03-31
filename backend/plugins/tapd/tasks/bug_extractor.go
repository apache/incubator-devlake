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
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"strings"
)

var _ plugin.SubTaskEntryPoint = ExtractBugs

var ExtractBugMeta = plugin.SubTaskMeta{
	Name:             "extractBugs",
	EntryPoint:       ExtractBugs,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractBugs(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUG_TABLE)
	db := taskCtx.GetDal()
	statusList := make([]models.TapdBugStatus, 0)
	statusLanguageMap, getStdStatus, err := getDefaultStdStatusMapping(data, db, statusList)
	if err != nil {
		return err
	}
	customStatusMap := getStatusMapping(data)
	stdTypeMappings := getStdTypeMappings(data)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		BatchSize:          100,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var bugBody struct {
				Bug models.TapdBug
			}
			err = errors.Convert(json.Unmarshal(row.Data, &bugBody))
			if err != nil {
				return nil, err
			}
			toolL := bugBody.Bug

			toolL.Status = statusLanguageMap[toolL.Status]
			toolL.ConnectionId = data.Options.ConnectionId
			toolL.Type = "BUG"
			toolL.StdType = stdTypeMappings[toolL.Type]
			if toolL.StdType == "" {
				toolL.StdType = ticket.BUG
			}
			if len(customStatusMap) != 0 {
				toolL.StdStatus = customStatusMap[toolL.Status]
			} else {
				toolL.StdStatus = getStdStatus(toolL.Status)
			}
			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/bugtrace/bugs/view?bug_id=%d", toolL.WorkspaceId, toolL.Id)
			if strings.Contains(toolL.CurrentOwner, ";") {
				toolL.CurrentOwner = strings.Split(toolL.CurrentOwner, ";")[0]
			}
			workSpaceBug := &models.TapdWorkSpaceBug{
				ConnectionId: data.Options.ConnectionId,
				WorkspaceId:  toolL.WorkspaceId,
				BugId:        toolL.Id,
			}
			results := make([]interface{}, 0, 3)
			results = append(results, &toolL, workSpaceBug)
			if toolL.IterationId != 0 {
				iterationBug := &models.TapdIterationBug{
					ConnectionId:   data.Options.ConnectionId,
					IterationId:    toolL.IterationId,
					WorkspaceId:    toolL.WorkspaceId,
					BugId:          toolL.Id,
					ResolutionDate: toolL.Resolved,
					BugCreatedDate: toolL.Created,
				}
				results = append(results, iterationBug)
			}
			if toolL.Label != "" {
				labelList := strings.Split(toolL.Label, "|")
				for _, v := range labelList {
					toolLIssueLabel := &models.TapdBugLabel{
						BugId:     toolL.Id,
						LabelName: v,
					}
					results = append(results, toolLIssueLabel)
				}
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
