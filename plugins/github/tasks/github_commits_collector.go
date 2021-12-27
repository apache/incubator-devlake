package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

type ApiCommitsResponse []CommitsResponse
type CommitsResponse struct {
	Sha       string `json:"sha"`
	Commit    Commit
	Url       string
	Author    *models.GithubUser
	Committer *models.GithubUser
}

type Commit struct {
	Author struct {
		Name  string
		Email string
		Date  core.Iso8601Time
	}
	Committer struct {
		Name  string
		Email string
		Date  core.Iso8601Time
	}
	Message string
}

func CollectCommits(owner string, repositoryName string, repositoryId int, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/commits", owner, repositoryName)
	return githubApiClient.FetchWithPaginationAnts(getUrl, nil, 100, 20, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiCommitsResponse{}
			err := core.UnmarshalResponse(res, githubApiResponse)
			if err != nil || res.StatusCode == 401 {
				logger.Error("Error: ", err)
				return err
			}
			for _, commit := range *githubApiResponse {
				githubCommit, err := convertGithubCommit(&commit, repositoryId)
				if err != nil {
					return err
				}
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&githubCommit).Error
				if err != nil {
					logger.Error("Could not upsert: ", err)
				}
				// save author and committer
				if commit.Author != nil {
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&commit.Author).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				}
				if commit.Committer != nil {
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&commit.Committer).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				}
			}
			return nil
		})
}
func convertGithubCommit(commit *CommitsResponse, repoId int) (*models.GithubCommit, error) {
	githubCommit := &models.GithubCommit{
		Sha:            commit.Sha,
		RepositoryId:   repoId,
		Message:        commit.Commit.Message,
		AuthorId:       commit.Author.Id,
		AuthorName:     commit.Commit.Author.Name,
		AuthorEmail:    commit.Commit.Author.Email,
		AuthoredDate:   commit.Commit.Author.Date.ToTime(),
		CommitterId:    commit.Committer.Id,
		CommitterName:  commit.Commit.Committer.Name,
		CommitterEmail: commit.Commit.Committer.Email,
		CommittedDate:  commit.Commit.Committer.Date.ToTime(),
		Url:            commit.Url,
	}
	return githubCommit, nil
}
