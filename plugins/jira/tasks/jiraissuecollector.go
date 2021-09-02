package tasks

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

var epicKeyField, workloadField string

type JiraApiIssue struct {
	Id     string                 `json:"id"`
	Self   string                 `json:"self"`
	Key    string                 `json:"key"`
	Fields map[string]interface{} `json:"fields"`
}

type JiraApiIssuesResponse struct {
	Issues []JiraApiIssue `json:"issues"`
}

func init() {
	epicKeyField = config.V.GetString("JIRA_ISSUE_EPIC_KEY_FIELD")
	workloadField = config.V.GetString("JIRA_ISSUE_WORKLOAD_FIELD")
}

func CollectIssues(boardId uint64) error {
	jiraApiClient := GetJiraApiClient()
	return jiraApiClient.FetchWithPagination(fmt.Sprintf("/agile/1.0/board/%v/issue", boardId), nil,
		func(res *http.Response) error {
			// parse response
			jiraApiIssuesResponse := &JiraApiIssuesResponse{}
			err := core.UnmarshalResponse(res, jiraApiIssuesResponse)
			if err != nil {
				return err
			}

			// process issues
			for _, jiraApiIssue := range jiraApiIssuesResponse.Issues {

				jiraIssue, err := convertIssue(&jiraApiIssue)
				if err != nil {
					return err
				}
				// issue
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&jiraIssue).Error
				if err != nil {
					return err
				}

				// board / issue relationship
				lakeModels.Db.Create(&models.JiraBoardIssue{
					BoardId: boardId,
					IssueId: jiraIssue.ID,
				})
			}
			return nil
		})
}

func convertIssue(jiraApiIssue *JiraApiIssue) (*models.JiraIssue, error) {

	id, err := strconv.ParseUint(jiraApiIssue.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	created, err := time.Parse(core.ISO_8601_FORMAT, jiraApiIssue.Fields["created"].(string))
	if err != nil {
		return nil, err
	}
	updated, err := time.Parse(core.ISO_8601_FORMAT, jiraApiIssue.Fields["updated"].(string))
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
	if epicKeyField != "" {
		epicKey, _ = jiraApiIssue.Fields[epicKeyField].(string)
	}
	resolutionDate := sql.NullTime{}
	if rd, ok := jiraApiIssue.Fields["resolutiondate"]; ok && rd != nil {
		if resolutionDate.Time, err = time.Parse(core.ISO_8601_FORMAT, rd.(string)); err == nil {
			resolutionDate.Valid = true
		}
	}
	workload := 0.0
	if workloadField != "" {
		workload, _ = jiraApiIssue.Fields[workloadField].(float64)
	}
	jiraIssue := &models.JiraIssue{
		Model:          lakeModels.Model{ID: id},
		ProjectId:      projectId,
		Self:           jiraApiIssue.Self,
		Key:            jiraApiIssue.Key,
		Summary:        jiraApiIssue.Fields["summary"].(string),
		Type:           jiraApiIssue.Fields["issuetype"].(map[string]interface{})["name"].(string),
		StatusName:     statusName,
		StatusKey:      statusKey,
		EpicKey:        epicKey,
		ResolutionDate: resolutionDate,
		Workload:       workload,
		Created:        created,
		Updated:        updated,
	}
	return jiraIssue, nil
}
