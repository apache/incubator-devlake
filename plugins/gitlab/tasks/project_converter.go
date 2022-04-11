package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertProjectMeta = core.SubTaskMeta{
	Name:             "convertApiProject",
	EntryPoint:       ConvertApiProjects,
	EnabledByDefault: true,
	Description:      "Update domain layer Repo according to GitlabProject",
}

func ConvertApiProjects(taskCtx core.SubTaskContext) error {

	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)
	db := taskCtx.GetDb()

	//Find all piplines associated with the current projectid
	cursor, err := db.Model(&models.GitlabProject{}).Where("gitlab_id=?", data.Options.ProjectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabProject{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabProject := inputRow.(*models.GitlabProject)

			domainRepository := convertToRepositoryModel(gitlabProject)

			return []interface{}{
				domainRepository,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
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

func convertToRepositoryModel(project *models.GitlabProject) *code.Repo {
	domainRepository := &code.Repo{
		DomainEntity: domainlayer.DomainEntity{
			Id: didgen.NewDomainIdGenerator(project).Generate(project.GitlabId),
		},
		Name:        project.Name,
		Url:         project.WebUrl,
		Description: project.Description,
		ForkedFrom:  project.ForkedFromProjectWebUrl,
		CreatedDate: project.CreatedDate,
		UpdatedDate: project.UpdatedDate,
	}
	return domainRepository
}
