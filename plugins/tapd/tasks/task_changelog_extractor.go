package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"strings"
)

var _ core.SubTaskEntryPoint = ExtractTaskChangelog

var ExtractTaskChangelogMeta = core.SubTaskMeta{
	Name:             "extractTaskChangelog",
	EntryPoint:       ExtractTaskChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdTaskChangelogRes struct {
	WorkitemChange models.TapdTaskChangelog
}

func ExtractTaskChangelog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_TASK_CHANGELOG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {

			results := make([]interface{}, 0, 2)

			var taskChangelogBody TapdTaskChangelogRes
			err := json.Unmarshal(row.Data, &taskChangelogBody)
			if err != nil {
				return nil, err
			}
			taskChangelog := taskChangelogBody.WorkitemChange

			taskChangelog.SourceId = data.Source.ID
			for _, fc := range taskChangelog.FieldChanges {
				var item *models.TapdTaskChangelogItem
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
						item.ChangelogId = taskChangelog.ID
						item.Field = k
						item.ValueAfterParsed = v.(string)
						item.ValueBeforeParsed = valueBeforeMap[k]
						results = append(results, item)
					}
				default:
					item.SourceId = data.Source.ID
					item.ChangelogId = taskChangelog.ID
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
				results = append(results, item)
			}
			results = append(results, &taskChangelog)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
