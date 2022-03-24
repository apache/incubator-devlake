package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractBoard

func ExtractBoard(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				SourceId: data.Source.ID,
				BoardId:  data.Options.BoardId,
			},
			Table: RAW_BOARD_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var board apiv2models.Board
			err := json.Unmarshal(row.Data, &board)
			if err != nil {
				return nil, err
			}
			return []interface{}{board.ToToolLayer(data.Source.ID)}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
