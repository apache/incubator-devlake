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

type ApiPrLabelResponse []IssueLabel

type PrLabel struct {
	GithubId    int `json:"id"`
	Name        string
	Description string
	Color       string
}

func CollectPrLabels(owner string, repositoryName string, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
	var prs []models.GithubPullRequest
	lakeModels.Db.Find(&prs)
	for i := 0; i < len(prs); i++ {
		labelsErr := processPrLabelsCollection(owner, repositoryName, &prs[i], scheduler, githubApiClient)
		if labelsErr != nil {
			logger.Error("Could not collect Pr labels", labelsErr)
			return labelsErr
		}
	}
	return nil
}

func processPrLabelsCollection(owner string, repositoryName string, pr *models.GithubPullRequest, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
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
					githubIssueLabel := &models.GithubIssueLabel{
						IssueId:        pr.GithubId,
						IssueLabelName: label.Name,
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubIssueLabel).Error
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
