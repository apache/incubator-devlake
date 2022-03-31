package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractIterations

var ExtractIterationsMeta = core.SubTaskMeta{
	Name:             "extractIterations",
	EntryPoint:       ExtractIterations,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table tapd_iterations",
}

type TapdIterationRes struct {
	Iteration models.TapdIteration
}

func ExtractIterations(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_ITERATION_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var iterRes TapdIterationRes
			err := json.Unmarshal(row.Data, &iterRes)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 1)
			iterRes.Iteration.SourceId = data.Source.ID
			results = append(results, &iterRes.Iteration)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
