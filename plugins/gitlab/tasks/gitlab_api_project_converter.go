package tasks

import (
	"reflect"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

func ConvertApiProjects(taskCtx core.SubTaskContext) error {

	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)

	//Find all piplines associated with the current projectid
	cursor, err := lakeModels.Db.Model(&models.GitlabProject{}).Where("gitlab_id=?", data.Options.ProjectId).Rows()
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
