package tasks

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
)

const TIME_FORMAT = "2006-01-02T15:04:05-0700"

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

type JiraApiIssuesResponse struct {
	JiraPagination
	Issues []JiraApiIssue `json:"issues"`
}

func init() {
	epicFieldName = config.V.GetString("JIRA_ISSUE_EPIC_KEY_FIELD")
}

func CollectIssues(boardId uint64) error {
	jiraApiClient := GetJiraApiClient()
	return jiraApiClient.FetchPages(fmt.Sprintf("/agile/1.0/board/%v/issue", boardId), nil,
		func(res *http.Response) (*JiraPagination, error) {
			// parse response
			jiraApiIssuesResponse := &JiraApiIssuesResponse{}
			err := core.UnmarshalResponse(res, jiraApiIssuesResponse)
			if err != nil {
				logger.Error("Error: ", err)
				return nil, err
			}

			// process issues
			for _, jiraApiIssue := range jiraApiIssuesResponse.Issues {

				jiraIssue, err := convertIssue(&jiraApiIssue)
				if err != nil {
					logger.Error("Error: ", err)
					break
				}
				// issue
				err = lakeModels.Db.Save(jiraIssue).Error
				if err != nil {
					logger.Error("Error: ", err)
					break
				}

				// board / issue relationship
				jiraBoardIssue := &models.JiraBoardIssue{
					BoardId: boardId,
					IssueId: jiraIssue.ID,
				}
				err = lakeModels.Db.Save(jiraBoardIssue).Error
				if err != nil {
					logger.Error("Error: ", err)
					break
				}
			}

			// return pagination infomration
			return &jiraApiIssuesResponse.JiraPagination, nil
		})
}

func convertIssue(jiraApiIssue *JiraApiIssue) (*models.JiraIssue, error) {

	id, err := strconv.ParseUint(jiraApiIssue.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	created, err := time.Parse(TIME_FORMAT, jiraApiIssue.Fields["created"].(string))
	if err != nil {
		return nil, err
	}
	updated, err := time.Parse(TIME_FORMAT, jiraApiIssue.Fields["updated"].(string))
	if err != nil {
		return nil, err
	}
	projectId, err := strconv.ParseUint(
		jiraApiIssue.Fields["project"].(map[string]interface{})["id"].(string), 10, 64,
	)
	if err != nil {
		return nil, err
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
	return jiraIssue, nil
}
