package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdBug struct {
	SourceId    uint64 `gorm:"primaryKey"`
	ID          uint64 `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	EpicKey     string `gorm:"type:varchar(255)"`
	Title       string `gorm:"type:varchar(255)"`
	Description string
	WorkspaceId uint64     `json:"workspace_id"`
	Created     *time.Time `json:"created"`
	Modified    *time.Time `json:"modified" gorm:"index"`
	Status      string     `gorm:"type:varchar(255)"`
	Cc          string     `gorm:"type:varchar(255)"`
	Begin       *time.Time `json:"begin"`
	Due         *time.Time `json:"due"`
	Priority    string     `gorm:"type:varchar(255)"`
	IterationID uint64     `json:"iteration_id"`
	Source      string     `gorm:"type:varchar(255)"`
	Module      string     `gorm:"type:varchar(255)"`
	ReleaseID   uint64     `json:"release_id"`
	CreatedFrom string     `gorm:"type:varchar(255)"`
	Feature     string     `gorm:"type:varchar(255)"`
	common.NoPKModel

	Severity         string     `gorm:"type:varchar(255)"`
	Reporter         string     `gorm:"type:varchar(255)"`
	Resolved         *time.Time `json:"resolved"`
	Closed           *time.Time `json:"closed"`
	Lastmodify       string     `gorm:"type:varchar(255)"`
	Auditer          string     `gorm:"type:varchar(255)"`
	De               string     `gorm:"type:varchar(255);comment:developer"`
	Fixer            string     `gorm:"type:varchar(255)"`
	VersionTest      string     `gorm:"type:varchar(255)"`
	VersionReport    string     `gorm:"type:varchar(255)"`
	VersionClose     string     `gorm:"type:varchar(255)"`
	VersionFix       string     `gorm:"type:varchar(255)"`
	BaselineFind     string     `gorm:"type:varchar(255)"`
	BaselineJoin     string     `gorm:"type:varchar(255)"`
	BaselineClose    string     `gorm:"type:varchar(255)"`
	BaselineTest     string     `gorm:"type:varchar(255)"`
	Sourcephase      string     `gorm:"type:varchar(255)"`
	Te               string     `gorm:"type:varchar(255)"`
	CurrentOwner     string     `gorm:"type:varchar(255)"`
	Resolution       string     `gorm:"type:varchar(255)"`
	Originphase      string     `gorm:"type:varchar(255)"`
	Confirmer        string     `gorm:"type:varchar(255)"`
	Participator     string     `gorm:"type:varchar(255)"`
	Closer           string     `gorm:"type:varchar(255)"`
	Platform         string     `gorm:"type:varchar(255)"`
	Os               string     `gorm:"type:varchar(255)"`
	Testtype         string     `gorm:"type:varchar(255)"`
	Testphase        string     `gorm:"type:varchar(255)"`
	Frequency        string     `gorm:"type:varchar(255)"`
	RegressionNumber string     `gorm:"type:varchar(255)"`
	Flows            string     `gorm:"type:varchar(255)"`
	Testmode         string     `gorm:"type:varchar(255)"`
	IssueID          uint64     `json:"issue_id"`
	VerifyTime       *time.Time `json:"verify_time"`
	RejectTime       *time.Time `json:"reject_time"`
	ReopenTime       *time.Time `json:"reopen_time"`
	AuditTime        *time.Time `json:"audit_time"`
	SuspendTime      *time.Time `json:"suspend_time"`
	Deadline         *time.Time `json:"deadline"`
	InProgressTime   *time.Time `json:"in_progress_time"`
	AssignedTime     *time.Time `json:"assigned_time"`
	TemplateID       uint64     `json:"template_id"`
	StoryID          uint64     `json:"story_id"`
	StdStatus        string     `gorm:"type:varchar(255)"`
	StdType          string     `gorm:"type:varchar(255)"`
	Type             string     `gorm:"type:varchar(255)"`
	Url              string     `gorm:"type:varchar(255)"`
}

func (TapdBug) TableName() string {
	return "_tool_tapd_bugs"
}
