package models

import (
	"time"

	"github.com/merico-dev/lake/models"
)

type JiraIssue struct {
	models.Model

	// collected field
	ProjectId  uint64
	Self       string
	Key        string
	Summary    string
	Type       string
	EpicKey    string
	StatusName string
	StatusKey  string
	Created    time.Time
	Updated    time.Time
}
