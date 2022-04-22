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

var _ core.SubTaskEntryPoint = ExtractTasks

var ExtractTaskMeta = core.SubTaskMeta{
	Name:             "extractTasks",
	EntryPoint:       ExtractTasks,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdTaskRes struct {
	Task models.TapdTaskApiRes
}

func ExtractTasks(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	db := taskCtx.GetDb()
	sourceId := data.Source.ID
	// prepare getStdStatus function
	var statusMappingRows []*models.TapdIssueStatusMapping
	err := db.Find(&statusMappingRows, "source_id = ?", sourceId).Error
	if err != nil {
		return err
	}
	statusMappings := make(map[string]string)
	makeStatusMappingKey := func(userType string, userStatus string) string {
		return fmt.Sprintf("%v:%v", userType, userStatus)
	}
	for _, statusMappingRow := range statusMappingRows {
		k := makeStatusMappingKey(statusMappingRow.UserType, statusMappingRow.UserStatus)
		statusMappings[k] = statusMappingRow.StandardStatus
	}
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
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_TASK_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var taskBody TapdTaskRes
			err := json.Unmarshal(row.Data, &taskBody)
			if err != nil {
				return nil, err
			}
			taskRes := taskBody.Task

			i, err := VoToDTO(&taskRes, &models.TapdTask{})
			if err != nil {
				return nil, err
			}
			toolL := i.(*models.TapdTask)
			toolL.SourceId = data.Source.ID
			toolL.Type = "TASK"
			toolL.StdType = "TASK"
			toolL.StdStatus = getStdStatus(toolL.Status)
			if strings.Contains(toolL.Owner, ";") {
				toolL.Owner = strings.Split(toolL.Owner, ";")[0]
			}
			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/prong/stories/view/%d", toolL.WorkspaceId, toolL.ID)

			workSpaceIssue := &models.TapdWorkSpaceIssue{
				SourceId:    data.Source.ID,
				WorkspaceId: toolL.WorkspaceId,
				IssueId:     toolL.ID,
			}
			results := make([]interface{}, 0, 3)
			results = append(results, toolL, workSpaceIssue)
			if toolL.IterationID != 0 {
				iterationIssue := &models.TapdIterationIssue{
					SourceId:         data.Source.ID,
					IterationId:      toolL.IterationID,
					IssueId:          toolL.ID,
					ResolutionDate:   toolL.Completed,
					IssueCreatedDate: toolL.Created,
				}
				results = append(results, iterationIssue)
			}
			if toolL.Label != "" {
				labelList := strings.Split(toolL.Label, "|")
				for _, v := range labelList {
					toolLIssueLabel := &models.TapdIssueLabel{
						IssueId:   toolL.ID,
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
