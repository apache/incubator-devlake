package models

import (
	"github.com/merico-dev/lake/models"
)

type JiraIssue struct {
	models.Model
	// JiraId          string
	Key string
	// ProjectId       int
	// Url             string
	// Title string
	// Description     string
	// LeadTime        int
	// IssueType       string
	// EpicKey         string
	// Status          string
	// IssueCreatedAt  string
	// IssueUpdatedAt  string
	// IssueResolvedAt string
}
