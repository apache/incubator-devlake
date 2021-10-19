package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
)

func GetReviewRounds(commits []gitlabModels.GitlabMergeRequestCommit, notes []gitlabModels.GitlabMergeRequestNote) int {
	i := 0
	j := 0
	reviewRounds := 0

	if len(commits) == 0 && len(notes) == 0 {
		return 1
	}

	// state is used to keep track of previous activity
	// 0: init, 1: commit, 2: comment
	// whenever state is switched to comment, we increment reviewRounds by 1
	state := 0 // 0, 1, 2
	for i < len(commits) && j < len(notes) {
		if commits[i].AuthoredDate.Before(notes[j].GitlabCreatedAt) {
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

	// There's another implicit round of review in 2 scenarios
	// One: the last state is commit (state == 1)
	// Two: the last state is comment but there're still commits left
	if state == 1 || i < len(commits) {
		reviewRounds++
	}

	return reviewRounds
}

func setReviewRounds(mr *gitlabModels.GitlabMergeRequest) error {
	var notes []gitlabModels.GitlabMergeRequestNote

	// `system` = 0 is needed since we only care about human comments
	lakeModels.Db.Where("merge_request_id = ? AND `system` = 0", mr.GitlabId).Order("gitlab_created_at asc").Find(&notes)

	var commits []gitlabModels.GitlabMergeRequestCommit
	lakeModels.Db.Joins("join gitlab_merge_request_commit_merge_requests gmrcmr on gmrcmr.merge_request_commit_id = gitlab_merge_request_commits.commit_id").Where("merge_request_id = ?", mr.GitlabId).Order("authored_date asc").Find(&commits)

	reviewRounds := GetReviewRounds(commits, notes)

	err := lakeModels.Db.Model(&mr).Where("gitlab_id = ?", mr.GitlabId).Update("review_rounds", reviewRounds).Error

	if err != nil {
		return err
	}

	return nil
}

func EnrichMergeRequests() error {
	// get mrs from theDB
	var mrs []gitlabModels.GitlabMergeRequest
	lakeModels.Db.Find(&mrs)

	for _, mr := range mrs {
		err := setReviewRounds(&mr)
		if err != nil {
			return err
		}
	}
	return nil
}
