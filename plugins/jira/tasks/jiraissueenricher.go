package tasks

import (
	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/jira/models"
)

var workloadCoefficient float64

func init() {
	workloadCoefficient = config.V.GetFloat64("JIRA_WORKLOAD_COEFFICIENT")
}

func EnrichIssues(boardId uint64) error {
	jiraIssue := &models.JiraIssue{}

	// select all issues belongs to the board
	cursor, err := lakeModels.Db.Model(jiraIssue).
		Select("jira_issues.*").
		Joins("left join jira_board_issues on jira_board_issues.issue_id = jira_issues.id").
		Where("jira_board_issues.board_id = ?", boardId).
		Rows()
	if err != nil {
		return err
	}

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jiraIssue)
		if err != nil {
			return err
		}
		if jiraIssue.ResolutionDate.Valid {
			jiraIssue.LeadTime = uint(jiraIssue.ResolutionDate.Time.Unix()-jiraIssue.Created.Unix()) / 60
		}
		jiraIssue.StdWorkload = uint(jiraIssue.Workload * workloadCoefficient)
		err = lakeModels.Db.Save(jiraIssue).Error
		if err != nil {
			return err
		}
	}
	return nil
}
