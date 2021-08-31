package tasks

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
)

var epicFieldName string

type JiraPagination struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}

type JiraApiIssue struct {
	Id     string                 `json:"id"`
	Self   string                 `json:"self"`
	Key    string                 `json:"key"`
	Fields map[string]interface{} `json:"fields"`
}

type JiraApiResponse struct {
	JiraPagination
	Issues []JiraApiIssue `json:"issues"`
}

func init() {
	epicFieldName = config.V.GetString("JIRA_ISSUE_EPIC_KEY_FIELD")
}

func CollectIssues(boardId uint64) error {
	jiraApiClient := GetJiraApiClient()

	loaded, total, query := 0, 1, &url.Values{}
	for loaded < total {
		// fetch page
		query.Set("maxResults", "100")
		query.Set("startAt", strconv.Itoa(loaded))
		res, err := jiraApiClient.Get(fmt.Sprintf("/agile/1.0/board/%v/issue", boardId), query, nil)
		if err != nil {
			return err
		}

		// parse response
		jiraApiResponse := &JiraApiResponse{}
		err = core.UnmarshalResponse(res, jiraApiResponse)
		if err != nil {
			logger.Error("Error: ", err)
			return nil
		}

		// save issues
		SaveIssues(boardId, jiraApiResponse.Issues)

		// next page
		loaded += len(jiraApiResponse.Issues)
		total = jiraApiResponse.Total
		logger.Info("jira board issues collection", map[string]interface{}{
			"boardId": boardId,
			"loaded":  loaded,
			"total":   total,
		})
	}
	return nil
}

func SaveIssues(boardId uint64, issues []JiraApiIssue) {
	const TIME_FORMAT = "2006-01-02T15:04:05-0700"

	// convert and save
	for _, jiraApiIssue := range issues {

		// issue
		id, err := strconv.ParseUint(jiraApiIssue.Id, 10, 64)
		if err != nil {
			logger.Error("Error: ", err)
			break
		}
		created, err := time.Parse(TIME_FORMAT, jiraApiIssue.Fields["created"].(string))
		if err != nil {
			logger.Error("Error: ", err)
			break
		}
		updated, err := time.Parse(TIME_FORMAT, jiraApiIssue.Fields["updated"].(string))
		if err != nil {
			logger.Error("Error: ", err)
			break
		}
		projectId, err := strconv.ParseUint(
			jiraApiIssue.Fields["project"].(map[string]interface{})["id"].(string), 10, 64,
		)
		if err != nil {
			logger.Error("Error: ", err)
			break
		}
		status := jiraApiIssue.Fields["status"].(map[string]interface{})
		statusName := status["name"].(string)
		statusKey := status["statusCategory"].(map[string]interface{})["key"].(string)
		epicKey := ""
		if epicFieldName != "" {
			epicKey, _ = jiraApiIssue.Fields[epicFieldName].(string)
		}
		jiraIssue := &models.JiraIssue{
			Model:      lakeModels.Model{ID: id},
			ProjectId:  projectId,
			Self:       jiraApiIssue.Self,
			Key:        jiraApiIssue.Key,
			Summary:    jiraApiIssue.Fields["summary"].(string),
			Type:       jiraApiIssue.Fields["issuetype"].(map[string]interface{})["name"].(string),
			StatusName: statusName,
			StatusKey:  statusKey,
			EpicKey:    epicKey,
			Created:    created,
			Updated:    updated,
		}
		err = lakeModels.Db.Save(jiraIssue).Error
		if err != nil {
			logger.Error("Error: ", err)
			break
		}

		// board / issue relationship
		jiraBoardIssue := &models.JiraBoardIssue{
			BoardId: boardId,
			IssueId: id,
		}
		err = lakeModels.Db.Save(jiraBoardIssue).Error
		if err != nil {
			logger.Error("Error: ", err)
			break
		}
	}
}
