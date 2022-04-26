package tasks

import (
	"reflect"
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

type ChangelogItemResult struct {
	models.JiraChangelogItem
	IssueId           uint64 `gorm:"index"`
	AuthorAccountId   string
	AuthorDisplayName string
	Created           time.Time
}

func ConvertChangelogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	sourceId := data.Source.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("covert changelog")
	// select all changelogs belongs to the board
	cursor, err := db.Table("_tool_jira_changelog_items").
		Joins(`left join _tool_jira_changelogs on (
			_tool_jira_changelogs.source_id = _tool_jira_changelog_items.source_id
			AND _tool_jira_changelogs.changelog_id = _tool_jira_changelog_items.changelog_id
		)`).
		Joins(`left join _tool_jira_board_issues on (
			_tool_jira_board_issues.source_id = _tool_jira_changelogs.source_id
			AND _tool_jira_board_issues.issue_id = _tool_jira_changelogs.issue_id
		)`).
		Select("_tool_jira_changelog_items.*, _tool_jira_changelogs.issue_id, author_account_id, author_display_name, created").
		Where("_tool_jira_changelog_items.source_id = ? AND _tool_jira_board_issues.board_id = ?", sourceId, boardId).
		Rows()
	if err != nil {
		logger.Info(err.Error())
		return err
	}
	defer cursor.Close()
	sprintIssueConverter := NewSprintIssueConverter(taskCtx)
	issueIdGenerator := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	changelogIdGenerator := didgen.NewDomainIdGenerator(&models.JiraChangelogItem{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				SourceId: sourceId,
				BoardId:  boardId,
			},
			Table: RAW_CHANGELOG_TABLE,
		},
		InputRowType: reflect.TypeOf(ChangelogItemResult{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			row := inputRow.(*ChangelogItemResult)
			changelog := &ticket.Changelog{
				DomainEntity: domainlayer.DomainEntity{Id: changelogIdGenerator.Generate(
					row.SourceId,
					row.ChangelogId,
					row.Field,
				)},
				IssueId:     issueIdGenerator.Generate(row.SourceId, row.IssueId),
				AuthorId:    row.AuthorAccountId,
				AuthorName:  row.AuthorDisplayName,
				FieldId:     row.FieldId,
				FieldName:   row.Field,
				From:        row.FromString,
				To:          row.ToString,
				CreatedDate: row.Created,
			}
			sprintIssueConverter.FeedIn(sourceId, *row)
			return []interface{}{changelog}, nil
		},
	})
	if err != nil {
		logger.Info(err.Error())
		return err
	}

	err = converter.Execute()
	if err != nil {
		return err
	}
	return sprintIssueConverter.UpdateSprintIssue()
}
