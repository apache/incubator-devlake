package models

import (
	"time"

	"github.com/merico-dev/lake/models"
)

// This Model is intended to save commits that are associated to a merge request
// for the purpose of calculating PR Review Rounds and other metrics that
// rely on associating commits to merge requests that may or may not
// exist on the main branch of a project.
// Thus a "Merge Request Commit" needs to be considered as distinct from a "Commit"

type GitlabMergeRequestCommit struct {
	CommitId       string `gorm:"primaryKey"`
	Title          string
	Message        string
	ShortId        string
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   time.Time
	CommitterName  string
	CommitterEmail string
	CommittedDate  time.Time
	WebUrl         string
	Additions      int
	Deletions      int
	Total          int
	models.NoPKModel
}
