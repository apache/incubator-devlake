package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

// This has to be called JiraUserProjects since it is a different api call than JiraProjects.
// This call retreives all projects that the user has access to
type JiraUserProjectApiRes []JiraProject
type JiraProject struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

func CollectProjects(
	jiraApiClient *JiraApiClient,
	sourceId uint64,
) error {
	res, err := jiraApiClient.Get("/api/3/project", nil, nil)
	if err != nil {
		return err
	}

	jiraApiProjects := &JiraUserProjectApiRes{}

	err = core.UnmarshalResponse(res, jiraApiProjects)
	if err != nil {
		return err
	}
	// process issues
	for _, project := range *jiraApiProjects {
		project, _ := convertProject(&project, sourceId)
		err = lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(project).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func convertProject(user *JiraProject, sourceId uint64) (*models.JiraProject, error) {
	jiraProject := &models.JiraProject{
		SourceId: sourceId,
		Id:       user.Id,
		Key:      user.Key,
		Name:     user.Name,
	}
	return jiraProject, nil
}
