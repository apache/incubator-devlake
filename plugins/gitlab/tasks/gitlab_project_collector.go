package tasks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

type GitlabApiProject struct {
	GitlabId          int    `json:"id"`
	Name              string `josn:"name"`
	Description       string `json:"description"`
	DefaultBranch     string `json:"default_branch"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebUrl            string `json:"web_url"`
	CreatorId         int
	Visibility        string            `json:"visibility"`
	OpenIssuesCount   int               `json:"open_issues_count"`
	StarCount         int               `json:"star_count"`
	ForkedFromProject *GitlabApiProject `json:"forked_from_project"`
	CreatedAt         core.Iso8601Time  `json:"created_at"`
	LastActivityAt    *core.Iso8601Time `json:"last_activity_at"`
}

type GitlabApiProjectResponse GitlabApiProject

func CollectProject(ctx context.Context, projectId int, gitlabApiClient *GitlabApiClient) error {

	res, err := gitlabApiClient.Get(fmt.Sprintf("projects/%v", projectId), nil, nil)
	if err != nil {
		logger.Error("Error: ", err)
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	gitlabApiResponse := &GitlabApiProjectResponse{}
	err = core.UnmarshalResponse(res, gitlabApiResponse)
	if err != nil {
		logger.Error("Error: ", err)
		return err
	}
	gitlabProject := &models.GitlabProject{
		GitlabId:          gitlabApiResponse.GitlabId,
		Name:              gitlabApiResponse.Name,
		Description:       gitlabApiResponse.Description,
		DefaultBranch:     gitlabApiResponse.DefaultBranch,
		CreatorId:         gitlabApiResponse.CreatorId,
		PathWithNamespace: gitlabApiResponse.PathWithNamespace,
		WebUrl:            gitlabApiResponse.WebUrl,
		Visibility:        gitlabApiResponse.Visibility,
		OpenIssuesCount:   gitlabApiResponse.OpenIssuesCount,
		StarCount:         gitlabApiResponse.StarCount,
		CreatedDate:       gitlabApiResponse.CreatedAt.ToTime(),
		UpdatedDate:       core.Iso8601TimeToTime(gitlabApiResponse.LastActivityAt),
	}
	if gitlabApiResponse.ForkedFromProject != nil {
		gitlabProject.ForkedFromProjectId = gitlabApiResponse.ForkedFromProject.GitlabId
		gitlabProject.ForkedFromProjectWebUrl = gitlabApiResponse.ForkedFromProject.WebUrl
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&gitlabProject).Error
	if err != nil {
		logger.Error("Could not upsert: ", err)
		return err
	}
	return nil
}

// Convert the API response to our DB model instance
func convertProject(gitlabApiProject *GitlabApiProject) *models.GitlabProject {
	gitlabProject := &models.GitlabProject{
		GitlabId:          gitlabApiProject.GitlabId,
		Name:              gitlabApiProject.Name,
		Description:       gitlabApiProject.Description,
		DefaultBranch:     gitlabApiProject.DefaultBranch,
		CreatorId:         gitlabApiProject.CreatorId,
		PathWithNamespace: gitlabApiProject.PathWithNamespace,
		WebUrl:            gitlabApiProject.WebUrl,
		Visibility:        gitlabApiProject.Visibility,
		OpenIssuesCount:   gitlabApiProject.OpenIssuesCount,
		StarCount:         gitlabApiProject.StarCount,
		CreatedDate:       gitlabApiProject.CreatedAt.ToTime(),
		UpdatedDate:       core.Iso8601TimeToTime(gitlabApiProject.LastActivityAt),
	}
	if gitlabApiProject.ForkedFromProject != nil {
		gitlabProject.ForkedFromProjectId = gitlabApiProject.ForkedFromProject.GitlabId
		gitlabProject.ForkedFromProjectWebUrl = gitlabApiProject.ForkedFromProject.WebUrl
	}
	return gitlabProject
}
