package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"strings"
	"time"
)

var _ core.SubTaskEntryPoint = ExtractBugs

var ExtractBugMeta = core.SubTaskMeta{
	Name:             "extractBugs",
	EntryPoint:       ExtractBugs,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table tapd_iterations",
}

type TapdBugRes struct {
	Bug models.TapdBugApiRes
}

func ExtractBugs(taskCtx core.SubTaskContext) error {
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
			Table: RAW_BUG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var bugBody TapdBugRes
			err := json.Unmarshal(row.Data, &bugBody)
			if err != nil {
				return nil, err
			}
			bugRes := bugBody.Bug

			i, err := VoToDTO(&bugRes, &models.TapdBug{})
			if err != nil {
				return nil, err
			}
			toolL := i.(*models.TapdBug)
			toolL.SourceId = data.Source.ID
			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/prong/stories/view/%d", toolL.WorkspaceId, toolL.ID)
			if strings.Contains(toolL.CurrentOwner, ";") {
				toolL.CurrentOwner = strings.Split(toolL.CurrentOwner, ";")[0]
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
					ResolutionDate:   toolL.Resolved,
					IssueCreatedDate: toolL.Created,
				}
				results = append(results, iterationIssue)
			}
			changelogs := make([]*models.ChangelogTmp, 0)
			err = db.Table("tapd_changelog_items").
				Joins("left join tapd_changelogs tc on tc.id = tapd_changelog_items.changelog_id ").
				Where("tc.source_id = ? AND tc.workspace_id = ? AND tc.issue_id = ?",
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
				if v.FieldName == "current_owner" {
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
