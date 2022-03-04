package v8models

import (
	"encoding/json"
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

func (s Sprint) toToolLayer(sourceId uint64) *models.JiraSprint {
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
func (s Sprint) FromAPI(sourceId uint64, raw json.RawMessage) (interface{}, error) {
	v, err := s.GetJiraSprints(sourceId, raw)
	return v, err
}
func (s Sprint) GetJiraSprints(sourceId uint64, raw json.RawMessage) ([]*models.JiraSprint, error) {
	var vv []Sprint
	err := json.Unmarshal(raw, &vv)
	if err != nil {
		return nil, err
	}
	list := make([]*models.JiraSprint, len(vv))
	for i, item := range vv {
		list[i] = item.toToolLayer(sourceId)
	}
	return list, nil
}

func (Sprint) ExtractRawMessage(blob []byte) (json.RawMessage, error) {
	var resp struct {
		MaxResults int             `json:"maxResults"`
		StartAt    int             `json:"startAt"`
		IsLast     bool            `json:"isLast"`
		Values     json.RawMessage `json:"values"`
	}
	err := json.Unmarshal(blob, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Values, nil
}
