package apiv2models

import (
	"github.com/merico-dev/lake/plugins/core"

	"github.com/merico-dev/lake/plugins/jira/models"
)

type Changelog struct {
	ID      uint64           `json:"id,string"`
	Author  User             `json:"author"`
	Created core.Iso8601Time `json:"created"`
	Items   []ChangelogItem  `json:"items"`
}

func (c Changelog) ToToolLayer(connectionId, issueId uint64) (*models.JiraChangelog, *models.JiraUser) {
	return &models.JiraChangelog{
		ConnectionId:      connectionId,
		ChangelogId:       c.ID,
		IssueId:           issueId,
		AuthorAccountId:   c.Author.getAccountId(),
		AuthorDisplayName: c.Author.DisplayName,
		AuthorActive:      c.Author.Active,
		Created:           c.Created.ToTime(),
	}, c.Author.ToToolLayer(connectionId)
}

type ChangelogItem struct {
	Field      string `json:"field"`
	Fieldtype  string `json:"fieldtype"`
	From       string `json:"from"`
	FromString string `json:"fromString"`
	To         string `json:"to"`
	ToString   string `json:"toString"`
}

func (c ChangelogItem) ToToolLayer(connectionId, changelogId uint64) *models.JiraChangelogItem {
	return &models.JiraChangelogItem{
		ConnectionId: connectionId,
		ChangelogId:  changelogId,
		Field:        c.Field,
		FieldType:    c.Fieldtype,
		From:         c.From,
		FromString:   c.FromString,
		To:           c.To,
		ToString:     c.ToString,
	}
}

func (c ChangelogItem) ExtractUser(connectionId uint64) []*models.JiraUser {
	if c.Field != "assignee" {
		return nil
	}
	var result []*models.JiraUser
	if c.From != "" {
		result = append(result, &models.JiraUser{ConnectionId: connectionId, AccountId: c.From})
	}
	if c.To != "" {
		result = append(result, &models.JiraUser{ConnectionId: connectionId, AccountId: c.To})
	}
	return result
}
