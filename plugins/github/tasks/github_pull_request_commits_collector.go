package tasks

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

type ApiPullRequestCommitResponse []PrCommitsResponse
type PrCommitsResponse struct {
	Sha    string `json:"sha"`
	Commit PrCommit
	Url    string
}

type PrCommit struct {
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

func CollectPullRequestCommits(pull *models.GithubPullRequest) error {
	githubApiClient := CreateApiClient()
	getUrl := strings.Replace(pull.CommitsUrl, config.V.GetString("GITHUB_ENDPOINT"), "", 1)
	return githubApiClient.FetchWithPaginationAnts(getUrl, 100,
		func(res *http.Response) error {
			githubApiResponse := &ApiPullRequestCommitResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, prCommit := range *githubApiResponse {
					githubCommit := &models.GithubPullRequestCommit{
						Sha:            prCommit.Sha,
						PullRequestId:  pull.GithubId,
						Message:        prCommit.Commit.Message,
						AuthorName:     prCommit.Commit.Author.Name,
						AuthorEmail:    prCommit.Commit.Author.Email,
						AuthoredDate:   prCommit.Commit.Author.Date,
						CommitterName:  prCommit.Commit.Committer.Name,
						CommitterEmail: prCommit.Commit.Committer.Email,
						CommittedDate:  prCommit.Commit.Committer.Date,
						Url:            prCommit.Url,
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubCommit).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				}
			} else {
				fmt.Println("INFO: PR PrCommit collection >>> res.Status: ", res.Status)
			}
			return nil
		})
}
