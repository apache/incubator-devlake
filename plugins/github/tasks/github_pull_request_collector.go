package tasks

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

type ApiPullRequestResponse []GithubApiPullRequest

type GithubApiPullRequest struct {
	GithubId int `json:"id"`
	Number   int
	State    string
	Title    string
	Body     string
	Labels   []struct {
		Name string `json:"name"`
	} `json:"labels"`
	Assignee *struct {
		Login string
		Id    int
	}
	ClosedAt        *core.Iso8601Time `json:"closed_at"`
	MergedAt        *core.Iso8601Time `json:"merged_at"`
	GithubCreatedAt core.Iso8601Time  `json:"created_at"`
	GithubUpdatedAt *core.Iso8601Time `json:"updated_at"`
	MergeCommitSha  string            `json:"merge_commit_sha"`
}

func CollectPullRequests(
	owner string,
	repo string,
	repoId int,
	scheduler *utils.WorkerScheduler,
	apiClient *GithubApiClient,
) error {
	getUrl := fmt.Sprintf("repos/%v/%v/pulls", owner, repo)
	queryParams := &url.Values{}
	queryParams.Set("state", "all")
	return apiClient.FetchWithPaginationAnts(getUrl, queryParams, 100, 20, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiPullRequestResponse{}
			err := core.UnmarshalResponse(res, githubApiResponse)
			if err != nil {
				return err
			}

			for _, pull := range *githubApiResponse {
				if pull.GithubId == 0 {
					return nil
				}
				// save pull request labels
				err = lakeModels.Db.
					Where("pull_id = ?", pull.GithubId).
					Delete(&models.GithubPullRequestLabel{}).Error
				if err != nil {
					return err
				}
				var labels []*models.GithubPullRequestLabel
				for _, lable := range pull.Labels {
					labels = append(labels, &models.GithubPullRequestLabel{
						PullId:    pull.GithubId,
						LabelName: lable.Name,
					})
				}
				err = lakeModels.Db.Clauses(clause.OnConflict{
					DoNothing: true,
				}).CreateInBatches(labels, BatchSize).Error
				if err != nil {
					logger.Error("save github_issue_label error:", err)
					return err
				}
				// save pull request detail
				githubPull, err := convertGithubPullRequest(&pull, repoId)
				if err != nil {
					return err
				}
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&githubPull).Error
				if err != nil {
					logger.Error("Could not upsert: ", err)
				}
			}
			return nil
		})
}

func convertGithubPullRequest(pull *GithubApiPullRequest, repoId int) (*models.GithubPullRequest, error) {
	githubPull := &models.GithubPullRequest{
		GithubId:        pull.GithubId,
		RepositoryId:    repoId,
		Number:          pull.Number,
		State:           pull.State,
		Title:           pull.Title,
		GithubCreatedAt: pull.GithubCreatedAt.ToTime(),
		GithubUpdatedAt: core.Iso8601TimeToTime(pull.GithubUpdatedAt),
		ClosedAt:        core.Iso8601TimeToTime(pull.ClosedAt),
		MergedAt:        core.Iso8601TimeToTime(pull.MergedAt),
		MergeCommitSha:  pull.MergeCommitSha,
	}
	return githubPull, nil
}
