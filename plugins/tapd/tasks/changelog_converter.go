package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"reflect"
)

func ConvertChangelog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert changelog :%d", data.Options.WorkspaceId)
	clIdGen := didgen.NewDomainIdGenerator(&models.TapdChangelogItem{})

	cursor, err := db.Table("_tool_tapd_changelog_items").
		Joins("left join _tool_tapd_changelogs tc on tc.id = _tool_tapd_changelog_items.changelog_id ").
		Where("tc.source_id = ? AND tc.workspace_id = ?", data.Source.ID, data.Options.WorkspaceId).
		Select("tc.issue_id as issue_id, " +
			"tc.creator as author_name," +
			"tc.created as created_date," +
			"tc.id as id," +
			"_tool_tapd_changelog_items.iteration_id_from," +
			"_tool_tapd_changelog_items.iteration_id_to," +
			"_tool_tapd_changelog_items.field as field_id, " +
			"_tool_tapd_changelog_items.field as field_name," +
			"_tool_tapd_changelog_items.value_before_parsed as 'from'," +
			"_tool_tapd_changelog_items.value_after_parsed as 'to'," +
			"_tool_tapd_changelog_items._raw_data_params as _raw_data_params," +
			"_tool_tapd_changelog_items._raw_data_table as _raw_data_table," +
			"_tool_tapd_changelog_items._raw_data_id as _raw_data_id," +
			"_tool_tapd_changelog_items._raw_data_remark as _raw_data_remark").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	changelogToHistoryConverter := NewChangelogToHistoryConverter(taskCtx)
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId:   data.Source.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: "_tool_tapd_api_%_changelogs",
		},
		InputRowType: reflect.TypeOf(models.ChangelogTmp{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			cl := inputRow.(*models.ChangelogTmp)
			domainCl := &ticket.Changelog{
				DomainEntity: domainlayer.DomainEntity{
					Id: clIdGen.Generate(data.Source.ID, cl.Id, cl.FieldId),
				},
				IssueId:     IssueIdGen.Generate(data.Source.ID, cl.IssueId),
				AuthorId:    UserIdGen.Generate(data.Source.ID, data.Options.WorkspaceId, cl.AuthorName),
				AuthorName:  cl.AuthorName,
				FieldId:     cl.FieldId,
				FieldName:   cl.FieldName,
				From:        cl.From,
				To:          cl.To,
				CreatedDate: cl.CreatedDate,
			}
			changelogToHistoryConverter.FeedIn(data.Source.ID, *cl)

			return []interface{}{
				domainCl,
			}, nil
		},
	})
	if err != nil {
		logger.Info(err.Error())
		return err
	}

	err = converter.Execute()
	if err != nil {
		return err
	}
	return changelogToHistoryConverter.UpdateSprintIssue()
}

var ConvertChangelogMeta = core.SubTaskMeta{
	Name:             "convertChangelog",
	EntryPoint:       ConvertChangelog,
	EnabledByDefault: true,
	Description:      "convert Tapd changelog",
}
