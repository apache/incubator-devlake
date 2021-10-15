package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

type MergeRequestRes struct {
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
	Reviewers        []Reviewer
	FirstCommentTime string
}

type ApiMergeRequestResponse []MergeRequestRes

func CollectMergeRequests(projectId int, scheduler *utils.WorkerScheduler) error {
	gitlabApiClient := CreateApiClient()

	return gitlabApiClient.FetchWithPaginationAnts(scheduler, fmt.Sprintf("projects/%v/merge_requests", projectId), 100,
		func(res *http.Response) error {
			gitlabApiResponse := &ApiMergeRequestResponse{}

			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}

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
					MergedAt:         utils.ConvertStringToSqlNullTime(mr.MergedAt),
					GitlabCreatedAt:  utils.ConvertStringToSqlNullTime(mr.GitlabCreatedAt),
					ClosedAt:         utils.ConvertStringToSqlNullTime(mr.ClosedAt),
					MergedByUsername: mr.MergedBy.Username,
					AuthorUsername:   mr.Author.Username,
				}

				result := lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&gitlabMergeRequest)

				if result.Error != nil {
					logger.Error("Could not upsert: ", result.Error)
				}

				CreateReviewers(projectId, mr.GitlabId, mr.Reviewers)

			}

			return nil

		})
}
