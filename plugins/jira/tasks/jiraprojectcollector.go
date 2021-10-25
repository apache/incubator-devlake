package tasks

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

// This has to be called JiraUserProjects since it is a different api call than JiraProjects.
// This call retreives all projects that the user has access to
type JiraUserProjects struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

func CollectProjects(
	jiraApiClient *JiraApiClient,
	source *models.JiraSource,
	boardId uint64,
) error {
	logger.Info("JON >>> attempting to get projects", true)
	res, err := jiraApiClient.Get("/rest/api/2/project", nil, nil)
	if err != nil {
		return err
	}

	// parse response
	var response []JiraUserProjects
	err = core.UnmarshalResponse(res, response)
	if err != nil {
		return err
	}

	// process issues
	for _, project := range response {
		err = lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(project).Error
		if err != nil {
			return err
		}
	}
	return nil
}
