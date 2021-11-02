package tasks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

type JiraApiSprintIssueRelation struct {
	Expand     string `json:"expand"`
	StartAt    int    `json:"startAt"`
	MaxResults int    `json:"maxResults"`
	Total      int    `json:"total"`
	Issues     []struct {
		ID uint64 `json:"id,string"`
	} `json:"issues"`
}

func collectSprintIssueRelation(ctx context.Context, jiraApiClient *JiraApiClient, source *models.JiraSource, boardId, sprintId uint64) error {
	err := jiraApiClient.FetchWithoutPaginationHeaders(fmt.Sprintf("/agile/1.0/board/%v/sprint/%d/issue", boardId, sprintId), nil, func(res *http.Response) (bool, error) {
		rel := &JiraApiSprintIssueRelation{}
		err := core.UnmarshalResponse(res, rel)
		if err != nil {
			logger.Error("Error: ", err)
			return false, err
		}
		logger.Info("jira sprint issue relation ", rel)
		for _, value := range rel.Issues {
			sprintIssueRel := &models.JiraSprintIssue{
				SourceId: source.ID,
				SprintId: sprintId,
				IssueId:  value.ID,
			}
			err = lakeModels.Db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(sprintIssueRel).Error
			if err != nil {
				logger.Error("Error: ", err)
				return false, err
			}
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	return nil
}
