package tasks

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
)

type JiraApiResponse struct {
	Issues []struct {
		Key string `json:"key"`
	} `json:"issues"`
}

func CollectIssues(boardId int) error {
	jiraApiClient := core.NewApiClient(
		config.V.GetString("JIRA_ENDPOINT"),
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", config.V.GetString("JIRA_BASIC_AUTH_ENCODED")),
		},
		10*time.Second,
		3,
	)

	res, err := jiraApiClient.Get(fmt.Sprintf("/agile/1.0/board/%v/issue", boardId), nil, nil)
	if err != nil {
		return err
	}

	// TODO: this should not be set to 999
	jiraApiResponse := &JiraApiResponse{}

	logger.Info("res", res.Body)

	err = core.UnmarshalResponse(res, jiraApiResponse)

	if err != nil {
		logger.Error("Error: ", err)
		return nil
	}
	logger.Info("jiraIssues ", jiraApiResponse)

	// TODO: save more than one
	jiraIssue := &models.JiraIssue{
		Key: jiraApiResponse.Issues[0].Key,
	}
	err = lakeModels.Db.Save(jiraIssue).Error
	if err != nil {
		logger.Error("Error: ", err)
	}
	return nil
}
