package apiv2models

import (
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
	IssueID          uint64           `json:"issueId,string"`
}

func (w Worklog) ToToolLayer(sourceId uint64) *models.JiraWorklog {
	result := &models.JiraWorklog{
		SourceId:         sourceId,
		IssueId:          w.IssueID,
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
