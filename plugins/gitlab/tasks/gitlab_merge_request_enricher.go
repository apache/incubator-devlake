package tasks

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
)

func calculateReviewRounds(mr *gitlabModels.GitlabMergeRequest) error {
	var notes []gitlabModels.GitlabMergeRequestNote
	tx := lakeModels.Db.Debug().Where("merge_request_id = ?", mr.Iid).Find(&notes)
	logger.Info("Find Notes -> Rows Affected: ", tx.RowsAffected)
	fmt.Println("KEVIN >>> mr.id: ", mr.GitlabId)
	fmt.Println("KEVIN >>> Notes count: ", len(notes))
	var mrCommits []gitlabModels.GitlabMergeRequestCommit
	tx = lakeModels.Db.Debug().Where("merge_request_id = ?", mr.GitlabId).Find(&mrCommits)
	logger.Info("Find MR_Commits -> Rows Affected: ", tx.RowsAffected)
	// var commits []gitlabModels.GitlabCommit
	// for _, mrCommit := range mrCommits {
	// 	tx = lakeModels.Db.Debug().Where("gitlab_id = ?", mrCommit.CommitId).Find(&commits)
	// 	logger.Info("Find Commits -> Rows Affected: ", tx.RowsAffected)
	// }
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
