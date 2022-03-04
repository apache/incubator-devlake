package tasks

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
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

func CollectCommits(ctx context.Context, projectId int, gitlabApiClient *GitlabApiClient) error {
	relativePath := fmt.Sprintf("projects/%v/repository/commits", projectId)
	queryParams := &url.Values{}
	queryParams.Set("with_stats", "true")
	gitlabUser := &models.GitlabUser{}
	return gitlabApiClient.FetchWithPaginationAnts(relativePath, queryParams, 100,
		func(res *http.Response) error {

			gitlabApiResponse := &ApiCommitResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
			gitlabProjectCommit := &models.GitlabProjectCommit{GitlabProjectId: projectId}
			for _, gitlabApiCommit := range *gitlabApiResponse {
				gitlabCommit, err := ConvertCommit(&gitlabApiCommit)
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

				// create project/commits relationship
				gitlabProjectCommit.CommitSha = gitlabCommit.Sha
				err = lakeModels.Db.Clauses(clause.OnConflict{
					DoNothing: true,
				}).Create(&gitlabProjectCommit).Error
				if err != nil {
					logger.Error("Could not upsert: ", err)
					return err
				}

				// create gitlab user
				gitlabUser.Email = gitlabCommit.AuthorEmail
				gitlabUser.Name = gitlabCommit.AuthorName
				err = lakeModels.Db.Clauses(clause.OnConflict{
					DoNothing: true,
				}).Create(&gitlabUser).Error
				if err != nil {
					logger.Error("Could not upsert: ", err)
					return err
				}
				if gitlabCommit.CommitterEmail != gitlabUser.Email {
					gitlabUser.Email = gitlabCommit.CommitterEmail
					gitlabUser.Name = gitlabCommit.CommitterName
					err = lakeModels.Db.Clauses(clause.OnConflict{
						DoNothing: true,
					}).Create(&gitlabUser).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
						return err
					}
				}
			}

			return nil
		})
}

// Convert the API response to our DB model instance
func ConvertCommit(commit *GitlabApiCommit) (*models.GitlabCommit, error) {
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
