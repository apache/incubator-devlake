package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractIterations

var ExtractIterationMeta = core.SubTaskMeta{
	Name:             "extractIterations",
	EntryPoint:       ExtractIterations,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdIterationRes struct {
	Iteration models.TapdIterationRes
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
			var iterBody TapdIterationRes
			err := json.Unmarshal(row.Data, &iterBody)
			if err != nil {
				return nil, err
			}
			iterRes := iterBody.Iteration

			i, err := VoToDTO(&iterRes, &models.TapdIteration{})
			if err != nil {
				return nil, err
			}
			iter := i.(*models.TapdIteration)
			iter.SourceId = data.Source.ID
			workspaceIter := &models.TapdWorkspaceIteration{
				SourceId:    data.Source.ID,
				WorkspaceId: iter.WorkspaceId,
				IterationId: iter.ID,
			}
			return []interface{}{
				iter, workspaceIter,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
