package tasks

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
)

type JiraApiLocation struct {
	ProjectId      int    `json:"projectId"`
	DisplayName    string `json:"displayName"`
	ProjectName    string `json:"projectName"`
	ProjectKey     string `json:"projectKey"`
	ProjectTypeKey string `json:"projectTypeKey"`
	AvatarURI      string `json:"avatarURI"`
	Name           string `json:"name"`
}

type JiraApiBoard struct {
	Id       int    `json:"id"`
	Self     string `json:"self"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location *JiraApiLocation
}

func CollectBoard(boardId int) error {
	jiraApiClient := core.NewApiClient(
		config.V.GetString("JIRA_ENDPOINT"),
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", config.V.GetString("JIRA_BASIC_AUTH_ENCODED")),
		},
		10*time.Second,
		3,
	)

	res, err := jiraApiClient.Get(fmt.Sprintf("/agile/1.0/board/%v", boardId), nil, nil)
	if err != nil {
		return err
	}

	jiraApiBoard := &JiraApiBoard{}
	err = core.UnmarshalResponse(res, jiraApiBoard)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil
	}
	fmt.Printf("jiraboard %v", jiraApiBoard)
	jiraBoard := &models.JiraBoard{
		JiraId: jiraApiBoard.Id,
		Name:   jiraApiBoard.Name,
	}
	err = lakeModels.Db.Save(jiraBoard).Error
	if err != nil {
		fmt.Println("Error: %w", err)
	}
	return nil
}
