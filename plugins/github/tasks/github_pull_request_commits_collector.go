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

func CollectPullRequestCommits(owner string, repositoryName string, pull *models.GithubPullRequest, scheduler *utils.WorkerScheduler) error {
	githubApiClient := CreateApiClient()
	getUrl := fmt.Sprintf("repos/%v/%v/pulls/%v/commits", owner, repositoryName, pull.Number)
	return githubApiClient.FetchWithPaginationAnts(getUrl, 100, 1, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiPullRequestCommitResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, prCommit := range *githubApiResponse {
					githubCommit, err := convertPullRequestCommit(&prCommit, pull.GithubId)
					if err != nil {
						return err
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubCommit).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
					GithubPullRequestCommitPullRequest := &models.GithubPullRequestCommitPullRequest{
						PullRequestCommitSha: prCommit.Sha,
						PullRequestId:        pull.GithubId,
					}
					result := lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&GithubPullRequestCommitPullRequest)

					if result.Error != nil {
						logger.Error("Could not upsert: ", result.Error)
					}
				}
			} else {
				fmt.Println("INFO: PR PrCommit collection >>> res.Status: ", res.Status)
			}
			return nil
		})
}
func convertPullRequestCommit(prCommit *PrCommitsResponse, pullId int) (*models.GithubPullRequestCommit, error) {
	convertedAuthoredDate, err := utils.ConvertStringToTime(prCommit.Commit.Author.Date)
	if err != nil {
		return nil, err
	}
	convertedCommittedDate, err := utils.ConvertStringToTime(prCommit.Commit.Committer.Date)
	if err != nil {
		return nil, err
	}
	githubCommit := &models.GithubPullRequestCommit{
		Sha:            prCommit.Sha,
		PullRequestId:  pullId,
		Message:        prCommit.Commit.Message,
		AuthorName:     prCommit.Commit.Author.Name,
		AuthorEmail:    prCommit.Commit.Author.Email,
		AuthoredDate:   *convertedAuthoredDate,
		CommitterName:  prCommit.Commit.Committer.Name,
		CommitterEmail: prCommit.Commit.Committer.Email,
		CommittedDate:  *convertedCommittedDate,
		Url:            prCommit.Url,
	}
	return githubCommit, nil
}
