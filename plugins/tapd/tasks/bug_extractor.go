package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractBugs

var ExtractBugMeta = core.SubTaskMeta{
	Name:             "extractBugs",
	EntryPoint:       ExtractBugs,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table tapd_iterations",
}

type TapdBugRes struct {
	Bug models.TapdBugApiRes
}

func ExtractBugs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_BUG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var bugBody TapdBugRes
			err := json.Unmarshal(row.Data, &bugBody)
			if err != nil {
				return nil, err
			}
			bugRes := bugBody.Bug

			i, err := ResToDb(&bugRes, &models.TapdBug{})
			if err != nil {
				return nil, err
			}
			bug := i.(*models.TapdBug)
			bug.SourceId = data.Source.ID
			return []interface{}{
				bug,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
