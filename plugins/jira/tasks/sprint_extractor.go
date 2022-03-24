package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/jira/models"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractSprints

func ExtractSprints(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				SourceId: data.Source.ID,
				BoardId:  data.Options.BoardId,
			},
			Table: RAW_SPRINT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var sprint apiv2models.Sprint
			err := json.Unmarshal(row.Data, &sprint)
			if err != nil {
				return nil, err
			}
			boardSprint := models.JiraBoardSprint{
				SourceId: data.Source.ID,
				BoardId:  data.Options.BoardId,
				SprintId: sprint.ID,
			}
			return []interface{}{sprint.ToToolLayer(data.Source.ID), &boardSprint}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
