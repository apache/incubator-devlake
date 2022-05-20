package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/jira/models"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractWorklogs

func ExtractWorklogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDb()
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_WORKLOGS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var input apiv2models.Input
			err := json.Unmarshal(row.Input, &input)
			if err != nil {
				return nil, err
			}
			issue := &models.JiraIssue{ConnectionId: data.Connection.ID, IssueId: input.IssueId}
			err = db.Model(issue).Update("worklog_updated", input.UpdateTime).Error
			if err != nil {
				return nil, err
			}
			var worklog apiv2models.Worklog
			err = json.Unmarshal(row.Data, &worklog)
			if err != nil {
				return nil, err
			}
			return []interface{}{worklog.ToToolLayer(data.Connection.ID)}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
