package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractTasks

var ExtractTaskMeta = core.SubTaskMeta{
	Name:             "extractTasks",
	EntryPoint:       ExtractTasks,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table tapd_iterations",
}

type TapdTaskRes struct {
	Task models.TapdTaskApiRes
}

func ExtractTasks(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_TASK_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var taskBody TapdTaskRes
			err := json.Unmarshal(row.Data, &taskBody)
			if err != nil {
				return nil, err
			}
			taskRes := taskBody.Task

			i, err := ResToDb(&taskRes, &models.TapdTask{})
			if err != nil {
				return nil, err
			}
			task := i.(*models.TapdTask)
			task.SourceId = data.Source.ID
			return []interface{}{
				task,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
