package tasks

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

type JiraApiLocation struct {
	ProjectId      uint   `json:"projectId"`
	DisplayName    string `json:"displayName"`
	ProjectName    string `json:"projectName"`
	ProjectKey     string `json:"projectKey"`
	ProjectTypeKey string `json:"projectTypeKey"`
	AvatarURI      string `json:"avatarURI"`
	Name           string `json:"name"`
}

type JiraApiBoard struct {
	Id       uint64 `json:"id"`
	Self     string `json:"self"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location *JiraApiLocation
}

func CollectBoard(jiraApiClient *JiraApiClient, source *models.JiraSource, boardId uint64) error {
	res, err := jiraApiClient.Get(fmt.Sprintf("rest/agile/1.0/board/%v", boardId), nil, nil)
	if err != nil {
		return err
	}

	jiraApiBoard := &JiraApiBoard{}
	err = core.UnmarshalResponse(res, jiraApiBoard)
	if err != nil {
		logger.Error("Error: ", err)
		return nil
	}
	logger.Info("jiraboard ", jiraApiBoard)
	jiraBoard := &models.JiraBoard{
		SourceId:  source.ID,
		BoardId:   jiraApiBoard.Id,
		ProjectId: jiraApiBoard.Location.ProjectId,
		Name:      jiraApiBoard.Name,
		Self:      jiraApiBoard.Self,
		Type:      jiraApiBoard.Type,
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(jiraBoard).Error
	if err != nil {
		logger.Error("Error: ", err)
	}
	return nil
}
