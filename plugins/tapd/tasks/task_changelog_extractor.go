package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractTaskChangelog

var ExtractTaskChangelogMeta = core.SubTaskMeta{
	Name:             "extractTaskChangelog",
	EntryPoint:       ExtractTaskChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdTaskChangelogRes struct {
	WorkitemChange models.TapdTaskChangelogApiRes
}

func ExtractTaskChangelog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
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
			taskChangelogRes := taskChangelogBody.WorkitemChange

			i, err := VoToDTO(&taskChangelogRes, &models.TapdTaskChangelog{})
			if err != nil {
				return nil, err
			}
			v := i.(*models.TapdTaskChangelog)
			v.SourceId = data.Source.ID
			for _, fc := range taskChangelogRes.FieldChanges {
				item := &models.TapdTaskChangelogItem{
					SourceId:          data.Source.ID,
					ChangelogId:       v.ID,
					Field:             fc.Field,
					ValueBeforeParsed: fc.ValueBeforeParsed,
					ValueAfterParsed:  fc.ValueAfterParsed,
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
			results = append(results, v)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
