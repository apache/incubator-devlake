package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdBug struct {
	SourceId    uint64 `gorm:"primaryKey"`
	ID          uint64 `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	EpicKey     string
	Title       string     `json:"name"`
	Description string     `json:"description"`
	WorkspaceID uint64     `json:"workspace_id"`
	Created     *time.Time `json:"created"`
	Modified    *time.Time `json:"modified" gorm:"index"`
	Status      string     `json:"status"`
	Cc          string     `json:"cc"`
	Begin       *time.Time `json:"begin"`
	Due         *time.Time `json:"due"`
	Priority    string     `json:"priority"`
	IterationID uint64     `json:"iteration_id"`
	Source      string     `json:"source"`
	Module      string     `json:"module"`
	ReleaseID   uint64     `json:"release_id"`
	CreatedFrom string     `json:"created_from"`
	Feature     string     `json:"feature"`
	common.NoPKModel

	Severity         string     `json:"severity"`
	Reporter         string     `json:"reporter"`
	Resolved         *time.Time `json:"resolved"`
	Closed           *time.Time `json:"closed"`
	Lastmodify       string     `json:"lastmodify"`
	Auditer          string     `json:"auditer"`
	De               string     `json:"De" gorm:"comment:developer"`
	Fixer            string     `json:"fixer"`
	VersionTest      string     `json:"version_test"`
	VersionReport    string     `json:"version_report"`
	VersionClose     string     `json:"version_close"`
	VersionFix       string     `json:"version_fix"`
	BaselineFind     string     `json:"baseline_find"`
	BaselineJoin     string     `json:"baseline_join"`
	BaselineClose    string     `json:"baseline_close"`
	BaselineTest     string     `json:"baseline_test"`
	Sourcephase      string     `json:"sourcephase"`
	Te               string     `json:"te"`
	CurrentOwner     string     `json:"current_owner"`
	Resolution       string     `json:"resolution"`
	Originphase      string     `json:"originphase"`
	Confirmer        string     `json:"confirmer"`
	Participator     string     `json:"participator"`
	Closer           string     `json:"closer"`
	Platform         string     `json:"platform"`
	Os               string     `json:"os"`
	Testtype         string     `json:"testtype"`
	Testphase        string     `json:"testphase"`
	Frequency        string     `json:"frequency"`
	RegressionNumber string     `json:"regression_number"`
	Flows            string     `json:"flows"`
	Testmode         string     `json:"testmode"`
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
}

type TapdBugApiRes struct {
	ID          string `gorm:"primaryKey" json:"id"`
	EpicKey     string
	Title       string `json:"name"`
	Description string `json:"description"`
	WorkspaceID string `json:"workspace_id"`
	Created     string `json:"created"`
	Modified    string `json:"modified" gorm:"index"`
	Status      string `json:"status"`
	Cc          string `json:"cc"`
	Begin       string `json:"begin"`
	Due         string `json:"due"`
	Priority    string `json:"priority"`
	IterationID string `json:"iteration_id"`
	Source      string `json:"source"`
	Module      string `json:"module"`
	ReleaseID   string `json:"release_id"`
	CreatedFrom string `json:"created_from"`
	Feature     string `json:"feature"`

	Severity         string `json:"severity"`
	Reporter         string `json:"reporter"`
	Bugtype          string `json:"bugtype"`
	Resolved         string `json:"resolved"`
	Closed           string `json:"closed"`
	Lastmodify       string `json:"lastmodify"`
	Auditer          string `json:"auditer"`
	De               string `json:"de"`
	Fixer            string `json:"fixer"`
	VersionTest      string `json:"version_test"`
	VersionReport    string `json:"version_report"`
	VersionClose     string `json:"version_close"`
	VersionFix       string `json:"version_fix"`
	BaselineFind     string `json:"baseline_find"`
	BaselineJoin     string `json:"baseline_join"`
	BaselineClose    string `json:"baseline_close"`
	BaselineTest     string `json:"baseline_test"`
	Sourcephase      string `json:"sourcephase"`
	Te               string `json:"te" gorm:"comment:developer"`
	CurrentOwner     string `json:"current_owner"`
	Resolution       string `json:"resolution"`
	Originphase      string `json:"originphase"`
	Confirmer        string `json:"confirmer"`
	Participator     string `json:"participator"`
	Closer           string `json:"closer"`
	Platform         string `json:"platform"`
	Os               string `json:"os"`
	Testtype         string `json:"testtype"`
	Testphase        string `json:"testphase"`
	Frequency        string `json:"frequency"`
	RegressionNumber string `json:"regression_number"`
	Flows            string `json:"flows"`
	Testmode         string `json:"testmode"`
	IssueID          string `json:"issue_id"`
	VerifyTime       string `json:"verify_time"`
	RejectTime       string `json:"reject_time"`
	ReopenTime       string `json:"reopen_time"`
	AuditTime        string `json:"audit_time"`
	SuspendTime      string `json:"suspend_time"`
	Deadline         string `json:"deadline"`
	InProgressTime   string `json:"in_progress_time"`
	AssignedTime     string `json:"assigned_time"`
	TemplateID       string `json:"template_id"`
	StoryID          string `json:"story_id"`
}
