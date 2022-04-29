package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"strings"
)

var _ core.SubTaskEntryPoint = ExtractStoryChangelog

var ExtractStoryChangelogMeta = core.SubTaskMeta{
	Name:             "extractStoryChangelog",
	EntryPoint:       ExtractStoryChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdStoryChangelogRes struct {
	WorkitemChange models.TapdStoryChangelog
}

func ExtractStoryChangelog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_STORY_CHANGELOG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {

			results := make([]interface{}, 0, 2)

			var storyChangelogBody TapdStoryChangelogRes
			err := json.Unmarshal(row.Data, &storyChangelogBody)
			if err != nil {
				return nil, err
			}
			storyChangelog := storyChangelogBody.WorkitemChange

			storyChangelog.SourceId = data.Source.ID
			for _, fc := range storyChangelog.FieldChanges {
				var item models.TapdStoryChangelogItem
				var valueAfterMap interface{}
				if err = json.Unmarshal(fc.ValueAfterParsed, &valueAfterMap); err != nil {
					return nil, err
				}
				switch valueAfterMap.(type) {
				case map[string]interface{}:
					valueBeforeMap := map[string]string{}
					err = json.Unmarshal(fc.ValueBeforeParsed, &valueBeforeMap)
					if err != nil {
						return nil, err
					}
					for k, v := range valueAfterMap.(map[string]interface{}) {
						item.SourceId = data.Source.ID
						item.ChangelogId = storyChangelog.ID
						item.Field = k
						item.ValueAfterParsed = v.(string)
						item.ValueBeforeParsed = valueBeforeMap[k]
					}
				default:
					item.SourceId = data.Source.ID
					item.ChangelogId = storyChangelog.ID
					item.Field = fc.Field
					item.ValueAfterParsed = strings.Trim(string(fc.ValueAfterParsed), `"`)
					item.ValueBeforeParsed = strings.Trim(string(fc.ValueBeforeParsed), `"`)
				}
				if item.Field == "iteration_id" {
					iterationFrom, iterationTo, err := parseIterationChangelog(taskCtx, item.ValueBeforeParsed, item.ValueAfterParsed)
					if err != nil {
						return nil, err
					}
					item.IterationIdFrom = iterationFrom
					item.IterationIdTo = iterationTo
				}
				results = append(results, &item)
			}
			results = append(results, &storyChangelog)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
