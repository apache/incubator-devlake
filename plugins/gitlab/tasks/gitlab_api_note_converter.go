package tasks

import (
	"reflect"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

func ConvertApiNotes(taskCtx core.SubTaskContext) error {

	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)

	cursor, err := lakeModels.Db.Model(&models.GitlabMergeRequestNote{}).
		Joins("left join gitlab_merge_requests on gitlab_merge_requests.gitlab_id = gitlab_merge_request_notes.merge_request_id").
		Where("gitlab_merge_requests.project_id = ?", data.Options.ProjectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	domainIdGeneratorNote := didgen.NewDomainIdGenerator(&models.GitlabMergeRequestNote{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabMergeRequestNote{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabNotes := inputRow.(*models.GitlabMergeRequestNote)
			domainNote := convertToNoteModel(gitlabNotes, domainIdGeneratorNote)

			return []interface{}{
				domainNote,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
