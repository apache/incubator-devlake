package apiv2models

import (
	"github.com/merico-dev/lake/plugins/core"

	"github.com/merico-dev/lake/plugins/jira/models"
)

type Changelog struct {
	ID     uint64 `json:"id,string"`
	Author struct {
		Self         string `json:"self"`
		Name         string `json:"name"`
		Key          string `json:"key"`
		EmailAddress string `json:"emailAddress"`
		AvatarUrls   struct {
			Four8X48  string `json:"48x48"`
			Two4X24   string `json:"24x24"`
			One6X16   string `json:"16x16"`
			Three2X32 string `json:"32x32"`
		} `json:"avatarUrls"`
		DisplayName string `json:"displayName"`
		Active      bool   `json:"active"`
		TimeZone    string `json:"timeZone"`
	} `json:"author"`
	Created core.Iso8601Time `json:"created"`
	Items   []ChangelogItem  `json:"items"`
}

func (c Changelog) ToToolLayer(sourceId, issueId uint64) *models.JiraChangelog {
	return &models.JiraChangelog{
		SourceId:          sourceId,
		ChangelogId:       c.ID,
		IssueId:           issueId,
		AuthorAccountId:   c.Author.EmailAddress,
		AuthorDisplayName: c.Author.DisplayName,
		AuthorActive:      c.Author.Active,
		Created:           c.Created.ToTime(),
	}
}

type ChangelogItem struct {
	Field      string `json:"field"`
	Fieldtype  string `json:"fieldtype"`
	From       string `json:"from"`
	FromString string `json:"fromString"`
	To         string `json:"to"`
	ToString   string `json:"toString"`
}

func (c ChangelogItem) ToToolLayer(sourceId, changelogId uint64) *models.JiraChangelogItem {
	return &models.JiraChangelogItem{
		SourceId:    sourceId,
		ChangelogId: changelogId,
		Field:       c.Field,
		FieldType:   c.Fieldtype,
		From:        c.From,
		FromString:  c.FromString,
		To:          c.To,
		ToString:    c.ToString,
	}
}
