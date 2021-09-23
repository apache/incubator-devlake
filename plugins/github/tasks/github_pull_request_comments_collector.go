package tasks

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

type ApiPullRequestCommentResponse []Comment

type Comment struct {
	GithubId int `json:"id"`
	Body     string
	User     struct {
		Login string
	}
}

func CollectPullRequestComments(pull *models.GithubPullRequest) error {
	githubApiClient := CreateApiClient()
	getUrl := strings.Replace(pull.CommentsUrl, config.V.GetString("GITHUB_ENDPOINT"), "", 1)
	return githubApiClient.FetchWithPaginationAnts(getUrl, 100,
		func(res *http.Response) error {
			githubApiResponse := &ApiPullRequestCommentResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, comment := range *githubApiResponse {
					githubComment := &models.GithubPullRequestComment{
						GithubId:       comment.GithubId,
						PullRequestId:  pull.GithubId,
						Body:           comment.Body,
						AuthorUsername: comment.User.Login,
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubComment).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				}
			} else {
				fmt.Println("INFO: PR Comment collection >>> res.Status: ", res.Status)
			}
			return nil
		})
}
