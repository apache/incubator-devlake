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
	Commit PullRequestCommit
	Url    string
}

type PullRequestCommit struct {
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

func CollectPullRequestCommits(owner string, repo string, scheduler *utils.WorkerScheduler, apiClient *GithubApiClient) error {
	var prs []models.GithubPullRequest
	lakeModels.Db.Find(&prs)
	for i := 0; i < len(prs); i++ {
		err := ProcessCollection(owner, repo, &prs[i], scheduler, apiClient)
		if err != nil {
			return err
		}
	}
	return nil
}
func convertPullRequestCommit(prCommit *PrCommitsResponse) (*models.GithubCommit, error) {
	githubCommit := &models.GithubCommit{
		Sha:            prCommit.Sha,
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

func ProcessCollection(owner string, repo string, pr *models.GithubPullRequest, scheduler *utils.WorkerScheduler, apiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/pulls/%v/commits", owner, repo, pr.Number)
	return apiClient.FetchWithPaginationAnts(getUrl, nil, 100, 1, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiPullRequestCommitResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				err = lakeModels.Db.Where("pull_request_id = ?",
					pr.GithubId).Delete(&models.GithubPullRequestCommit{}).Error
				if err != nil {
					logger.Error("Could not delete: ", err)
					return err
				}
				for _, pullRequestCommit := range *githubApiResponse {
					githubCommit, err := convertPullRequestCommit(&pullRequestCommit)
					if err != nil {
						return err
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						DoNothing: true,
					}).Create(&githubCommit).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
					githubPullRequestCommit := &models.GithubPullRequestCommit{
						CommitSha:     pullRequestCommit.Sha,
						PullRequestId: pr.GithubId,
					}
					result := lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubPullRequestCommit)

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
