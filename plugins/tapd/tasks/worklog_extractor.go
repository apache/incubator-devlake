package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractWorklogs

var ExtractWorklogMeta = core.SubTaskMeta{
	Name:             "extractWorklogs",
	EntryPoint:       ExtractWorklogs,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table tapd_iterations",
}

type TapdWorklogRes struct {
	Timesheet models.TapdWorklogApiRes
}

func ExtractWorklogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_WORKLOG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var worklogBody TapdWorklogRes
			err := json.Unmarshal(row.Data, &worklogBody)
			if err != nil {
				return nil, err
			}
			worklogRes := worklogBody.Timesheet

			i, err := VoToDTO(&worklogRes, &models.TapdWorklog{})
			if err != nil {
				return nil, err
			}
			toolL := i.(*models.TapdWorklog)
			toolL.SourceId = data.Source.ID
			results := make([]interface{}, 0, 1)
			results = append(results, toolL)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
