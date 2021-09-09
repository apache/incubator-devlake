package tasks

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
)

func GetReviewRounds(commits []gitlabModels.GitlabMergeRequestCommit, notes []gitlabModels.GitlabMergeRequestNote) int {
	i := 0
	j := 0
	reviewRounds := 0
	state := 0 // 0, 1, 2
	for i < len(commits) && j < len(notes) {
		if commits[i].AuthoredDate < notes[j].GitlabCreatedAt {
			i++
			state = 1
		} else {
			j++
			if state != 2 {
				reviewRounds++
			}
			state = 2
		}
	}

	if state == 1 {
		reviewRounds++
	} else if i < len(commits) {
		reviewRounds++
	}

	return reviewRounds
}

func setReviewRounds(mr *gitlabModels.GitlabMergeRequest) error {
	var notes []gitlabModels.GitlabMergeRequestNote

	lakeModels.Db.Where("merge_request_id = ? AND `system` = 0", mr.GitlabId).Order("gitlab_created_at asc").Find(&notes)

	var commits []gitlabModels.GitlabMergeRequestCommit
	lakeModels.Db.Joins("join gitlab_merge_request_commit_merge_requests gmrcmr on gmrcmr.merge_request_commit_id = gitlab_merge_request_commits.commit_id").Where("merge_request_id = ?", mr.GitlabId).Order("authored_date asc").Find(&commits)

	reviewRounds := GetReviewRounds(commits, notes)

	err := lakeModels.Db.Model(&mr).Where("gitlab_id = ?", mr.GitlabId).Update("review_rounds", reviewRounds).Error

	if err != nil {
		return nil
	}

	return nil
}

func EnrichMergeRequests() error {
	// get mrs from theDB
	var mrs []gitlabModels.GitlabMergeRequest
	lakeModels.Db.Find(&mrs)

	for _, mr := range mrs {
		setReviewRounds(&mr)
	}
	return nil
}
