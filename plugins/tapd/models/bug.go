package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdBug struct {
	SourceId    Uint64s `gorm:"primaryKey"`
	ID          Uint64s `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	EpicKey     string
	Title       string `json:"name"`
	Description string
	WorkspaceID Uint64s           `json:"workspace_id"`
	Created     *core.Iso8601Time `json:"created"`
	Modified    *core.Iso8601Time `json:"modified" gorm:"index"`
	Status      string            `json:"status"`
	Cc          string            `json:"cc"`
	Begin       *core.Iso8601Time `json:"begin"`
	Due         *core.Iso8601Time `json:"due"`
	Priority    string            `json:"priority"`
	IterationID Uint64s           `json:"iteration_id"`
	Source      string            `json:"source"`
	Module      string            `json:"module"`
	ReleaseID   Uint64s           `json:"release_id"`
	CreatedFrom string            `json:"created_from"`
	Feature     string            `json:"feature"`
	common.NoPKModel

	Severity         string            `json:"severity"`
	Reporter         string            `json:"reporter"`
	Resolved         *core.Iso8601Time `json:"resolved"`
	Closed           *core.Iso8601Time `json:"closed"`
	Lastmodify       string            `json:"lastmodify"`
	Auditer          string            `json:"auditer"`
	De               string            `json:"De" gorm:"comment:developer"`
	Fixer            string            `json:"fixer"`
	VersionTest      string            `json:"version_test"`
	VersionReport    string            `json:"version_report"`
	VersionClose     string            `json:"version_close"`
	VersionFix       string            `json:"version_fix"`
	BaselineFind     string            `json:"baseline_find"`
	BaselineJoin     string            `json:"baseline_join"`
	BaselineClose    string            `json:"baseline_close"`
	BaselineTest     string            `json:"baseline_test"`
	Sourcephase      string            `json:"sourcephase"`
	Te               string            `json:"te"`
	CurrentOwner     string            `json:"current_owner"`
	Resolution       string            `json:"resolution"`
	Originphase      string            `json:"originphase"`
	Confirmer        string            `json:"confirmer"`
	Participator     string            `json:"participator"`
	Closer           string            `json:"closer"`
	Platform         string            `json:"platform"`
	Os               string            `json:"os"`
	Testtype         string            `json:"testtype"`
	Testphase        string            `json:"testphase"`
	Frequency        string            `json:"frequency"`
	RegressionNumber string            `json:"regression_number"`
	Flows            string            `json:"flows"`
	Testmode         string            `json:"testmode"`
	IssueID          Uint64s           `json:"issue_id"`
	VerifyTime       *core.Iso8601Time `json:"verify_time"`
	RejectTime       *core.Iso8601Time `json:"reject_time"`
	ReopenTime       *core.Iso8601Time `json:"reopen_time"`
	AuditTime        *core.Iso8601Time `json:"audit_time"`
	SuspendTime      *core.Iso8601Time `json:"suspend_time"`
	Deadline         *core.Iso8601Time `json:"deadline"`
	InProgressTime   *core.Iso8601Time `json:"in_progress_time"`
	AssignedTime     *core.Iso8601Time `json:"assigned_time"`
	TemplateID       Uint64s           `json:"template_id"`
	StoryID          Uint64s           `json:"story_id"`
	StdStatus        string
	StdType          string
	Type             string
	Url              string

	SupportID       Uint64s `json:"support_id"`
	SupportForumID  Uint64s `json:"support_forum_id"`
	TicketID        Uint64s `json:"ticket_id"`
	Follower        string  `json:"follower"`
	SyncType        string  `json:"sync_type"`
	Label           string  `json:"label"`
	Effort          Floats  `json:"effort"`
	EffortCompleted Floats  `json:"effort_completed"`
	Exceed          Floats  `json:"exceed"`
	Remain          Floats  `json:"remain"`
	Progress        string  `json:"progress"`
	Estimate        Floats  `json:"estimate"`
	Bugtype         string  `json:"bugtype"`

	Milestone string `json:"milestone"`
}

func (TapdBug) TableName() string {
	return "_tool_tapd_bugs"
}
