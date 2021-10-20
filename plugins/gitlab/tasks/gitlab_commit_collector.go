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

type ApiCommitResponse []GitlabApiCommit

type GitlabApiCommit struct {
	GitlabId       string `json:"id"`
	Title          string
	Message        string
	ProjectId      int
	ShortId        string `json:"short_id"`
	AuthorName     string `json:"author_name"`
	AuthorEmail    string `json:"author_email"`
	AuthoredDate   string `json:"authored_date"`
	CommitterName  string `json:"committer_name"`
	CommitterEmail string `json:"committer_email"`
	CommittedDate  string `json:"committed_date"`
	WebUrl         string `json:"web_url"`
	Stats          struct {
		Additions int
		Deletions int
		Total     int
	}
}

func CollectCommits(projectId int, scheduler *utils.WorkerScheduler) error {
	gitlabApiClient := CreateApiClient()

	return gitlabApiClient.FetchWithPaginationAnts(scheduler, fmt.Sprintf("projects/%v/repository/commits?with_stats=true", projectId), 100,
		func(res *http.Response) error {

			gitlabApiResponse := &ApiCommitResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
			for _, gitlabApiCommit := range *gitlabApiResponse {
				gitlabCommit, err := convertCommit(&gitlabApiCommit, projectId)
				if err != nil {
					return err
				}

				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&gitlabCommit).Error

				if err != nil {
					logger.Error("Could not upsert: ", err)
				}
			}

			return nil
		})
}

// Convert the API response to our DB model instance
func convertCommit(commit *GitlabApiCommit, projectId int) (*models.GitlabCommit, error) {
	convertedAuthoredDate, err := utils.ConvertStringToTime(commit.AuthoredDate)
	if err != nil {
		logger.Error("Error >>> authored date must be valid: ", err)
		return nil, err
	}
	convertedCommittedDate, err := utils.ConvertStringToTime(commit.CommittedDate)
	if err != nil {
		logger.Error("Error >>> committed date must be valid: ", err)
		return nil, err
	}
	gitlabCommit := &models.GitlabCommit{
		GitlabId:       commit.GitlabId,
		Title:          commit.Title,
		Message:        commit.Message,
		ProjectId:      projectId,
		ShortId:        commit.ShortId,
		AuthorName:     commit.AuthorName,
		AuthorEmail:    commit.AuthorEmail,
		AuthoredDate:   *convertedAuthoredDate,
		CommitterName:  commit.CommitterName,
		CommitterEmail: commit.CommitterEmail,
		CommittedDate:  *convertedCommittedDate,
		WebUrl:         commit.WebUrl,
		Additions:      commit.Stats.Additions,
		Deletions:      commit.Stats.Deletions,
		Total:          commit.Stats.Total,
	}
	return gitlabCommit, nil
}
