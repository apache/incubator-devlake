package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdBug struct {
	SourceId    models.Uint64s `gorm:"primaryKey"`
	ID          models.Uint64s `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	EpicKey     string
	Title       string `json:"name"`
	Description string
	WorkspaceId models.Uint64s    `json:"workspace_id"`
	Created     *core.Iso8601Time `json:"created"`
	Modified    *core.Iso8601Time `json:"modified" gorm:"index"`
	Status      string            `json:"status"`
	Cc          string            `json:"cc"`
	Begin       *core.Iso8601Time `json:"begin"`
	Due         *core.Iso8601Time `json:"due"`
	Priority    string            `json:"priority"`
	IterationID models.Uint64s    `json:"iteration_id"`
	Source      string            `json:"source"`
	Module      string            `json:"module"`
	ReleaseID   models.Uint64s    `json:"release_id"`
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
	IssueID          models.Uint64s    `json:"issue_id"`
	VerifyTime       *core.Iso8601Time `json:"verify_time"`
	RejectTime       *core.Iso8601Time `json:"reject_time"`
	ReopenTime       *core.Iso8601Time `json:"reopen_time"`
	AuditTime        *core.Iso8601Time `json:"audit_time"`
	SuspendTime      *core.Iso8601Time `json:"suspend_time"`
	Deadline         *core.Iso8601Time `json:"deadline"`
	InProgressTime   *core.Iso8601Time `json:"in_progress_time"`
	AssignedTime     *core.Iso8601Time `json:"assigned_time"`
	TemplateID       models.Uint64s    `json:"template_id"`
	StoryID          models.Uint64s    `json:"story_id"`
	StdStatus        string
	StdType          string
	Type             string
	Url              string

	SupportID       models.Uint64s `json:"support_id"`
	SupportForumID  models.Uint64s `json:"support_forum_id"`
	TicketID        models.Uint64s `json:"ticket_id"`
	Follower        string         `json:"follower"`
	SyncType        string         `json:"sync_type"`
	Label           string         `json:"label"`
	Effort          models.Ints    `json:"effort"`
	EffortCompleted models.Ints    `json:"effort_completed"`
	Exceed          models.Ints    `json:"exceed"`
	Remain          models.Ints    `json:"remain"`
	Progress        string         `json:"progress"`
}

func (TapdBug) TableName() string {
	return "_tool_tapd_bugs"
}
