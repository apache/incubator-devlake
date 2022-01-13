package tasks

import (
	"fmt"
	"net/http"
	"net/url"

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
	ShortId        string           `json:"short_id"`
	AuthorName     string           `json:"author_name"`
	AuthorEmail    string           `json:"author_email"`
	AuthoredDate   core.Iso8601Time `json:"authored_date"`
	CommitterName  string           `json:"committer_name"`
	CommitterEmail string           `json:"committer_email"`
	CommittedDate  core.Iso8601Time `json:"committed_date"`
	WebUrl         string           `json:"web_url"`
	Stats          struct {
		Additions int
		Deletions int
		Total     int
	}
}

var commitsSlice = []models.GitlabCommit{}
var projectCommitsSlice = []models.GitlabProjectCommit{}
var usersSlice = []models.GitlabUser{}

func CollectCommits(projectId int, scheduler *utils.WorkerScheduler) error {
	gitlabApiClient := CreateApiClient()
	relativePath := fmt.Sprintf("projects/%v/repository/commits", projectId)
	queryParams := &url.Values{}
	queryParams.Set("with_stats", "true")
	gitlabUser := &models.GitlabUser{}
	return gitlabApiClient.FetchWithPaginationAnts(scheduler, relativePath, queryParams, 100,
		func(res *http.Response) error {

			gitlabApiResponse := &ApiCommitResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
			gitlabProjectCommit := &models.GitlabProjectCommit{GitlabProjectId: projectId}
			for _, gitlabApiCommit := range *gitlabApiResponse {
				gitlabCommit, err := convertCommit(&gitlabApiCommit, projectId)
				if err != nil {
					return err
				}

				commitsSlice = append(commitsSlice, *gitlabCommit)

				// create project/commits relationship
				gitlabProjectCommit.CommitSha = gitlabCommit.Sha
				projectCommitsSlice = append(projectCommitsSlice, *gitlabProjectCommit)

				// create gitlab user
				gitlabUser.Email = gitlabCommit.AuthorEmail
				gitlabUser.Name = gitlabCommit.AuthorName

				usersSlice = append(usersSlice, *gitlabUser)

				if gitlabCommit.CommitterEmail != gitlabUser.Email {
					gitlabUser.Email = gitlabCommit.CommitterEmail
					gitlabUser.Name = gitlabCommit.CommitterName
					usersSlice = append(usersSlice, *gitlabUser)
				}
			}
			err = saveCommitsInBatches()
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}
			return nil
		})
}

func saveCommitsInBatches() error {
	err := lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&commitsSlice).Error
	if err != nil {
		return err
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&projectCommitsSlice).Error
	if err != nil {
		return err
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&usersSlice).Error
	if err != nil {
		return err
	}
	return nil
}

// Convert the API response to our DB model instance
func convertCommit(commit *GitlabApiCommit, projectId int) (*models.GitlabCommit, error) {
	gitlabCommit := &models.GitlabCommit{
		Sha:            commit.GitlabId,
		Title:          commit.Title,
		Message:        commit.Message,
		ShortId:        commit.ShortId,
		AuthorName:     commit.AuthorName,
		AuthorEmail:    commit.AuthorEmail,
		AuthoredDate:   commit.AuthoredDate.ToTime(),
		CommitterName:  commit.CommitterName,
		CommitterEmail: commit.CommitterEmail,
		CommittedDate:  commit.CommittedDate.ToTime(),
		WebUrl:         commit.WebUrl,
		Additions:      commit.Stats.Additions,
		Deletions:      commit.Stats.Deletions,
		Total:          commit.Stats.Total,
	}
	return gitlabCommit, nil
}
