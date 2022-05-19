package tasks

import (
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"reflect"
	"time"
)

func ConvertWorklog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert board:%d", data.Options.WorkspaceID)
	worklogIdGen := didgen.NewDomainIdGenerator(&models.TapdWorklog{})
	cursor, err := db.Model(&models.TapdWorklog{}).Where("connection_id = ? AND workspace_id = ?", data.Connection.ID, data.Options.WorkspaceID).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,

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
					Id: worklogIdGen.Generate(data.Connection.ID, toolL.ID),
				},
				AuthorId:         UserIdGen.Generate(data.Connection.ID, toolL.WorkspaceID, toolL.Owner),
				Comment:          toolL.Memo,
				TimeSpentMinutes: int(toolL.Timespent),
				LoggedDate:       (*time.Time)(toolL.Created),
				//IssueId:          toolL.EntityID,
			}
			switch toolL.EntityType {
			case "TASK":
				domainL.IssueId = didgen.
					NewDomainIdGenerator(&models.TapdTask{}).Generate(toolL.EntityID)
			case "BUG":
				domainL.IssueId = didgen.
					NewDomainIdGenerator(&models.TapdBug{}).Generate(toolL.EntityID)
			case "STORY":
				domainL.IssueId = didgen.
					NewDomainIdGenerator(&models.TapdStory{}).Generate(toolL.EntityID)
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
