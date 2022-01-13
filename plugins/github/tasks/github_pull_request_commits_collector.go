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
		Date  core.Iso8601Time
	}
	Committer struct {
		Name  string
		Email string
		Date  core.Iso8601Time
	}
	Message string
}

var pullRequestCommitSlice []models.GithubPullRequestCommit
var pullRequestAssociationSlice []models.GithubPullRequestCommitPullRequest

func CollectPullRequestCommits(owner string, repositoryName string, pull *models.GithubPullRequest, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/pulls/%v/commits", owner, repositoryName, pull.Number)
	return githubApiClient.FetchWithPaginationAnts(getUrl, nil, 100, 1, scheduler,
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

					pullRequestCommitSlice = append(pullRequestCommitSlice, *githubCommit)

					GithubPullRequestCommitPullRequest := &models.GithubPullRequestCommitPullRequest{
						PullRequestCommitSha: prCommit.Sha,
						PullRequestId:        pull.GithubId,
					}

					pullRequestAssociationSlice = append(pullRequestAssociationSlice, *GithubPullRequestCommitPullRequest)
				}
			} else {
				fmt.Println("INFO: PR PrCommit collection >>> res.Status: ", res.Status)
			}
			err := savePullRequestCommitsInBatches()
			if err != nil {
				return err
			}
			return nil
		})
}

func savePullRequestCommitsInBatches() error {
	err := lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&pullRequestCommitSlice).Error
	if err != nil {
		logger.Error("Could not upsert: ", err)
		return err
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&pullRequestAssociationSlice).Error
	if err != nil {
		logger.Error("Could not upsert: ", err)
		return err
	}
	return nil
}

func convertPullRequestCommit(prCommit *PrCommitsResponse, pullId int) (*models.GithubPullRequestCommit, error) {
	githubCommit := &models.GithubPullRequestCommit{
		Sha:            prCommit.Sha,
		PullRequestId:  pullId,
		Message:        prCommit.Commit.Message,
		AuthorName:     prCommit.Commit.Author.Name,
		AuthorEmail:    prCommit.Commit.Author.Email,
		AuthoredDate:   prCommit.Commit.Author.Date.ToTime(),
		CommitterName:  prCommit.Commit.Committer.Name,
		CommitterEmail: prCommit.Commit.Committer.Email,
		CommittedDate:  prCommit.Commit.Committer.Date.ToTime(),
		Url:            prCommit.Url,
	}
	return githubCommit, nil
}
