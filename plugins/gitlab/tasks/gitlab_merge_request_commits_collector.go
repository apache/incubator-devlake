package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/utils"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

type ApiMergeRequestCommitResponse []GitlabApiCommit

func CollectMergeRequestCommits(ctx context.Context, projectId int, rateLimitPerSecondInt int, gitlabApiClient *GitlabApiClient) error {
	scheduler, err := utils.NewWorkerScheduler(rateLimitPerSecondInt*2, rateLimitPerSecondInt, ctx)
	if err != nil {
		return nil
	}
	defer scheduler.Release()

	cursor, err := lakeModels.Db.Model(&models.GitlabMergeRequest{}).Where("project_id = ?", projectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	gitlabMr := &models.GitlabMergeRequest{}
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, gitlabMr)
		if err != nil {
			return nil
		}
		getUrl := fmt.Sprintf("projects/%v/merge_requests/%v/commits", projectId, gitlabMr.Iid)
		err = scheduler.Submit(func() error {
			return gitlabApiClient.FetchWithPagination(getUrl, nil, 100,
				func(res *http.Response) error {
					gitlabApiResponse := &ApiMergeRequestCommitResponse{}
					err = core.UnmarshalResponse(res, gitlabApiResponse)

					if err != nil {
						logger.Error("Error: ", err)
						return err
					}

					for _, commit := range *gitlabApiResponse {
						gitlabCommit, err := ConvertCommit(&commit)
						if err != nil {
							return err
						}

						err = lakeModels.Db.Clauses(clause.OnConflict{
							UpdateAll: true,
						}).Create(&gitlabCommit).Error
						if err != nil {
							logger.Error("Could not upsert: ", err)
							return err
						}

						gitlabMrCommit := &models.GitlabMergeRequestCommit{
							CommitSha:      commit.GitlabId,
							MergeRequestId: gitlabMr.GitlabId,
						}
						err = lakeModels.Db.Clauses(clause.OnConflict{
							UpdateAll: true,
						}).Create(&gitlabMrCommit).Error

						if err != nil {
							logger.Error("Could not upsert: ", err)
							return err
						}
					}
					return nil
				})
		})
		if err != nil {
			return nil
		}
	}
	scheduler.WaitUntilFinish()
	return nil
}
