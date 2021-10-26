package tasks

import (
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/logger"
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
	sourceId uint64,
) error {
	var jiraProjects []models.JiraProject
	err := lakeModels.Db.Debug().Find(&jiraProjects).Error
	if err != nil {
		return err
	}

	for _, project := range jiraProjects {
		query := &url.Values{}
		query.Set("project", project.Key)
		err := jiraApiClient.FetchWithoutPagination("/api/3/user/assignable/search", query,
			func(res *http.Response) (bool, error) {
				jiraApiUsersResponse := &JiraUserApiRes{}
				err := core.UnmarshalResponse(res, jiraApiUsersResponse)
				if err != nil {
					return false, err
				}

				// there is no more data to fetch
				if len(*jiraApiUsersResponse) == 0 {
					return false, nil
				}

				// process Users
				for _, jiraApiUser := range *jiraApiUsersResponse {
					jiraUser, err := convertUser(&jiraApiUser, project.Id, sourceId)
					if err != nil {
						return false, err
					}
					// User
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(jiraUser).Error
					if err != nil {
						logger.Info("Error saving user", jiraUser)
						return false, err
					}
				}
				return true, nil
			})
		if err != nil {
			return err
		}
	}
	return nil
}

// convert api response into the model for the db
func convertUser(user *JiraApiUser, projectId string, sourceId uint64) (*models.JiraUser, error) {
	jiraUser := &models.JiraUser{
		SourceId:  sourceId,
		ProjectId: projectId,
		Name:      user.Name,
		Email:     user.Email,
		Timezone:  user.Timezone,
		AvatarUrl: user.AvatarUrls.Url,
	}
	return jiraUser, nil
}
