package apiv2models

import (
	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type Worklog struct {
	Self             string             `json:"self"`
	Author           *User              `json:"author"`
	UpdateAuthor     *User              `json:"updateAuthor"`
	Comment          string             `json:"comment"`
	Created          string             `json:"created"`
	Updated          helper.Iso8601Time `json:"updated"`
	Started          helper.Iso8601Time `json:"started"`
	TimeSpent        string             `json:"timeSpent"`
	TimeSpentSeconds int                `json:"timeSpentSeconds"`
	ID               string             `json:"id"`
	IssueID          uint64             `json:"issueId,string"`
}

func (w Worklog) ToToolLayer(connectionId uint64) *models.JiraWorklog {
	result := &models.JiraWorklog{
		ConnectionId:     connectionId,
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
