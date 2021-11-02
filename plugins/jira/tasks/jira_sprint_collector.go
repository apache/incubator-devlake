package tasks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

type Sprint struct {
	ID            uint64     `json:"id"`
	Self          string     `json:"self"`
	State         string     `json:"state"`
	Name          string     `json:"name"`
	StartDate     *time.Time `json:"startDate,omitempty"`
	EndDate       *time.Time `json:"endDate,omitempty"`
	CompleteDate  *time.Time `json:"completeDate,omitempty"`
	OriginBoardID uint64     `json:"originBoardId,omitempty"`
}
type JiraApiSprint struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	Total      int      `json:"total"`
	IsLast     bool     `json:"isLast"`
	Values     []Sprint `json:"values"`
}

func CollectSprint(jiraApiClient *JiraApiClient, source *models.JiraSource, boardId uint64) error {
	err := jiraApiClient.FetchWithoutPaginationHeaders(fmt.Sprintf("/agile/1.0/board/%v/sprint", boardId), nil, func(res *http.Response) (bool, error) {
		jiraApiSprint := &JiraApiSprint{}
		err := core.UnmarshalResponse(res, jiraApiSprint)
		if err != nil {
			logger.Error("Error: ", err)
			return false, err
		}
		if len(jiraApiSprint.Values) == 0 {
			return false, nil
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
				return false, err
			}

			err = lakeModels.Db.FirstOrCreate(&models.JiraBoardSprint{
				SourceId: source.ID,
				BoardId:  boardId,
				SprintId: value.ID,
			}).Error
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
