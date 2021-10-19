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
	Sha    string `json:"sha"`
	Commit Commit
	Url    string
}

type Commit struct {
	Author struct {
		Name  string
		Email string
		Date  string
	}
	Committer struct {
		Name  string
		Email string
		Date  string
	}
	Message string
}

func CollectCommits(owner string, repositoryName string, repositoryId int, scheduler *utils.WorkerScheduler) error {
	githubApiClient := CreateApiClient()
	getUrl := fmt.Sprintf("repos/%v/%v/commits", owner, repositoryName)
	return githubApiClient.FetchWithPaginationAnts(getUrl, 100, 20, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiCommitsResponse{}
			err := core.UnmarshalResponse(res, githubApiResponse)
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}
			for _, commit := range *githubApiResponse {
				githubCommit := &models.GithubCommit{
					Sha:            commit.Sha,
					RepositoryId:   repositoryId,
					Message:        commit.Commit.Message,
					AuthorName:     commit.Commit.Author.Name,
					AuthorEmail:    commit.Commit.Author.Email,
					AuthoredDate:   utils.ConvertStringToTime(commit.Commit.Author.Date),
					CommitterName:  commit.Commit.Committer.Name,
					CommitterEmail: commit.Commit.Committer.Email,
					CommittedDate:  utils.ConvertStringToTime(commit.Commit.Committer.Date),
					Url:            commit.Url,
				}
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&githubCommit).Error
				if err != nil {
					logger.Error("Could not upsert: ", err)
				}
			}
			return nil
		})
}
