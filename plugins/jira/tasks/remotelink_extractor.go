package tasks

import (
	"encoding/json"
	"regexp"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks/apiv2models"
)

func ExtractRemotelinks(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	sourceId := data.Source.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("extract remote links")
	var commitShaRegex *regexp.Regexp
	if pattern := data.Source.RemotelinkCommitShaPattern; pattern != "" {
		commitShaRegex = regexp.MustCompile(pattern)
	}

	// clean up issue_commits relationship for board
	err := db.Exec(`
	DELETE
	FROM _tool_jira_issue_commits
	WHERE source_id = ?
	  AND issue_id IN
	    (SELECT issue_id
	     FROM _tool_jira_board_issues
	     WHERE source_id = ?
	       AND board_id = ?)
	`, sourceId, sourceId, boardId).Error
	if err != nil {
		return err
	}

	// select all remotelinks belongs to the board, cursor is important for low memory footprint
	cursor, err := db.Model(&models.JiraRemotelink{}).
		Select("_tool_jira_remotelinks.*").
		Joins("left join _tool_jira_board_issues on _tool_jira_board_issues.issue_id = _tool_jira_remotelinks.issue_id").
		Where("_tool_jira_board_issues.board_id = ? AND _tool_jira_board_issues.source_id = ?", boardId, sourceId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				SourceId: sourceId,
				BoardId:  boardId,
			},
			Table: RAW_REMOTELINK_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var result []interface{}
			var raw apiv2models.RemoteLink
			err := json.Unmarshal(row.Data, &raw)
			if err != nil {
				return nil, err
			}
			var input apiv2models.Input
			err = json.Unmarshal(row.Input, &input)
			if err != nil {
				return nil, err
			}
			issue := &models.JiraIssue{SourceId: sourceId, IssueId: input.IssueId}
			err = db.Model(issue).Update("remotelink_updated", input.UpdateTime).Error
			if err != nil {
				return nil, err
			}
			remotelink := &models.JiraRemotelink{
				SourceId:     sourceId,
				RemotelinkId: raw.ID,
				IssueId:      input.IssueId,
				Self:         raw.Self,
				Title:        raw.Object.Title,
				Url:          raw.Object.URL,
			}
			result = append(result, remotelink)
			if commitShaRegex != nil {
				groups := commitShaRegex.FindStringSubmatch(remotelink.Url)
				if len(groups) > 1 {
					issueCommit := &models.JiraIssueCommit{
						SourceId:  sourceId,
						IssueId:   remotelink.IssueId,
						CommitSha: groups[1],
						CommitUrl: remotelink.Url,
					}
					result = append(result, issueCommit)
				}
			}
			return result, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}
