package tasks

import (
	"fmt"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
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
	logger.Info("convert board:%d", data.Options.WorkspaceId)
	issueIdGen := didgen.NewDomainIdGenerator(&models.TapdBug{})
	cursor, err := db.Model(&models.TapdBug{}).Where("source_id = ? AND workspace_id = ?", data.Source.ID, data.Options.WorkspaceId).Rows()
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
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_WORKSPACE_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdBug{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			bug := inputRow.(*models.TapdBug)
			issue := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: issueIdGen.Generate(bug.SourceId, bug.ID),
				},
				Url:            fmt.Sprintf("https://www.tapd.cn/%d/prong/Bugs/view/%d", bug.WorkspaceID, bug.ID),
				Key:            strconv.FormatUint(bug.ID, 10),
				Title:          bug.Title,
				Summary:        bug.Title,
				EpicKey:        bug.EpicKey,
				Type:           "BUG",
				Status:         bug.Status,
				ResolutionDate: bug.Resolved,
				CreatedDate:    bug.Created,
				UpdatedDate:    bug.Modified,
				ParentIssueId:  issueIdGen.Generate(bug.SourceId, bug.IssueID),
				Priority:       bug.Priority,
				CreatorId:      UserIdGen.Generate(data.Options.SourceId, bug.Reporter),
				AssigneeId:     UserIdGen.Generate(data.Options.SourceId, bug.De),
				AssigneeName:   bug.De,
				Severity:       bug.Severity,
				Component:      bug.Feature, // todo not sure about this
			}
			if issue.ResolutionDate != nil && issue.CreatedDate != nil {
				issue.TimeSpentMinutes = int64(issue.ResolutionDate.Minute() - issue.CreatedDate.Minute())
			}
			return []interface{}{
				issue,
			}, nil
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
