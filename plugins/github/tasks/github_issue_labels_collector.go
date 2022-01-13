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

type ApiIssueLabelResponse []IssueLabel

type IssueLabel struct {
	GithubId    int `json:"id"`
	Name        string
	Description string
	Color       string
}

var issueLabelsSlice []models.GithubIssueLabelIssue

func CollectIssueLabelsForSinglePullRequest(owner string, repositoryName string, pr *models.GithubPullRequest, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/issues/%v/labels", owner, repositoryName, pr.Number)
	return githubApiClient.FetchWithPaginationAnts(getUrl, nil, 100, 1, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiIssueLabelResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, label := range *githubApiResponse {
					githubLabel := &models.GithubIssueLabelIssue{
						IssueId:      pr.GithubId,
						IssueLabelId: label.GithubId,
					}
					issueLabelsSlice = append(issueLabelsSlice, *githubLabel)
				}
			} else {
				fmt.Println("INFO: PR Label collection >>> res.Status: ", res.Status)
			}
			err := saveLabelsInBatches()
			if err != nil {
				return err
			}
			return nil
		})
}

func saveLabelsInBatches() error {
	err := lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&issueLabelsSlice).Error
	if err != nil {
		logger.Error("Could not upsert: ", err)
		return err
	} else {
		return nil
	}
}

func CollectIssueLabelsForSingleIssue(owner string, repositoryName string, issue *models.GithubIssue, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/issues/%v/labels", owner, repositoryName, issue.Number)
	return githubApiClient.FetchWithPaginationAnts(getUrl, nil, 100, 1, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiIssueLabelResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, label := range *githubApiResponse {
					githubLabel := &models.GithubIssueLabelIssue{
						IssueId:      issue.GithubId,
						IssueLabelId: label.GithubId,
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubLabel).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				}
			} else {
				fmt.Println("INFO: PR Label collection >>> res.Status: ", res.Status)
			}
			return nil
		})
}
func CollectRepositoryIssueLabels(owner string, repositoryName string, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/labels", owner, repositoryName)
	return githubApiClient.FetchWithPaginationAnts(getUrl, nil, 100, 1, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiIssueLabelResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, label := range *githubApiResponse {
					githubLabel := &models.GithubIssueLabel{
						GithubId:    label.GithubId,
						Name:        label.Name,
						Description: label.Description,
						Color:       label.Color,
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubLabel).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				}
			} else {
				fmt.Println("INFO: PR Label collection >>> res.Status: ", res.Status)
			}
			return nil
		})
}
