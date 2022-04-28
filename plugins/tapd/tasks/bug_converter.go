package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"reflect"
	"strconv"
)

func ConvertBug(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert board:%d", data.Options.WorkspaceID)

	cursor, err := db.Model(&models.TapdBug{}).Where("source_id = ? AND workspace_id = ?", data.Source.ID, data.Options.WorkspaceID).Rows()
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
			Table: RAW_BUG_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdBug{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			toolL := inputRow.(*models.TapdBug)
			domainL := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: IssueIdGen.Generate(toolL.SourceId, toolL.ID),
				},
				Url:            toolL.Url,
				Number:         strconv.FormatUint(uint64(toolL.ID), 10),
				Title:          toolL.Title,
				EpicKey:        toolL.EpicKey,
				Type:           "BUG",
				Status:         toolL.StdStatus,
				ResolutionDate: toolL.Resolved.ToNullableTime(),
				CreatedDate:    toolL.Created.ToNullableTime(),
				UpdatedDate:    toolL.Modified.ToNullableTime(),
				ParentIssueId:  IssueIdGen.Generate(toolL.SourceId, toolL.IssueID),
				Priority:       toolL.Priority,
				CreatorId:      UserIdGen.Generate(data.Options.SourceId, toolL.WorkspaceID, toolL.Reporter),
				AssigneeId:     UserIdGen.Generate(data.Options.SourceId, toolL.WorkspaceID, toolL.CurrentOwner),
				AssigneeName:   toolL.CurrentOwner,
				Severity:       toolL.Severity,
				Component:      toolL.Feature, // todo not sure about this
				OriginalStatus: toolL.Status,
			}
			if domainL.ResolutionDate != nil && domainL.CreatedDate != nil {
				domainL.LeadTimeMinutes = uint(int64(domainL.ResolutionDate.Minute() - domainL.CreatedDate.Minute()))
			}
			results := make([]interface{}, 0, 2)
			boardIssue := &ticket.BoardIssue{
				BoardId: WorkspaceIdGen.Generate(data.Options.WorkspaceID),
				IssueId: domainL.Id,
			}
			results = append(results, domainL, boardIssue)
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertBugMeta = core.SubTaskMeta{
	Name:             "convertBug",
	EntryPoint:       ConvertBug,
	EnabledByDefault: true,
	Description:      "convert Tapd Bug",
}
