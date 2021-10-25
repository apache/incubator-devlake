package tasks

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

type JiraUserApiRes []JiraApiUser
type JiraApiUser struct {
	Name     string `json:"displayName"`
	Email    string `json:"emailAddress"`
	Timezone string `json:"timeZone"`
}

func CollectUsers(jiraApiClient *JiraApiClient,
	source *models.JiraSource,
	boardId uint64,
	ctx context.Context,
) error {
	var jiraProjects []models.JiraProject
	err := lakeModels.Db.Debug().Find(&jiraProjects).Error
	if err != nil {
		return err
	}

	scheduler, err := utils.NewWorkerScheduler(10, 50, ctx)
	if err != nil {
		return err
	}
	defer scheduler.Release()

	for _, project := range jiraProjects {
		query := &url.Values{}
		query.Set("project", project.Key)
		err := jiraApiClient.FetchWithoutPagination(scheduler, "/rest/api/3/user/assignable/search", query,
			func(res *http.Response) error {
				jiraApiUsersResponse := &JiraUserApiRes{}
				err := core.UnmarshalResponse(res, jiraApiUsersResponse)
				if err != nil {
					return err
				}

				// process Users
				for _, jiraApiUser := range *jiraApiUsersResponse {

					jiraUser, err := convertUser(&jiraApiUser, project.Id)
					if err != nil {
						return err
					}
					// User
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(jiraUser).Error
					if err != nil {
						return err
					}

					// Project / User relationship
					// lakeModels.Db.FirstOrCreate(&models.JiraProjectUser{
					// 	SourceId: source.ID,
					// 	BoardId:  boardId,
					// 	UserId:  jiraUser.UserId,
					// })
				}
				return nil
			})
		if err != nil {
			fmt.Println("KEVIN >>> err", err)
			return err
		}
	}
	return nil
}

func convertUser(user *JiraApiUser, projectId string) (*models.JiraUser, error) {
	jiraUser := &models.JiraUser{
		ProjectId: projectId,
		Name:      user.Name,
		Email:     user.Email,
		Timezone:  user.Timezone,
	}
	return jiraUser, nil
}
