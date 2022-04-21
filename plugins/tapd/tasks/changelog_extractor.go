package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractChangelog

var ExtractChangelogMeta = core.SubTaskMeta{
	Name:             "extractStoryChangelog",
	EntryPoint:       ExtractChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdStoryChangelogRes struct {
	WorkitemChange models.TapdStoryChangelogApiRes
}

type TapdBugChangelogRes struct {
	BugChange models.TapdBugChangelogApiRes
}
type TapdTaskChangelogRes struct {
	WorkitemChange models.TapdTaskChangelogApiRes
}

func ExtractChangelog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_CHANGELOG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			changelogType := &Type{}
			err := json.Unmarshal(row.Input, changelogType)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 2)
			switch changelogType.Type {
			case "bug":
				var bugChangelogBody TapdBugChangelogRes
				err := json.Unmarshal(row.Data, &bugChangelogBody)
				if err != nil {
					return nil, err
				}
				bugChangelogRes := bugChangelogBody.BugChange

				i, err := VoToDTO(&bugChangelogRes, &models.TapdBugChangelog{})
				if err != nil {
					return nil, err
				}
				v := i.(*models.TapdBugChangelog)
				toolL := &models.TapdChangelog{
					SourceId:      data.Source.ID,
					ID:            v.ID,
					WorkspaceId:   v.WorkspaceId,
					Creator:       v.Author,
					Created:       v.Created,
					ChangeSummary: v.Memo,
					EntityType:    "BUG",
					ChangeType:    "BUG",
					IssueId:       v.BugID,
				}

				item := &models.TapdChangelogItem{
					SourceId:          data.Source.ID,
					ChangelogId:       v.ID,
					Field:             v.Field,
					ValueBeforeParsed: v.OldValue,
					ValueAfterParsed:  v.NewValue,
					NoPKModel:         common.NoPKModel{},
				}
				if item.Field == "iteration_id" {
					item, err = parseIterationChangelog(taskCtx, item)
					if err != nil {
						return nil, err
					}
				}
				results = append(results, toolL, item)
			case "story":
				var storyChangelogBody TapdStoryChangelogRes
				err = json.Unmarshal(row.Data, &storyChangelogBody)
				if err != nil {
					return nil, err
				}
				storyChangelogRes := storyChangelogBody.WorkitemChange

				i, err := VoToDTO(&storyChangelogRes, &models.TapdStoryChangelog{})
				if err != nil {
					return nil, err
				}
				v := i.(*models.TapdStoryChangelog)
				toolL := &models.TapdChangelog{
					SourceId:       data.Source.ID,
					ID:             v.ID,
					WorkspaceId:    v.WorkspaceId,
					WorkitemTypeID: v.WorkitemTypeID,
					Creator:        v.Creator,
					Created:        v.Created,
					ChangeSummary:  v.ChangeSummary,
					Comment:        v.Comment,
					EntityType:     v.EntityType,
					ChangeType:     v.EntityType,
					IssueId:        v.StoryID,
				}
				for _, fc := range v.FieldChanges {
					item := &models.TapdChangelogItem{
						SourceId:          data.Source.ID,
						ChangelogId:       toolL.ID,
						Field:             fc.Field,
						ValueBeforeParsed: fc.ValueBeforeParsed,
						ValueAfterParsed:  fc.ValueAfterParsed,
					}
					if item.Field == "iteration_id" {
						item, err = parseIterationChangelog(taskCtx, item)
						if err != nil {
							return nil, err
						}
					}
					results = append(results, item)
				}
				results = append(results, toolL)
			case "task":
				var taskChangelogBody TapdTaskChangelogRes
				err := json.Unmarshal(row.Data, &taskChangelogBody)
				if err != nil {
					return nil, err
				}
				taskChangelogRes := taskChangelogBody.WorkitemChange

				i, err := VoToDTO(&taskChangelogRes, &models.TapdTaskChangelog{})
				if err != nil {
					return nil, err
				}
				v := i.(*models.TapdTaskChangelog)
				toolL := &models.TapdChangelog{
					SourceId:       data.Source.ID,
					ID:             v.ID,
					WorkspaceId:    v.WorkspaceId,
					WorkitemTypeID: v.WorkitemTypeID,
					Creator:        v.Creator,
					Created:        v.Created,
					ChangeSummary:  v.ChangeSummary,
					Comment:        v.Comment,
					EntityType:     v.EntityType,
					ChangeType:     v.EntityType,
					IssueId:        v.TaskID,
				}
				for _, fc := range v.FieldChanges {
					item := &models.TapdChangelogItem{
						SourceId:          data.Source.ID,
						ChangelogId:       toolL.ID,
						Field:             fc.Field,
						ValueBeforeParsed: fc.ValueBeforeParsed,
						ValueAfterParsed:  fc.ValueAfterParsed,
					}
					if item.Field == "iteration_id" {
						item, err = parseIterationChangelog(taskCtx, item)
						if err != nil {
							return nil, err
						}
					}
					results = append(results, item)
				}
				results = append(results, toolL)
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
