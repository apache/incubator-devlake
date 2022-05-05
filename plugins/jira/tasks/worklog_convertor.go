package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

func ConvertWorklogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDb()
	connectionId := data.Connection.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("convert worklog")
	// select all worklogs belongs to the board
	cursor, err := db.Model(&models.JiraWorklog{}).
		Select("_tool_jira_worklogs.*").
		Joins(`left join _tool_jira_board_issues on (_tool_jira_board_issues.issue_id = _tool_jira_worklogs.issue_id)`).
		Where("_tool_jira_board_issues.connection_id = ? AND _tool_jira_board_issues.board_id = ?", connectionId, boardId).
		Rows()
	if err != nil {
		logger.Error("convert worklog error:", err)
		return err
	}
	defer cursor.Close()

	worklogIdGen := didgen.NewDomainIdGenerator(&models.JiraWorklog{})
	userIdGen := didgen.NewDomainIdGenerator(&models.JiraUser{})
	issueIdGen := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_WORKLOGS_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraWorklog{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			jiraWorklog := inputRow.(*models.JiraWorklog)
			worklog := &ticket.IssueWorklog{
				DomainEntity:     domainlayer.DomainEntity{Id: worklogIdGen.Generate(jiraWorklog.ConnectionId, jiraWorklog.IssueId, jiraWorklog.WorklogId)},
				IssueId:          issueIdGen.Generate(jiraWorklog.ConnectionId, jiraWorklog.IssueId),
				TimeSpentMinutes: jiraWorklog.TimeSpentSeconds / 60,
				StartedDate:      &jiraWorklog.Started,
				LoggedDate:       &jiraWorklog.Updated,
			}
			if jiraWorklog.AuthorId != "" {
				worklog.AuthorId = userIdGen.Generate(connectionId, jiraWorklog.AuthorId)
			}
			return []interface{}{worklog}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
