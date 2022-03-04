package v8models

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"

	"github.com/merico-dev/lake/plugins/jira/models"
)

type Worklog struct {
	Self             string           `json:"self"`
	Author           *User            `json:"author"`
	UpdateAuthor     *User            `json:"updateAuthor"`
	Comment          string           `json:"comment"`
	Created          string           `json:"created"`
	Updated          core.Iso8601Time `json:"updated"`
	Started          core.Iso8601Time `json:"started"`
	TimeSpent        string           `json:"timeSpent"`
	TimeSpentSeconds int              `json:"timeSpentSeconds"`
	ID               string           `json:"id"`
	IssueID          string           `json:"issueId"`
}

func (w Worklog) toToolLayer(sourceId, issueId uint64) *models.JiraWorklog {
	result := &models.JiraWorklog{
		SourceId:         sourceId,
		IssueId:          issueId,
		WorklogId:        w.ID,
		TimeSpent:        w.TimeSpent,
		TimeSpentSeconds: w.TimeSpentSeconds,
		Updated:          w.Updated.ToTime(),
		Started:          w.Started.ToTime(),
	}
	if w.Author != nil {
		result.AuthorId = w.Author.EmailAddress
	}
	if w.UpdateAuthor != nil {
		result.UpdateAuthorId = w.UpdateAuthor.EmailAddress
	}
	return result
}

func (Worklog) FromAPI(sourceId, issueId uint64, raw json.RawMessage) (interface{}, error) {
	var vv []Worklog
	err := json.Unmarshal(raw, &vv)
	if err != nil {
		return nil, err
	}
	list := make([]*models.JiraWorklog, len(vv))
	for i, item := range vv {
		list[i] = item.toToolLayer(sourceId, issueId)
	}
	return list, nil
}

func (Worklog) ExtractRawMessage(blob []byte) (json.RawMessage, error) {
	var resp struct {
		StartAt    int             `json:"startAt"`
		MaxResults int             `json:"maxResults"`
		Total      int             `json:"total"`
		Worklogs   json.RawMessage `json:"worklogs"`
	}
	err := json.Unmarshal(blob, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Worklogs, nil
}
