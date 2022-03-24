package apiv2models

import (
	"time"

	"github.com/merico-dev/lake/plugins/jira/models"
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

func (s Sprint) ToToolLayer(sourceId uint64) *models.JiraSprint {
	return &models.JiraSprint{
		SourceId:      sourceId,
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
