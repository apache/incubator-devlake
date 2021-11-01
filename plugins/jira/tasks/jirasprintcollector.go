package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/utils"
	"net/http"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

type Sprint struct {
	ID            uint64    `json:"id"`
	Self          string    `json:"self"`
	State         string    `json:"state"`
	Name          string    `json:"name"`
	StartDate     time.Time `json:"startDate,omitempty"`
	EndDate       time.Time `json:"endDate,omitempty"`
	CompleteDate  time.Time `json:"completeDate,omitempty"`
	OriginBoardID int       `json:"originBoardId,omitempty"`
}
type JiraApiSprint struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	Total      int      `json:"total"`
	IsLast     bool     `json:"isLast"`
	Values     []Sprint `json:"values"`
}

func CollectSprint(ctx context.Context, jiraApiClient *JiraApiClient, source *models.JiraSource, boardId uint64) error {
	scheduler, err := utils.NewWorkerScheduler(10, 50, ctx)
	if err != nil {
		return err
	}
	defer scheduler.Release()
	err = jiraApiClient.FetchPages(scheduler, fmt.Sprintf("/agile/1.0/board/%v/sprint", boardId), nil, func(res *http.Response) error {
	jiraApiSprint := &JiraApiSprint{}
	err = core.UnmarshalResponse(res, jiraApiSprint)
	if err != nil {
		logger.Error("Error: ", err)
		return nil
	}
	logger.Info("jirasprint ", jiraApiSprint)
	for _, value := range jiraApiSprint.Values {
		jiraSprint := &models.JiraSprint{
			SourceId:      source.ID,
			SprintId:      value.ID,
			Self:          value.Self,
			State:         value.State,
			Name:          value.Name,
			StartDate:     value.StartDate,
			EndDate:       value.EndDate,
			CompleteDate:  value.CompleteDate,
			OriginBoardID: value.OriginBoardID,
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(jiraSprint).Error
		if err != nil {
			logger.Error("Error: ", err)
			return err
		}
		boardSprintRel := &models.JiraBoardSprint{
			SourceId: source.ID,
			BoardId:  boardId,
			SprintId: value.ID,
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(boardSprintRel).Error
		if err != nil {
			logger.Error("Error: ", err)
			return err
		}
		err = collectSprintIssueRelation(ctx, scheduler, jiraApiClient, source, boardId, value.ID)
		if err != nil {
			logger.Error("Error: ", err)
			return err
		}
	}
	return nil
})
	if err != nil {
		return err
	}
	scheduler.WaitUntilFinish()
	return nil
}