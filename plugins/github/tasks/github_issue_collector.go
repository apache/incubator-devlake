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

const BatchSize = 100

type ApiIssuesResponse []IssuesResponse

type IssuesResponse struct {
	GithubId    int `json:"id"`
	Number      int
	State       string
	Title       string
	Body        string
	PullRequest struct {
		Url     string `json:"url"`
		HtmlUrl string `json:"html_url"`
	} `json:"pull_request"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`

	Assignee *struct {
		Login string
		Id    int
	}
	ClosedAt        *core.Iso8601Time `json:"closed_at"`
	GithubCreatedAt core.Iso8601Time  `json:"created_at"`
	GithubUpdatedAt core.Iso8601Time  `json:"updated_at"`
}

func CollectIssues(owner string, repo string, repoId int, scheduler *utils.WorkerScheduler, apiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/issues", owner, repo)
	queryParams := &url.Values{}
	queryParams.Set("state", "all")
	return apiClient.FetchWithPaginationAnts(getUrl, queryParams, 100, 20, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiIssuesResponse{}
			err := core.UnmarshalResponse(res, githubApiResponse)
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}

			for _, issue := range *githubApiResponse {
				if issue.GithubId == 0 {
					return nil
				}
				//If this is a pr, ignore
				if issue.PullRequest.Url != "" {
					continue
				}

				err = lakeModels.Db.Where("issue_id = ?", issue.GithubId).Delete(&models.GithubIssueLabel{}).Error
				if err != nil {
					logger.Error("delete github_issue_label error:", err)
					return err
				}
				var labels []*models.GithubIssueLabel
				for _, lable := range issue.Labels {
					labels = append(labels, &models.GithubIssueLabel{
						IssueId:   issue.GithubId,
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

				// This is an issue from github
				githubIssue, err := convertGithubIssue(&issue, repoId)
				if err != nil {
					return err
				}
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&githubIssue).Error
				if err != nil {
					logger.Error("Could not upsert: ", err)
				}

			}
			return nil
		})
}

func convertGithubIssue(issue *IssuesResponse, repositoryId int) (*models.GithubIssue, error) {
	githubIssue := &models.GithubIssue{
		GithubId:        issue.GithubId,
		RepoId:          repositoryId,
		Number:          issue.Number,
		State:           issue.State,
		Title:           issue.Title,
		Body:            issue.Body,
		ClosedAt:        core.Iso8601TimeToTime(issue.ClosedAt),
		GithubCreatedAt: issue.GithubCreatedAt.ToTime(),
		GithubUpdatedAt: issue.GithubUpdatedAt.ToTime(),
	}

	if issue.Assignee != nil {
		githubIssue.AssigneeId = issue.Assignee.Id
		githubIssue.AssigneeName = issue.Assignee.Login
	}
	if issue.ClosedAt != nil {
		githubIssue.LeadTimeMinutes = uint(issue.ClosedAt.ToTime().Sub(issue.GithubCreatedAt.ToTime()).Minutes())
	}

	return githubIssue, nil
}
