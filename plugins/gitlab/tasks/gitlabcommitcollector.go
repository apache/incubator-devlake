package tasks

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
)

type ApiCommitResponse []struct {
	Title string `json:"title"`
	// Message string `json:"message"`
}

func createApiClient() *core.ApiClient {
	return core.NewApiClient(
		config.V.GetString("GITLAB_ENDPOINT"),
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", config.V.GetString("GITLAB_AUTH")),
		},
		10*time.Second,
		3,
	)
}

func CollectCommits(projectId int) error {
	gitlabApiClient := createApiClient()

	res, err := gitlabApiClient.Get(fmt.Sprintf("projects/%v/repository/commits?with_stats=true", projectId), nil, nil)
	if err != nil {
		return err
	}

	gitlabApiResponse := &ApiCommitResponse{}

	logger.Info("res", res)

	err = core.UnmarshalResponse(res, gitlabApiResponse)

	if err != nil {
		logger.Error("Error: ", err)
		return nil
	}

	for _, value := range *gitlabApiResponse {
		fmt.Println(value.Title)
		gitlabCommit := &models.GitlabCommit{
			Title: value.Title,
		}
		err = lakeModels.Db.Save(gitlabCommit).Error
	}

	if err != nil {
		logger.Error("Error: ", err)
	}
	return nil
}
