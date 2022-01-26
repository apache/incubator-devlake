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

func CollectIssueLabels(owner string, repositoryName string, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
	var issues []models.GithubIssue
	lakeModels.Db.Find(&issues)
	for i := 0; i < len(issues); i++ {
		labelsErr := processIssueLabelsCollection(owner, repositoryName, &issues[i], scheduler, githubApiClient)
		if labelsErr != nil {
			logger.Error("Could not collect issue labels", labelsErr)
			return labelsErr
		}
	}
	return nil
}

func processIssueLabelsCollection(owner string, repositoryName string, issue *models.GithubIssue, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
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
				//delete labels before insertion to keep the table having no outdated data
				err = lakeModels.Db.Where("issue_id = ?",
					issue.GithubId).Delete(&models.GithubIssueLabel{}).Error
				if err != nil {
					logger.Error("Could not delete: ", err)
					return err
				}
				for _, label := range *githubApiResponse {
					githubLabel := &models.GithubIssueLabel{
						IssueId:   issue.GithubId,
						LabelName: label.Name,
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
