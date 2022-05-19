package tasks

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractProjects

func ExtractProjects(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_PROJECT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var project apiv2models.Project
			err := json.Unmarshal(row.Data, &project)
			if err != nil {
				return nil, err
			}
			return []interface{}{project.ToToolLayer(data.Connection.ID)}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
