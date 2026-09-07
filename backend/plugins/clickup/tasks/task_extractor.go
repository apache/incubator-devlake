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
	"strconv"
	"strings"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/clickup/models"
)

var ExtractTaskMeta = plugin.SubTaskMeta{
	Name:             "Extract Tasks",
	EntryPoint:       ExtractTasks,
	EnabledByDefault: true,
	Description:      "Extract raw task data into the tool layer table _tool_clickup_tasks",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

var _ plugin.SubTaskEntryPoint = ExtractTasks

// ClickUpApiTask is the subset of the ClickUp task JSON that the extractor reads.
type ClickUpApiTask struct {
	Id          string `json:"id"`
	CustomId    string `json:"custom_id"`
	Name        string `json:"name"`
	TextContent string `json:"text_content"`
	Description string `json:"markdown_description"`
	Status      *struct {
		Status string `json:"status"`
		Type   string `json:"type"`
	} `json:"status"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	DateClosed  string `json:"date_closed"`
	Creator     *struct {
		Id       json.Number `json:"id"`
		Username string      `json:"username"`
	} `json:"creator"`
	Assignees []struct {
		Id       json.Number `json:"id"`
		Username string      `json:"username"`
	} `json:"assignees"`
	Priority *struct {
		Priority string `json:"priority"`
	} `json:"priority"`
	Parent string `json:"parent"`
	Url    string `json:"url"`
	// Points is ClickUp's native sprint points field (Fibonacci LOE for these
	// teams). It is the default story-point source.
	Points       *float64             `json:"points"`
	CustomFields []clickUpCustomField `json:"custom_fields"`
	TaskType     string               `json:"task_type"`
	List         *struct {
		Id string `json:"id"`
	} `json:"list"`
	Space *struct {
		Id string `json:"id"`
	} `json:"space"`
}

// clickUpCustomField is one entry of a task's custom_fields array. Value is left
// raw because ClickUp encodes it as a number or a string depending on the field.
type clickUpCustomField struct {
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
}

func ExtractTasks(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ClickUpTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickUpApiParams{
				ConnectionId: data.Options.ConnectionId,
				FolderId:     data.Options.FolderId,
			},
			Table: RAW_TASK_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			apiTask := &ClickUpApiTask{}
			if err := errors.Convert(json.Unmarshal(row.Data, apiTask)); err != nil {
				return nil, err
			}
			if apiTask.Id == "" {
				return nil, nil
			}
			description := apiTask.Description
			if description == "" {
				description = apiTask.TextContent
			}
			listId := ""
			if apiTask.List != nil {
				listId = apiTask.List.Id
			}
			task := &models.ClickUpTask{
				ConnectionId: data.Options.ConnectionId,
				Id:           apiTask.Id,
				ListId:       listId,
				FolderId:     data.Options.FolderId,
				CustomId:     apiTask.CustomId,
				Name:         apiTask.Name,
				Description:  description,
				Type:         apiTask.TaskType,
				ParentId:     parentOf(apiTask.Parent),
				StoryPoint:   storyPointOf(apiTask, data.ScopeConfig),
				Url:          apiTask.Url,
				CreatedDate:  parseClickUpTime(apiTask.DateCreated),
				UpdatedDate:  parseClickUpTime(apiTask.DateUpdated),
				ClosedDate:   parseClickUpTime(apiTask.DateClosed),
			}
			if apiTask.Space != nil {
				task.SpaceId = apiTask.Space.Id
			}
			if apiTask.Status != nil {
				task.Status = apiTask.Status.Status
				task.StatusType = apiTask.Status.Type
			}
			if apiTask.Priority != nil {
				task.Priority = apiTask.Priority.Priority
			}
			if apiTask.Creator != nil {
				task.CreatorId = apiTask.Creator.Id.String()
			}
			// TODO(clickup): ClickUp tasks can have multiple assignees. The MVP
			// keeps only the first for Issue.AssigneeId; emitting one
			// ticket.IssueAssignee row per assignee is a follow-up.
			if len(apiTask.Assignees) > 0 {
				task.AssigneeId = apiTask.Assignees[0].Id.String()
				task.AssigneeName = apiTask.Assignees[0].Username
			}
			return []interface{}{task}, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}

// parentOf normalizes ClickUp's parent field, which is the JSON literal null
// (decoded to an empty string) for top-level tasks.
func parentOf(parent string) string {
	if parent == "null" {
		return ""
	}
	return parent
}

// storyPointOf resolves a task's story points. Default source is ClickUp's
// native sprint points field. When the scope config names a custom field, that
// field's numeric value wins (teams that track Fibonacci LOE in a custom field).
func storyPointOf(apiTask *ClickUpApiTask, sc *models.ClickUpScopeConfig) *float64 {
	if sc != nil && sc.StoryPointField != "" {
		for _, cf := range apiTask.CustomFields {
			if strings.EqualFold(cf.Name, sc.StoryPointField) {
				return numericValue(cf.Value)
			}
		}
		return nil
	}
	return apiTask.Points
}

// numericValue coerces a raw custom-field value (JSON number or quoted string)
// into a float pointer; returns nil for empty / non-numeric values.
func numericValue(raw json.RawMessage) *float64 {
	s := strings.TrimSpace(strings.Trim(string(raw), `"`))
	if s == "" || s == "null" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &f
}
