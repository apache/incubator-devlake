package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"strconv"
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

			i, err := ResToDb(&storyRes, &models.TapdStory{})
			if err != nil {
				return nil, err
			}
			story := i.(*models.TapdStory)
			story.SourceId = data.Source.ID
			return []interface{}{
				story,
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
