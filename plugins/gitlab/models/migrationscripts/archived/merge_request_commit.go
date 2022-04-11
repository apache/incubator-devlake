package archived

import (
	"github.com/merico-dev/lake/models/common"
)

// This Model is intended to save commits that are associated to a merge request
// for the purpose of calculating PR Review Rounds and other metrics that
// rely on associating commits to merge requests that may or may not
// exist on the main branch of a project.
// Thus a "Merge Request Commit" needs to be considered as distinct from a "Commit"

type GitlabMergeRequestCommit struct {
	CommitSha      string `gorm:"primaryKey;type:varchar(40)"`
	MergeRequestId int    `gorm:"primaryKey;autoIncrement:false"`
	common.NoPKModel
}

func (GitlabMergeRequestCommit) TableName() string {
	return "_tool_gitlab_merge_request_commits"
}
