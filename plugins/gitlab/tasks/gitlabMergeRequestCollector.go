package tasks

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

type ApiMergeRequestResponse []struct {
	GitlabId        int `json:"id"`
	Iid             int
	ProjectId       int `json:"project_id"`
	State           string
	Title           string
	Description     string
	WebUrl          string `json:"web_url"`
	UserNotesCount  int    `json:"user_notes_count"`
	WorkInProgress  bool   `json:"work_in_progress"`
	SourceBranch    string `json:"source_branch"`
	MergedAt        string `json:"merged_at"`
	GitlabCreatedAt string `json:"created_at"`
	ClosedAt        string `json:"closed_at"`
	MergedBy        struct {
		Username string `json:"username"`
	} `json:"merged_by"`
	Author struct {
		Username string `json:"username"`
	}
	Reviewers []Reviewer
}

func CollectMergeRequests(projectId int) error {
	gitlabApiClient := CreateApiClient()

	res, err := gitlabApiClient.Get(fmt.Sprintf("projects/%v/merge_requests", projectId), nil, nil)
	if err != nil {
		return err
	}

	gitlabApiResponse := &ApiMergeRequestResponse{}

	logger.Info("res", res)

	err = core.UnmarshalResponse(res, gitlabApiResponse)

	if err != nil {
		logger.Error("Error: ", err)
		return nil
	}

	go func() {

	}()
	for _, mr := range *gitlabApiResponse {
		gitlabMergeRequest := &models.GitlabMergeRequest{
			GitlabId:         mr.GitlabId,
			Iid:              mr.Iid,
			ProjectId:        mr.ProjectId,
			State:            mr.State,
			Title:            mr.Title,
			Description:      mr.Description,
			WebUrl:           mr.WebUrl,
			UserNotesCount:   mr.UserNotesCount,
			WorkInProgress:   mr.WorkInProgress,
			SourceBranch:     mr.SourceBranch,
			MergedAt:         mr.MergedAt,
			GitlabCreatedAt:  mr.GitlabCreatedAt,
			ClosedAt:         mr.ClosedAt,
			MergedByUsername: mr.MergedBy.Username,
			AuthorUsername:   mr.Author.Username,
		}

		result := lakeModels.Db.Debug().Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&gitlabMergeRequest)

		if result.Error != nil {
			logger.Error("Could not upsert: ", result.Error)
		}

		logger.Info("Rows Affected: ", result.RowsAffected)
		CreateReviewers(mr.GitlabId, mr.Reviewers)
<<<<<<< HEAD

		collectErr := CollectMergeRequestNotes(projectId, mr.Iid)

		if collectErr != nil {
			logger.Error("Could not collect MR Notes", collectErr)
		}
=======
>>>>>>> 745abc6 (feat: merge request and reviewer collection)
	}

	return nil
}
