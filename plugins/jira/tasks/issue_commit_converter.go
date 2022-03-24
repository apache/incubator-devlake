package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

func ConvertIssueCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDb()
	sourceId := data.Source.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("convert issue commits")

	cursor, err := db.Table("jira_issue_commits jic").
		Joins(`left join jira_board_issues jbi on (
			jbi.source_id = jic.source_id
			AND jbi.issue_id = jic.issue_id
		)`).
		Select("jic.*").
		Where("jbi.source_id = ? AND jbi.board_id = ?", sourceId, boardId).
		Order("jbi.source_id, jbi.issue_id").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGenerator := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				SourceId: sourceId,
				BoardId:  boardId,
			},
			Table: RAW_REMOTELINK_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraIssueCommit{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			var result []interface{}
			issueCommit := inputRow.(*models.JiraIssueCommit)
			item := &crossdomain.IssueCommit{
				IssueId:   issueIdGenerator.Generate(sourceId, issueCommit.IssueId),
				CommitSha: issueCommit.CommitSha,
			}
			result = append(result, item)
			return result, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
