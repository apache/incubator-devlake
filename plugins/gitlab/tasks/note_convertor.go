package tasks

import (
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertApiNotesMeta = core.SubTaskMeta{
	Name:             "convertApiNotes",
	EntryPoint:       ConvertApiNotes,
	EnabledByDefault: true,
	Description:      "Update domain layer Note according to GitlabMergeRequestNote",
}

func ConvertApiNotes(taskCtx core.SubTaskContext) error {

	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)
	db := taskCtx.GetDb()

	cursor, err := db.Model(&models.GitlabMergeRequestNote{}).
		Joins("left join _tool_gitlab_merge_requests on _tool_gitlab_merge_requests.gitlab_id = _tool_gitlab_merge_request_notes.merge_request_id").
		Where("_tool_gitlab_merge_requests.project_id = ?", data.Options.ProjectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	domainIdGeneratorNote := didgen.NewDomainIdGenerator(&models.GitlabMergeRequestNote{})
	prIdGen := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})
	userIdGen := didgen.NewDomainIdGenerator(&models.GitlabUser{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabMergeRequestNote{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabNotes := inputRow.(*models.GitlabMergeRequestNote)
			domainNote := &code.Note{
				DomainEntity: domainlayer.DomainEntity{
					Id: domainIdGeneratorNote.Generate(gitlabNotes.GitlabId),
				},
				PrId:        prIdGen.Generate(gitlabNotes.MergeRequestId),
				Type:        gitlabNotes.NoteableType,
				Author:      userIdGen.Generate(gitlabNotes.AuthorUsername),
				Body:        gitlabNotes.Body,
				Resolvable:  gitlabNotes.Resolvable,
				IsSystem:    gitlabNotes.IsSystem,
				CreatedDate: gitlabNotes.GitlabCreatedAt,
			}
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
