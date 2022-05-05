package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractUsers

func ExtractUsers(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_USERS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var user apiv2models.User
			err := json.Unmarshal(row.Data, &user)
			if err != nil {
				return nil, err
			}
			return []interface{}{user.ToToolLayer(data.Connection.ID)}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
