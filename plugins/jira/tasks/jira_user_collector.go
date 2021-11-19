package tasks

import (
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

type JiraUserApiRes []JiraApiUser
type JiraApiUser struct {
	AccountId   string `json:"accountId"`
	AccountType string `json:"accountType"`
	DisplayName string `json:"displayName"`
	Email       string `json:"emailAddress"`
	Timezone    string `json:"timeZone"`
	AvatarUrls  struct {
		Url string `json:"48x48"`
	} `json:"avatarUrls"`
}

func CollectUsers(jiraApiClient *JiraApiClient,
	sourceId uint64,
) error {
	// The reason we use FetchWithoutPaginationHeaders is because this API endpoint does not
	// return pagination info in it's headers the same way that other endpoints do.
	// This method still uses pagination, but in a different way.
	err := jiraApiClient.FetchWithoutPaginationHeaders("rest/api/3/users/search", nil,
		func(res *http.Response) (int, error) {
			jiraApiUsersResponse := &JiraUserApiRes{}
			err := core.UnmarshalResponse(res, jiraApiUsersResponse)
			if err != nil {
				return 0, err
			}

			// there is no more data to fetch
			if len(*jiraApiUsersResponse) == 0 {
				return 0, nil
			}

			// process Users
			for _, jiraApiUser := range *jiraApiUsersResponse {
				jiraUser, err := convertUser(&jiraApiUser, sourceId)
				if err != nil {
					return 0, err
				}
				// User
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(jiraUser).Error
				if err != nil {
					logger.Info("Error saving user", jiraUser)
					return 0, err
				}
			}
			return len(*jiraApiUsersResponse), nil
		})
	if err != nil {
		return err
	}
	return nil
}

// convert api response into the model for the db
func convertUser(user *JiraApiUser, sourceId uint64) (*models.JiraUser, error) {
	jiraUser := &models.JiraUser{
		SourceId:    sourceId,
		AccountId:   user.AccountId,
		AccountType: user.AccountType,
		Name:        user.DisplayName,
		Email:       user.Email,
		Timezone:    user.Timezone,
		AvatarUrl:   user.AvatarUrls.Url,
	}
	return jiraUser, nil
}
