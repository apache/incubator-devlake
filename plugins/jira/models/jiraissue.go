package models

import (
	"database/sql"
	"time"

	"github.com/merico-dev/lake/models"
)

type JiraIssue struct {
	models.NoPKModel

	// collected fields
	SourceId       uint64 `gorm:"primaryKey"`
	IssueId        uint64 `gorm:"primarykey"`
	ProjectId      uint64
	Self           string
	Key            string
	Summary        string
	Type           string
	EpicKey        string
	StatusName     string
	StatusKey      string
	StatusCategory string
	StoryPoint     float64
	ResolutionDate sql.NullTime
	Created        time.Time
	Updated        time.Time

	// enriched fields
	// RequirementAnalsyisLeadTime uint
	// DesignLeadTime              uint
	// DevelopmentLeadTime         uint
	// TestLeadTime                uint
	// DeliveryLeadTime            uint
	LeadTime      uint
	StdStoryPoint uint
	StdType       string
	StdStatus     string

	// internal status tracking
	ChangelogUpdated sql.NullTime
}
