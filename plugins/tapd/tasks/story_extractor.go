package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"strconv"
	"strings"
)

var _ core.SubTaskEntryPoint = ExtractStories

var ExtractStoryMeta = core.SubTaskMeta{
	Name:             "extractStories",
	EntryPoint:       ExtractStories,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdStoryRes struct {
	Story models.TapdStory
}

func ExtractStories(taskCtx core.SubTaskContext) error {
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
		if statusKey == "resolved" {
			return ticket.DONE
		} else if statusKey == "planning" {
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
			Table: RAW_STORY_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var storyBody TapdStoryRes
			err := json.Unmarshal(row.Data, &storyBody)
			if err != nil {
				return nil, err
			}
			toolL := storyBody.Story

			toolL.SourceId = models.Uint64s(data.Source.ID)
			toolL.StdType = "REQUIREMENT"
			toolL.StdStatus = getStdStatus(toolL.Status)
			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/prong/stories/view/%d", toolL.WorkspaceId, toolL.ID)
			if strings.Contains(toolL.Owner, ";") {
				toolL.Owner = strings.Split(toolL.Owner, ";")[0]
			}
			workSpaceIssue := &models.TapdWorkSpaceIssue{
				SourceId:    models.Uint64s(data.Source.ID),
				WorkspaceId: toolL.WorkspaceId,
				IssueId:     toolL.ID,
			}
			results := make([]interface{}, 0, 3)
			results = append(results, &toolL, workSpaceIssue)
			if toolL.IterationID != 0 {
				iterationIssue := &models.TapdIterationIssue{
					SourceId:         models.Uint64s(data.Source.ID),
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
					toolLIssueLabel := &models.TapdStoryLabel{
						StoryId:   toolL.ID,
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

func AtoIIgnoreEmpty(s string) (int, error) {
	if len(s) == 0 {
		return 0, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}
