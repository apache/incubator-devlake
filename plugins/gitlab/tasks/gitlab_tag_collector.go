package tasks

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
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
	Commit struct {
		Id string
	}
}

func CollectTags(projectId int, gitlabApiClient *GitlabApiClient) error {
	relativePath := fmt.Sprintf("projects/%v/repository/tags", projectId)
	queryParams := &url.Values{}
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
				gitlabTag, err := convertTag(&gitlabApiTag, projectId)
				if err != nil {
					return err
				}

				err = lakeModels.Db.
					Where(`project_id=? AND name=?`, projectId, gitlabTag.Name).
					Delete(&models.GitlabTag{}).Error
				if err != nil {
					return err
				}

				err = lakeModels.Db.Create(&gitlabTag).Error
				if err != nil {
					logger.Error("Could not upsert: ", err)
					return err
				}
			}

			return nil
		})
}

// Convert the API response to our DB model instance
func convertTag(tag *GitlabApiTag, projectId int) (*models.GitlabTag, error) {
	gitlabTag := &models.GitlabTag{
		ProjectId:          projectId,
		Name:               tag.Name,
		Message:            tag.Message,
		Target:             tag.Commit.Id,
		Protected:          tag.Protected,
		ReleaseDescription: tag.Release.Description,
	}
	return gitlabTag, nil
}
