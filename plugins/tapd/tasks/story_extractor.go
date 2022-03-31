package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"strconv"
)

var _ core.SubTaskEntryPoint = ExtractStories

var ExtractStoriesMeta = core.SubTaskMeta{
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
			var iterRes TapdStoryRes
			err := json.Unmarshal(row.Data, &iterRes)
			if err != nil {
				return nil, err
			}
			storyRes := iterRes.Story

			idInt, err := AtoIIgnoreEmpty(storyRes.ID)
			if err != nil {
				return nil, err
			}
			workitemTypeID, err := AtoIIgnoreEmpty(storyRes.WorkitemTypeID)
			if err != nil {
				return nil, err
			}
			iterationID, err := AtoIIgnoreEmpty(storyRes.IterationID)
			if err != nil {
				return nil, err
			}
			categoryID, err := AtoIIgnoreEmpty(storyRes.CategoryID)
			if err != nil {
				return nil, err
			}
			parentID, err := AtoIIgnoreEmpty(storyRes.ParentID)
			if err != nil {
				return nil, err
			}
			ancestorID, err := AtoIIgnoreEmpty(storyRes.AncestorID)
			if err != nil {
				return nil, err
			}
			effort, err := AtoIIgnoreEmpty(storyRes.Effort)
			if err != nil {
				return nil, err
			}
			effortCompleted, err := AtoIIgnoreEmpty(storyRes.EffortCompleted)
			if err != nil {
				return nil, err
			}
			exceed, err := AtoIIgnoreEmpty(storyRes.Exceed)
			if err != nil {
				return nil, err
			}
			remain, err := AtoIIgnoreEmpty(storyRes.Remain)
			if err != nil {
				return nil, err
			}
			releaseID, err := AtoIIgnoreEmpty(storyRes.ReleaseID)
			if err != nil {
				return nil, err
			}
			templatedID, err := AtoIIgnoreEmpty(storyRes.TemplatedID)
			if err != nil {
				return nil, err
			}
			begin, err := core.ConvertStringToTimePtr(storyRes.Begin)
			if err != nil {
				return nil, err
			}
			due, err := core.ConvertStringToTimePtr(storyRes.Due)
			if err != nil {
				return nil, err
			}
			created, err := core.ConvertStringToTimePtr(storyRes.Created)
			if err != nil {
				return nil, err
			}
			modified, err := core.ConvertStringToTimePtr(storyRes.Modified)
			if err != nil {
				return nil, err
			}

			issue := &models.TapdStory{
				SourceId:        data.Source.ID,
				ID:              uint64(idInt),
				WorkitemTypeID:  uint64(workitemTypeID),
				Name:            storyRes.Name,
				Description:     storyRes.Description,
				WorkspaceID:     data.Options.WorkspaceId,
				Creator:         storyRes.Creator,
				Created:         created,
				Modified:        modified,
				Status:          storyRes.Status,
				Owner:           storyRes.Owner,
				Cc:              storyRes.Cc,
				Begin:           begin,
				Due:             due,
				Priority:        storyRes.Priority,
				Developer:       storyRes.Developer,
				IterationID:     uint64(iterationID),
				TestFocus:       storyRes.TestFocus,
				Type:            storyRes.Type,
				Source:          storyRes.Source,
				Module:          storyRes.Module,
				Version:         storyRes.Version,
				CategoryID:      uint64(categoryID),
				Path:            storyRes.Path,
				ParentID:        uint64(parentID),
				ChildrenID:      storyRes.ChildrenID,
				AncestorID:      uint64(ancestorID),
				BusinessValue:   storyRes.BusinessValue,
				Effort:          uint64(effort),
				EffortCompleted: uint64(effortCompleted),
				Exceed:          uint64(exceed),
				Remain:          uint64(remain),
				ReleaseID:       uint64(releaseID),
				Confidential:    storyRes.Confidential,
				TemplatedID:     uint64(templatedID),
				CreatedFrom:     storyRes.CreatedFrom,
				Feature:         storyRes.Feature,
				NoPKModel:       common.NoPKModel{},
			}
			if storyRes.Size != "" {
				v, err := AtoIIgnoreEmpty(storyRes.Size)
				if err != nil {
					return nil, err
				}
				issue.Size = v
			}
			if storyRes.Effort != "" {
				v, err := AtoIIgnoreEmpty(storyRes.Effort)
				if err != nil {
					return nil, err
				}
				issue.Effort = uint64(v)
			}
			if storyRes.Completed != "" {
				v, err := core.ConvertStringToTimePtr(storyRes.Completed)
				if err != nil {
					return nil, err
				}
				issue.Completed = v
			}
			if storyRes.EffortCompleted != "" {
				v, err := AtoIIgnoreEmpty(storyRes.EffortCompleted)
				if err != nil {
					return nil, err
				}
				issue.EffortCompleted = uint64(v)
			}
			if storyRes.EffortCompleted != "" {
				v, err := AtoIIgnoreEmpty(storyRes.EffortCompleted)
				if err != nil {
					return nil, err
				}
				issue.EffortCompleted = uint64(v)
			}

			return []interface{}{
				issue,
			}, nil
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
