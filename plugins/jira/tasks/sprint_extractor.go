package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/jira/models"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractSprints

func ExtractSprints(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
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
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
				SprintId:     sprint.ID,
			}
			return []interface{}{sprint.ToToolLayer(data.Connection.ID), &boardSprint}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
