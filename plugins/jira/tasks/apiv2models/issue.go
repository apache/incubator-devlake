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

package apiv2models

import (
	"encoding/json"

	"gorm.io/datatypes"

	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type Issue struct {
	Expand string `json:"expand"`
	ID     uint64 `json:"id,string"`
	Self   string `json:"self"`
	Key    string `json:"key"`
	Fields struct {
		AllFields map[string]interface{}
		Issuetype struct {
			Self        string `json:"self"`
			ID          string `json:"id"`
			Description string `json:"description"`
			IconURL     string `json:"iconUrl"`
			Name        string `json:"name"`
			Subtask     bool   `json:"subtask"`
			AvatarID    int    `json:"avatarId"`
		} `json:"issuetype"`
		Parent *struct {
			ID  uint64 `json:"id,string"`
			Key string `json:"key"`
		} `json:"parent"`
		Timespent     interface{} `json:"timespent"`
		Sprint        *Sprint     `json:"sprint"`
		ClosedSprints []Sprint    `json:"closedSprints"`
		Project       struct {
			Self           string `json:"self"`
			ID             uint64 `json:"id,string"`
			Key            string `json:"key"`
			Name           string `json:"name"`
			ProjectTypeKey string `json:"projectTypeKey"`
			AvatarUrls     struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
		} `json:"project"`
		FixVersions        []interface{}       `json:"fixVersions"`
		Aggregatetimespent interface{}         `json:"aggregatetimespent"`
		Resolution         interface{}         `json:"resolution"`
		Resolutiondate     *helper.Iso8601Time `json:"resolutiondate"`
		Workratio          int                 `json:"workratio"`
		LastViewed         string              `json:"lastViewed"`
		Watches            struct {
			Self       string `json:"self"`
			WatchCount int    `json:"watchCount"`
			IsWatching bool   `json:"isWatching"`
		} `json:"watches"`
		Created helper.Iso8601Time `json:"created"`
		Epic    *struct {
			ID      int    `json:"id"`
			Key     string `json:"key"`
			Self    string `json:"self"`
			Name    string `json:"name"`
			Summary string `json:"summary"`
			Color   struct {
				Key string `json:"key"`
			} `json:"color"`
			Done bool `json:"done"`
		} `json:"epic"`
		Priority *struct {
			Self    string `json:"self"`
			IconURL string `json:"iconUrl"`
			Name    string `json:"name"`
			ID      uint64 `json:"id,string"`
		} `json:"priority"`
		Labels                        []interface{}      `json:"labels"`
		Timeestimate                  interface{}        `json:"timeestimate"`
		Aggregatetimeoriginalestimate interface{}        `json:"aggregatetimeoriginalestimate"`
		Versions                      []interface{}      `json:"versions"`
		Issuelinks                    []interface{}      `json:"issuelinks"`
		Assignee                      *User              `json:"assignee"`
		Updated                       helper.Iso8601Time `json:"updated"`
		Status                        struct {
			Self           string `json:"self"`
			Description    string `json:"description"`
			IconURL        string `json:"iconUrl"`
			Name           string `json:"name"`
			ID             string `json:"id"`
			StatusCategory struct {
				Self      string `json:"self"`
				ID        int    `json:"id"`
				Key       string `json:"key"`
				ColorName string `json:"colorName"`
				Name      string `json:"name"`
			} `json:"statusCategory"`
		} `json:"status"`
		Timeoriginalestimate *int64      `json:"timeoriginalestimate"`
		Description          interface{} `json:"description"`
		Timetracking         *struct {
			RemainingEstimate        string `json:"remainingEstimate"`
			TimeSpent                string `json:"timeSpent"`
			RemainingEstimateSeconds int64  `json:"remainingEstimateSeconds"`
			TimeSpentSeconds         int    `json:"timeSpentSeconds"`
		} `json:"timetracking"`
		Archiveddate          interface{}   `json:"archiveddate"`
		Aggregatetimeestimate *int64        `json:"aggregatetimeestimate"`
		Summary               string        `json:"summary"`
		Creator               User          `json:"creator"`
		Subtasks              []interface{} `json:"subtasks"`
		Reporter              User          `json:"reporter"`
		Aggregateprogress     struct {
			Progress int `json:"progress"`
			Total    int `json:"total"`
		} `json:"aggregateprogress"`
		Environment interface{} `json:"environment"`
		Duedate     interface{} `json:"duedate"`
		Progress    struct {
			Progress int `json:"progress"`
			Total    int `json:"total"`
		} `json:"progress"`
		Worklog *struct {
			StartAt    int       `json:"startAt"`
			MaxResults int       `json:"maxResults"`
			Total      int       `json:"total"`
			Worklogs   []Worklog `json:"worklogs"`
		} `json:"worklog"`
	} `json:"fields"`
	Changelog *struct {
		StartAt    int         `json:"startAt"`
		MaxResults int         `json:"maxResults"`
		Total      int         `json:"total"`
		Histories  []Changelog `json:"histories"`
	} `json:"changelog"`
}

func (i Issue) toToolLayer(connectionId uint64) *models.JiraIssue {
	var workload float64
	result := &models.JiraIssue{
		ConnectionId:       connectionId,
		IssueId:            i.ID,
		ProjectId:          i.Fields.Project.ID,
		Self:               i.Self,
		IconURL:            i.Fields.Issuetype.IconURL,
		IssueKey:           i.Key,
		StoryPoint:         workload,
		Summary:            i.Fields.Summary,
		Type:               i.Fields.Issuetype.Name,
		StatusName:         i.Fields.Status.Name,
		StatusKey:          i.Fields.Status.StatusCategory.Key,
		ResolutionDate:     i.Fields.Resolutiondate.ToNullableTime(),
		CreatorAccountId:   i.Fields.Creator.getAccountId(),
		CreatorDisplayName: i.Fields.Creator.DisplayName,
		Created:            i.Fields.Created.ToTime(),
		Updated:            i.Fields.Updated.ToTime(),
	}
	if i.Fields.Epic != nil {
		result.EpicKey = i.Fields.Epic.Key
	}
	if i.Fields.Assignee != nil {
		result.AssigneeAccountId = i.Fields.Assignee.getAccountId()
		result.AssigneeDisplayName = i.Fields.Assignee.DisplayName
	}
	if i.Fields.Priority != nil {
		result.PriorityId = i.Fields.Priority.ID
		result.PriorityName = i.Fields.Priority.Name
	}
	if i.Fields.Timeoriginalestimate != nil {
		result.OriginalEstimateMinutes = *i.Fields.Timeoriginalestimate / 60
	}
	if i.Fields.Aggregatetimeestimate != nil {
		result.AggregateEstimateMinutes = *i.Fields.Aggregatetimeestimate / 60
	}
	if i.Fields.Timetracking != nil {
		result.RemainingEstimateMinutes = i.Fields.Timetracking.RemainingEstimateSeconds / 60
	}
	if i.Fields.Parent != nil {
		result.ParentId = i.Fields.Parent.ID
		result.ParentKey = i.Fields.Parent.Key
	}
	if i.Fields.Sprint != nil {
		result.SprintId = i.Fields.Sprint.ID
		result.SprintName = i.Fields.Sprint.Name
	}
	return result
}

func (i *Issue) SetAllFields(raw datatypes.JSON) error {
	var issue2 struct {
		Expand string          `json:"expand"`
		ID     uint64          `json:"id,string"`
		Self   string          `json:"self"`
		Key    string          `json:"key"`
		Fields json.RawMessage `json:"fields"`
	}
	err := json.Unmarshal(raw, &issue2)
	if err != nil {
		return err
	}
	err = json.Unmarshal(issue2.Fields, &i.Fields.AllFields)
	if err != nil {
		return err
	}
	return nil
}

func (i Issue) ExtractEntities(connectionId uint64) ([]uint64, *models.JiraIssue, bool, []*models.JiraWorklog, []*models.JiraChangelog, []*models.JiraChangelogItem, []*models.JiraUser) {
	issue := i.toToolLayer(connectionId)
	var worklogs []*models.JiraWorklog
	var changelogs []*models.JiraChangelog
	var changelogItems []*models.JiraChangelogItem
	var users []*models.JiraUser
	var needCollectWorklog bool
	var sprints []uint64
	if i.Fields.Worklog != nil {
		if i.Fields.Worklog.Total > len(i.Fields.Worklog.Worklogs) {
			needCollectWorklog = true
		} else {
			for _, w := range i.Fields.Worklog.Worklogs {
				worklogs = append(worklogs, w.ToToolLayer(connectionId))
			}
		}
	}
	if i.Changelog != nil {
		for _, changelog := range i.Changelog.Histories {
			cl, user := changelog.ToToolLayer(connectionId, i.ID)
			changelogs = append(changelogs, cl)
			users = append(users, user)
			for _, item := range changelog.Items {
				changelogItems = append(changelogItems, item.ToToolLayer(connectionId, changelog.ID))
				users = append(users, item.ExtractUser(connectionId)...)
			}
		}
	}
	if i.Fields.Sprint != nil {
		sprints = append(sprints, i.Fields.Sprint.ID)
	}
	for _, sprint := range i.Fields.ClosedSprints {
		sprints = append(sprints, sprint.ID)
	}
	users = append(users, i.Fields.Creator.ToToolLayer(connectionId), i.Fields.Reporter.ToToolLayer(connectionId))
	if i.Fields.Assignee != nil {
		users = append(users, i.Fields.Assignee.ToToolLayer(connectionId))
	}
	return sprints, issue, needCollectWorklog, worklogs, changelogs, changelogItems, users
}
