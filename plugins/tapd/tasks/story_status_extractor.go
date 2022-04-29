package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractStoryStatus

var ExtractStoryStatusMeta = core.SubTaskMeta{
	Name:             "extractStoryStatus",
	EntryPoint:       ExtractStoryStatus,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_bugStatus",
}

type TapdStoryStatusRes struct {
	Data map[string]string
}

func ExtractStoryStatus(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_STORY_STATUS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var statusRes TapdStoryStatusRes
			err := json.Unmarshal(row.Data, &statusRes)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0)
			for k, v := range statusRes.Data {
				toolL := &models.TapdStoryStatus{
					SourceId:    data.Source.ID,
					WorkspaceID: data.Options.WorkspaceID,
					EnglishName: k,
					ChineseName: v,
					IsLastStep:  false,
				}
				results = append(results, toolL)
			}

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
