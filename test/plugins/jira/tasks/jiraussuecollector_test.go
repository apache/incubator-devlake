package plugins

import (
	"testing"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks"
)

func TestPluginsLoading(t *testing.T) {
	err := lakeModels.Db.AutoMigrate(&models.JiraIssue{}, &models.JiraBoard{})
	if err != nil {
		t.Errorf("Failed to Migrate %v", err)
	}
	collectErr := tasks.CollectIssues(8)
	if collectErr != nil {
		t.Errorf("Failed to Collect Issue %v", collectErr)
	}
}
