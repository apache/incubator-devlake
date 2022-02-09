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

type ApiIssueEventResponse []IssueEvent

type IssueEvent struct {
	GithubId int `json:"id"`
	Event    string
	Actor    struct {
		Login string
	}
	GithubCreatedAt core.Iso8601Time `json:"created_at"`
}

func CollectIssueEvents(owner string, repo string, scheduler *utils.WorkerScheduler, apiClient *GithubApiClient) error {
	githubIssue := &models.GithubIssue{}
	cursor, err := lakeModels.Db.Model(githubIssue).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, githubIssue)
		if err != nil {
			return err
		}
		eventsErr := processEventsCollection(owner, repo, githubIssue, scheduler, apiClient)
		if eventsErr != nil {
			logger.Error("Could not collect issue events", eventsErr)
			return eventsErr
		}
	}
	return nil
}

func processEventsCollection(owner string, repo string, issue *models.GithubIssue, scheduler *utils.WorkerScheduler, apiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/issues/%v/events", owner, repo, issue.Number)
	return apiClient.FetchWithPaginationAnts(getUrl, nil, 100, 1, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiIssueEventResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, event := range *githubApiResponse {
					githubEvent, err := convertGithubEvent(&event, issue.GithubId)
					if err != nil {
						return err
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubEvent).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				}
			} else {
				fmt.Println("INFO: PR Event collection >>> res.Status: ", res.Status)
			}
			return nil
		})
}
func convertGithubEvent(event *IssueEvent, issueId int) (*models.GithubIssueEvent, error) {
	githubEvent := &models.GithubIssueEvent{
		GithubId:        event.GithubId,
		IssueId:         issueId,
		Type:            event.Event,
		AuthorUsername:  event.Actor.Login,
		GithubCreatedAt: event.GithubCreatedAt.ToTime(),
	}
	return githubEvent, nil
}
