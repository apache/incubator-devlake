package tasks

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

type ApiTagsResponse []GitlabApiTag

type GitlabApiTag struct {
	Name      string
	Message   string
	Target    string
	Protected bool
	Release   struct {
		TagName     string
		Description string
	}
}

func CollectTags(ctx context.Context, projectId int, gitlabApiClient *GitlabApiClient) error {
	relativePath := fmt.Sprintf("projects/%v/repository/tags", projectId)
	queryParams := url.Values{}
	queryParams.Set("with_stats", "true")
	return gitlabApiClient.FetchWithPaginationAnts(relativePath, queryParams, 100,
		func(res *http.Response) error {

			gitlabApiResponse := &ApiTagsResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
			for _, gitlabApiTag := range *gitlabApiResponse {
				gitlabTag, err := convertTag(&gitlabApiTag)
				if err != nil {
					return err
				}

				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&gitlabTag).Error

				if err != nil {
					logger.Error("Could not upsert: ", err)
					return err
				}
			}

			return nil
		})
}

// Convert the API response to our DB model instance
func convertTag(tag *GitlabApiTag) (*models.GitlabTag, error) {
	gitlabTag := &models.GitlabTag{
		Name:               tag.Name,
		Message:            tag.Message,
		Target:             tag.Target,
		Protected:          tag.Protected,
		ReleaseDescription: tag.Release.Description,
	}
	return gitlabTag, nil
}
