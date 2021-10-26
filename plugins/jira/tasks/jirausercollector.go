package tasks

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

type JiraUserApiRes []JiraApiUser
type JiraApiUser struct {
	Name       string `json:"displayName"`
	Email      string `json:"emailAddress"`
	Timezone   string `json:"timeZone"`
	AvatarUrls struct {
		Url string `json:"48x48"`
	} `json:"avatarUrls"`
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

	for _, project := range jiraProjects {
		query := &url.Values{}
		query.Set("project", project.Key)
		err := jiraApiClient.FetchWithoutPagination("/rest/api/3/user/assignable/search", query,
			func(res *http.Response) error {
				jiraApiUsersResponse := &JiraUserApiRes{}
				err := core.UnmarshalResponse(res, jiraApiUsersResponse)
				if err != nil {
					return err
				}

				// there is no more data to fetch
				if len(*jiraApiUsersResponse) == 0 {
					return errors.New("Done fetching")
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

// convert api response into the model for the db
func convertUser(user *JiraApiUser, projectId string) (*models.JiraUser, error) {
	jiraUser := &models.JiraUser{
		ProjectId: projectId,
		Name:      user.Name,
		Email:     user.Email,
		Timezone:  user.Timezone,
		AvatarUrl: user.AvatarUrls.Url,
	}
	return jiraUser, nil
}
