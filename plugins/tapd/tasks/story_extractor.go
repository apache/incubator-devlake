package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"strconv"
	"strings"
	"time"
)

var _ core.SubTaskEntryPoint = ExtractStories

var ExtractStoryMeta = core.SubTaskMeta{
	Name:             "extractStories",
	EntryPoint:       ExtractStories,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table tapd_iterations",
}

type TapdStoryRes struct {
	Story models.TapdStoryApiRes
}

func ExtractStories(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	db := taskCtx.GetDb()
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
			storyRes := storyBody.Story

			i, err := VoToDTO(&storyRes, &models.TapdStory{})
			if err != nil {
				return nil, err
			}
			toolL := i.(*models.TapdStory)
			toolL.SourceId = data.Source.ID
			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/prong/stories/view/%d", toolL.WorkspaceId, toolL.ID)
			if strings.Contains(toolL.Owner, ";") {
				toolL.Owner = strings.Split(toolL.Owner, ";")[0]
			}
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
			changelogs := make([]*models.ChangelogTmp, 0)
			err = db.Table("tapd_changelog_items").
				Joins("left join tapd_changelogs tc on tc.id = tapd_changelog_items.changelog_id ").
				Where("tc.source_id = ? AND tc.workspace_id = ? and tc.issue_id = ?",
					data.Source.ID, data.Options.WorkspaceId, toolL.ID).
				Order("tc.created desc").
				Pluck("tc.issue_id as issue_id, "+
					"tc.creator as author_name,"+
					"tc.created as created_date,"+
					"tc.id as id,"+
					"tapd_changelog_items.field as field_id, "+
					"tapd_changelog_items.field as field_name,"+
					"tapd_changelog_items.value_before_parsed as 'from',"+
					"tapd_changelog_items.value_after_parsed as 'to',"+
					"tapd_changelog_items._raw_data_params as _raw_data_params,"+
					"tapd_changelog_items._raw_data_table as _raw_data_table,"+
					"tapd_changelog_items._raw_data_id as _raw_data_id,"+
					"tapd_changelog_items._raw_data_remark as _raw_data_remark", &changelogs).Error
			if err != nil {
				return nil, err
			}
			// These three will be used as endDate
			lastSprintCreateDate := time.Now()
			lastStatusCreateDate := time.Now()
			lastAssignCreateDate := time.Now()
			for _, v := range changelogs {
				if v.FieldName == "iteration_id" {
					iteration := &models.TapdIteration{}
					err = db.Model(&models.TapdIteration{}).
						Where("source_id = ? and workspace_id = ? and name = ?",
							data.Source.ID, data.Options.WorkspaceId, v.To).Limit(1).Find(iteration).Error
					if err != nil {
						return nil, err
					}
					tapdIssueSprint := &models.TapdIssueSprintsHistory{
						SourceId:    data.Source.ID,
						WorkspaceId: data.Options.WorkspaceId,
						IssueId:     toolL.ID,
						SprintId:    iteration.ID,
						StartDate:   v.CreatedDate,
						EndDate:     lastSprintCreateDate,
					}
					results = append(results, tapdIssueSprint)
					lastSprintCreateDate = v.CreatedDate
				}
				if v.FieldName == "status" {
					tapdIssueStatus := &models.TapdIssueStatusHistory{
						SourceId:       data.Source.ID,
						WorkspaceId:    data.Options.WorkspaceId,
						IssueId:        toolL.ID,
						OriginalStatus: v.To,
						StartDate:      v.CreatedDate,
						EndDate:        lastStatusCreateDate,
					}
					lastSprintCreateDate = v.CreatedDate
					results = append(results, tapdIssueStatus)
				}
				if v.FieldName == "owner" {
					tapdIssueAssign := &models.TapdIssueAssigneeHistory{
						SourceId:    data.Source.ID,
						WorkspaceId: data.Options.WorkspaceId,
						IssueId:     toolL.ID,
						Assignee:    v.To,
						StartDate:   v.CreatedDate,
						EndDate:     lastAssignCreateDate,
					}
					lastAssignCreateDate = v.CreatedDate
					results = append(results, tapdIssueAssign)
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
