package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

// This has to be called JiraUserProjects since it is a different api call than JiraProjects.
// This call retreives all projects that the user has access to
type JiraUserProjectApiRes []JiraApiProject
type JiraApiProject struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

func CollectProjects(
	jiraApiClient *JiraApiClient,
	sourceId uint64,
) error {
	res, err := jiraApiClient.Get("rest/api/3/project", nil, nil)
	if err != nil {
		return err
	}

	jiraApiProjects := &JiraUserProjectApiRes{}

	err = core.UnmarshalResponse(res, jiraApiProjects)
	if err != nil {
		return err
	}
	// process issues
	for _, jiraApiProject := range *jiraApiProjects {
		jiraProject, _ := convertProject(&jiraApiProject, sourceId)
		err = lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(jiraProject).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func convertProject(jiraApiProject *JiraApiProject, sourceId uint64) (*models.JiraProject, error) {
	jiraProject := &models.JiraProject{
		SourceId: sourceId,
		Id:       jiraApiProject.Id,
		Key:      jiraApiProject.Key,
		Name:     jiraApiProject.Name,
	}
	return jiraProject, nil
}
