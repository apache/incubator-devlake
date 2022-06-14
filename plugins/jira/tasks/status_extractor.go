package tasks

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

func ExtractStatus(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Connection.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("extract Status, connection_id=%d, board_id=%d", connectionId, boardId)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_STATUS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var apiStatus apiv2models.Status
			err := json.Unmarshal(row.Data, &apiStatus)
			if err != nil {
				return nil, err
			}
			if apiStatus.Scope != nil {
				// FIXME: skip scope status
				return nil, nil
			}
			var jiraStatus = models.JiraStatus{
				ConnectionId:   connectionId,
				ID:             apiStatus.ID,
				Name:           apiStatus.Name,
				Self:           apiStatus.Self,
				StatusCategory: apiStatus.StatusCategory.Key,
			}
			var result []interface{}
			result = append(result, jiraStatus)
			return result, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
