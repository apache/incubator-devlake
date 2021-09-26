package tasks

import (
	"fmt"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/domainlayer/models/ticket"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
)

func ConvertIssues(boardId uint64) error {
	jiraBoard := &jiraModels.JiraBoard{}

	err := lakeModels.Db.First(jiraBoard, boardId).Error
	if err != nil {
		return err
	}

	board := &ticket.Board{
		DomainEntity: DomainEntity{
			OriginKey: fmt.Sprintf("jira:jira_board:%v", boardId),
		},
		Name: jiraBoard.Name,
		Url:  jiraBoard.Self,
	}

	err = lakeModels.Db.Clauses(clauses.OnConflict{UpdateAll: true}).Create(board)

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
		ticket.Board
		if jiraIssue.ResolutionDate.Valid {
			jiraIssue.LeadTime = uint(jiraIssue.ResolutionDate.Time.Unix()-jiraIssue.Created.Unix()) / 60
		}
		jiraIssue.StdStoryPoint = uint(jiraIssue.StoryPoint * storyPointCoefficient)
		jiraIssue.StdType = getStdType(jiraIssue)
		jiraIssue.StdStatus = getStdStatus(jiraIssue)
		err = lakeModels.Db.Save(jiraIssue).Error
		if err != nil {
			return err
		}
	}
	return nil
}
