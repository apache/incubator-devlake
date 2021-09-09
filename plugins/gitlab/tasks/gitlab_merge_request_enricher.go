package tasks

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
)

func calculateReviewRounds(mr *gitlabModels.GitlabMergeRequest) error {
	var notes []gitlabModels.GitlabMergeRequestNote
	mrId := 77373958
	tx := lakeModels.Db.Debug().Where("merge_request_id = ?", mrId).Find(&notes)
	logger.Info("Find Notes -> Rows Affected: ", tx.RowsAffected)
	fmt.Println("KEVIN >>> Notes : ", notes)

	var commits []gitlabModels.GitlabCommit
	tx = lakeModels.Db.Debug().Joins("join gitlab_merge_request_commits on gitlab_merge_request_commits.commit_id = gitlab_commits.gitlab_id").Where("merge_request_id = ?", mrId).Find(&commits)
	logger.Info("Find MR_Commits -> Rows Affected: ", tx.RowsAffected)
	fmt.Println("KEVIN >>> Commits: ", commits)

	return nil
}

func EnrichMergeRequests() error {
	var mrs []gitlabModels.GitlabMergeRequest
	tx := lakeModels.Db.Find(&mrs)
	logger.Info("Find MRs -> Rows Affected: ", tx.RowsAffected)
	calculateReviewRounds(&mrs[0])
	// for i := 0; i < len(mrs); i++ {
	// 	calculateReviewRounds(&mrs[i])
	// }
	return nil
}
