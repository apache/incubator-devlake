package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractChangelogs

func ExtractChangelogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	if data.JiraServerInfo.DeploymentType == models.DeploymentServer {
		return nil
	}
	db := taskCtx.GetDb()
	connectionId := data.Connection.ID
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_CHANGELOG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var input apiv2models.Input
			err := json.Unmarshal(row.Input, &input)
			if err != nil {
				return nil, err
			}
			var result []interface{}
			var changelog apiv2models.Changelog
			err = json.Unmarshal(row.Data, &changelog)
			if err != nil {
				return nil, err
			}
			issue := &models.JiraIssue{ConnectionId: connectionId, IssueId: input.IssueId}
			err = db.Model(issue).Update("changelog_updated", input.UpdateTime).Error
			if err != nil {
				return nil, err
			}
			cl, user := changelog.ToToolLayer(connectionId, input.IssueId)
			result = append(result, cl, user)
			for _, item := range changelog.Items {
				result = append(result, item.ToToolLayer(connectionId, changelog.ID))
			}
			return result, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
