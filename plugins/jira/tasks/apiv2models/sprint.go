package apiv2models

import (
	"time"

	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type Sprint struct {
	ID            uint64     `json:"id"`
	Self          string     `json:"self"`
	State         string     `json:"state"`
	Name          string     `json:"name"`
	StartDate     *time.Time `json:"startDate"`
	EndDate       *time.Time `json:"endDate"`
	CompleteDate  *time.Time `json:"completeDate"`
	OriginBoardID uint64     `json:"originBoardId"`
}

func (s Sprint) ToToolLayer(connectionId uint64) *models.JiraSprint {
	return &models.JiraSprint{
		ConnectionId:  connectionId,
		SprintId:      s.ID,
		Self:          s.Self,
		State:         s.State,
		Name:          s.Name,
		StartDate:     s.StartDate,
		EndDate:       s.EndDate,
		CompleteDate:  s.CompleteDate,
		OriginBoardID: s.OriginBoardID,
	}
}
