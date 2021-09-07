package tasks

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/utils"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

type JiraApiIssuesResponse struct {
	JiraPagination
	Issues []models.JiraIssue `json:"issues"`
}

func CollectIssues(boardId uint64) error {
	jiraApiClient := GetJiraApiClient()
	// diff sync
	lastestUpdated := &models.JiraIssue{}
	err := lakeModels.Db.Order("updated DESC").Select("id", "updated").Limit(1).Find(lastestUpdated).Error
	if err != nil {
		return err
	}
	jql := "ORDER BY updated ASC"
	if lastestUpdated != nil {
		jql = fmt.Sprintf("update >= %v %v", lastestUpdated.Fields.Updated.Format("2006/01/02 15:04"), jql)
	}
	query := &url.Values{}
	query.Set("jql", jql)

	scheduler, err := utils.NewWorkerScheduler(10, 50)
	if err != nil {
		return err
	}
	defer scheduler.Release()

	err = jiraApiClient.FetchPages(scheduler, fmt.Sprintf("/agile/1.0/board/%v/issue", boardId), query,
		func(res *http.Response) error {
			// parse response
			jiraApiIssuesResponse := &JiraApiIssuesResponse{}
			err := core.UnmarshalResponse(res, jiraApiIssuesResponse)
			if err != nil {
				return err
			}

			// process issues
			err = lakeModels.Db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&jiraApiIssuesResponse.Issues).Error
			if err != nil {
				return err
			}

			return nil
		})
	if err != nil {
		return err
	}
	scheduler.WaitUntilFinish()
	return nil
}
