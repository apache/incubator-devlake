package models

import (
	"github.com/merico-dev/lake/models/common"
)

// This Model is intended to be an association table between merge request commits and merge requests.
// It needs to exist because there is a many to many relationship between merge request commits
// (which are commits associated to a merge request) and merge requests.

type GitlabMergeRequestCommitMergeRequest struct {
	MergeRequestCommitId string `gorm:"primaryKey"`
	MergeRequestId       int    `gorm:"primaryKey;autoIncrement:false"`
	common.NoPKModel
}
