package tasks

import (
	"encoding/json"
	"regexp"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

func EnrichRemotelinks(source *models.JiraSource, boardId uint64) (err error) {
	// prepare for issue_commits enrichment
	var commitShaRegex *regexp.Regexp
	if source.RemotelinkCommitShaPattern != "" {
		commitShaRegex = regexp.MustCompile(source.RemotelinkCommitShaPattern)
	}
	// clean up issue_commits relationship for board
	lakeModels.Db.Exec(`
	DELETE ic
	FROM jira_issue_commits ic
	LEFT JOIN jira_board_issues bi ON (bi.source_id = ic.source_id AND bi.issue_id = ic.issue_id)
	WHERE ic.source_id = ? AND bi.board_id = ?
	`, source.ID, boardId)

	var remotelink *models.JiraRemotelink
	var issueCommit *models.JiraIssueCommit
	// select all remotelinks belongs to the board, cursor is important for low memory footprint
	cursor, err := lakeModels.Db.Model(&models.JiraRemotelink{}).
		Select("jira_remotelinks.*").
		Joins("left join jira_board_issues on jira_board_issues.issue_id = jira_remotelinks.issue_id").
		Where("jira_board_issues.board_id = ? AND jira_board_issues.source_id = ?", boardId, source.ID).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	// save in batch, this could improvement performance dramatically
	size := 1000
	i := 0
	j := 0
	batch := make([]models.JiraRemotelink, size)
	icBatch := make([]models.JiraIssueCommit, size)
	saveBatch := func() error {
		err := lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).CreateInBatches(batch[:i], size).Error
		if err != nil {
			return err
		}
		return nil
	}
	saveIcBatch := func() error {
		err := lakeModels.Db.Clauses(clause.OnConflict{
			DoNothing: true,
		}).CreateInBatches(icBatch[:j], size).Error
		if err != nil {
			return err
		}
		return nil
	}

	apiRemotelink := &JiraApiRemotelink{}
	// iterate all rows
	for cursor.Next() {
		if i >= size {
			err = saveBatch()
			if err != nil {
				return err
			}
			i = 0
		}
		if j >= size {
			err = saveIcBatch()
			if err != nil {
				return err
			}
			j = 0
		}

		remotelink = &batch[i]
		err = lakeModels.Db.ScanRows(cursor, remotelink)
		if err != nil {
			return err
		}
		err = json.Unmarshal(remotelink.RawJson, apiRemotelink)
		if err != nil {
			return err
		}
		remotelink.Self = apiRemotelink.Self
		remotelink.Url = apiRemotelink.Object.Url
		remotelink.Title = apiRemotelink.Object.Title

		// issue commit relationship
		if commitShaRegex != nil {
			groups := commitShaRegex.FindStringSubmatch(remotelink.Url)
			if len(groups) > 1 {
				issueCommit = &icBatch[j]
				issueCommit.SourceId = source.ID
				issueCommit.IssueId = remotelink.IssueId
				issueCommit.CommitSha = groups[1]
				j++
			}
		}

		i++
	}
	if i > 0 {
		err = saveBatch()
		if err != nil {
			return err
		}
	}
	if j > 0 {
		err = saveIcBatch()
		if err != nil {
			return err
		}
	}
	return nil
}
