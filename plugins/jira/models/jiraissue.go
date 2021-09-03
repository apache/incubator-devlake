package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
)

var storyPointFieldId string

func init() {
	if config.V.GetString("JIRA_ISSUE_STORY_POINT_FIELD") != "" {
		storyPointFieldId = config.V.GetString("JIRA_ISSUE_STORY_POINT_FIELD")
	} else {
		storyPointFieldId = "customfield_10024"
	}
}

type JiraType struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type JiraIssueType struct {
	JiraType
	Subtask bool `json:"subtask"`
}

type JiraTypeWithKey struct {
	JiraType
	Key string `json:"key"`
}

func (u *JiraTypeWithKey) MarshalJSON() ([]byte, error) {
	type Alias JiraTypeWithKey
	t := &struct {
		Id json.Number `json:"id,omitempty"`
		*Alias
	}{json.Number(u.Id), (*Alias)(u)}
	return json.Marshal(t)
}

func (u *JiraTypeWithKey) UnmarshalJSON(data []byte) (err error) {
	type Alias JiraTypeWithKey
	t := &struct {
		Id json.Number `json:"id,omitempty"`
		*Alias
	}{json.Number(u.Id), (*Alias)(u)}
	err = json.Unmarshal(data, t)
	if err != nil {
		return err
	}
	*u = JiraTypeWithKey(*t.Alias)
	return nil
}

type JiraStatus struct {
	JiraTypeWithKey
	Category JiraTypeWithKey `json:"statusCategory" gorm:"embedded;embeddedPrefix:category_"`
}

type JiraIssueFields struct {
	Issuetype      JiraIssueType   `json:"issuetype,omitempty" gorm:"embedded;embeddedPrefix:type_"`
	Status         JiraStatus      `json:"status,omitempty" gorm:"embedded;embeddedPrefix:status_"`
	Summary        string          `json:"summary,omitempty" `
	Epic           JiraTypeWithKey `json:"epic,omitempty" gorm:"embedded;embeddedPrefix:epic_"`
	Project        JiraTypeWithKey `json:"project,omitempty" gorm:"embedded;embeddedPrefix:project_"`
	Created        time.Time       `json:"created,omitempty" `
	Updated        time.Time       `json:"updated,omitempty" `
	ResolutionDate sql.NullTime    `json:"resolutiondate,omitempty" `
	StoryPoint     uint64
}

func (u *JiraIssueFields) MarshalJSON() ([]byte, error) {
	type Alias JiraIssueFields
	fields := &struct {
		Created core.Iso8601Time `json:"created"`
		Updated core.Iso8601Time `json:"updated"`
		*Alias
	}{core.Iso8601Time(u.Created), core.Iso8601Time(u.Updated), (*Alias)(u)}
	return json.Marshal(fields)
}

func (u *JiraIssueFields) UnmarshalJSON(data []byte) (err error) {
	fieldsMapping := make(map[string]interface{})
	err = json.Unmarshal(data, &fieldsMapping)
	if err != nil {
		return err
	}
	type Alias JiraIssueFields
	fields := &struct {
		Created core.Iso8601Time `json:"created"`
		Updated core.Iso8601Time `json:"updated"`
		*Alias
	}{core.Iso8601Time(u.Created), core.Iso8601Time(u.Updated), (*Alias)(u)}
	err = json.Unmarshal(data, fields)
	if err != nil {
		return err
	}
	fields.Alias.Created = time.Time(fields.Created)
	fields.Alias.Updated = time.Time(fields.Updated)
	if fieldsMapping[storyPointFieldId] != nil {
		points := fieldsMapping[storyPointFieldId].(float64)
		fields.Alias.StoryPoint = uint64(points)
	}
	*u = JiraIssueFields(*fields.Alias)
	return nil
}

type JiraIssue struct {
	models.Model
	// collected field
	ID     string          `json:"id,omitempty" gorm:"primaryKey"` //overrider id to string type, make json parse same to model type
	Self   string          `json:"self,omitempty" `
	Key    string          `json:"key,omitempty"`
	Fields JiraIssueFields `json:"fields" gorm:"embedded"`

	// enriched fields
	Workload    float64
	LeadTime    uint
	StdWorkload uint
	StdType     string
	StdStatus   string
	// RequirementAnalsyisLeadTime uint
	// DesignLeadTime              uint
	// DevelopmentLeadTime         uint
	// TestLeadTime                uint
	// DeliveryLeadTime            uint
}
