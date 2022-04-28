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

func ConvertWorklog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert board:%d", data.Options.WorkspaceID)
	worklogIdGen := didgen.NewDomainIdGenerator(&models.TapdWorklog{})
	cursor, err := db.Model(&models.TapdWorklog{}).Where("source_id = ? AND workspace_id = ?", data.Source.ID, data.Options.WorkspaceID).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId:   data.Source.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_WORKLOG_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdWorklog{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			toolL := inputRow.(*models.TapdWorklog)
			domainL := &ticket.IssueWorklog{
				DomainEntity: domainlayer.DomainEntity{
					Id: worklogIdGen.Generate(models.Uint64s(data.Source.ID), toolL.ID),
				},
				AuthorId:         UserIdGen.Generate(models.Uint64s(data.Source.ID), toolL.WorkspaceID, toolL.Owner),
				Comment:          toolL.Memo,
				TimeSpentMinutes: int(toolL.Timespent),
				LoggedDate:       toolL.Created.ToNullableTime(),
				//IssueId:          toolL.EntityID,
			}
			switch toolL.EntityType {
			case "TASK":
				domainL.IssueId = didgen.
					NewDomainIdGenerator(&models.TapdTask{}).Generate(toolL.EntityID)
				break
			case "BUG":
				domainL.IssueId = didgen.
					NewDomainIdGenerator(&models.TapdBug{}).Generate(toolL.EntityID)
				break
			case "STORY":
				domainL.IssueId = didgen.
					NewDomainIdGenerator(&models.TapdStory{}).Generate(toolL.EntityID)
				break
			}
			return []interface{}{
				domainL,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertWorklogMeta = core.SubTaskMeta{
	Name:             "convertWorklog",
	EntryPoint:       ConvertWorklog,
	EnabledByDefault: true,
	Description:      "convert Tapd Worklog",
}
