package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"strings"
)

var _ core.SubTaskEntryPoint = ExtractBugs

var ExtractBugMeta = core.SubTaskMeta{
	Name:             "extractBugs",
	EntryPoint:       ExtractBugs,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdBugRes struct {
	Bug models.TapdBug
}

func ExtractBugs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)

	getStdStatus := func(statusKey string) string {
		if statusKey == "done" {
			return ticket.DONE
		} else if statusKey == "new" {
			return ticket.TODO
		} else {
			return ticket.IN_PROGRESS
		}
	}
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_BUG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var bugBody TapdBugRes
			err := json.Unmarshal(row.Data, &bugBody)
			if err != nil {
				return nil, err
			}
			toolL := bugBody.Bug

			toolL.SourceId = models.Uint64s(data.Source.ID)
			toolL.Type = "BUG"
			toolL.StdType = "BUG"
			toolL.StdStatus = getStdStatus(toolL.Status)
			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/prong/stories/view/%d", toolL.WorkspaceID, toolL.ID)
			if strings.Contains(toolL.CurrentOwner, ";") {
				toolL.CurrentOwner = strings.Split(toolL.CurrentOwner, ";")[0]
			}
			workSpaceIssue := &models.TapdWorkSpaceIssue{
				SourceId:    models.Uint64s(data.Source.ID),
				WorkspaceID: toolL.WorkspaceID,
				IssueId:     toolL.ID,
			}
			results := make([]interface{}, 0, 3)
			results = append(results, &toolL, workSpaceIssue)
			if toolL.IterationID != 0 {
				iterationIssue := &models.TapdIterationIssue{
					SourceId:         models.Uint64s(data.Source.ID),
					IterationId:      toolL.IterationID,
					IssueId:          toolL.ID,
					ResolutionDate:   toolL.Resolved,
					IssueCreatedDate: toolL.Created,
				}
				results = append(results, iterationIssue)
			}
			if toolL.Label != "" {
				labelList := strings.Split(toolL.Label, "|")
				for _, v := range labelList {
					toolLIssueLabel := &models.TapdBugLabel{
						BugId:     toolL.ID,
						LabelName: v,
					}
					results = append(results, toolLIssueLabel)
				}
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
