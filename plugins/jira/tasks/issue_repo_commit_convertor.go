package tasks

import (
	"reflect"
	"regexp"

	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

// ConvertIssueRepoCommits is to extract issue_repo_commits from jira_issue_commits, nothing difference with
// issue_commits but added a RepoUrl. This task is needed by EE group.
func ConvertIssueRepoCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDb()
	sourceId := data.Source.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("convert issue repo commits")
	var commitRepoUrlRegex *regexp.Regexp
	commitRepoUrlPattern := `(.*)\-\/commit`
	commitRepoUrlRegex = regexp.MustCompile(commitRepoUrlPattern)

	cursor, err := db.Table("_tool_jira_issue_commits jic").
		Joins(`left join _tool_jira_board_issues jbi on (
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
			item := &crossdomain.IssueRepoCommit{
				IssueId:   issueIdGenerator.Generate(sourceId, issueCommit.IssueId),
				CommitSha: issueCommit.CommitSha,
			}
			if commitRepoUrlRegex != nil {
				groups := commitRepoUrlRegex.FindStringSubmatch(issueCommit.CommitUrl)
				if len(groups) > 1 {
					item.RepoUrl = groups[1]
				}
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
