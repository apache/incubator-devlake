package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

func ConvertIssues(sourceId uint64, boardId uint64) error {

	jiraIssue := &jiraModels.JiraIssue{}
	// select all issues belongs to the board
	cursor, err := lakeModels.Db.Model(jiraIssue).
		Select("jira_issues.*").
		Joins("left join jira_board_issues on jira_board_issues.issue_id = jira_issues.issue_id").
		Where("jira_board_issues.source_id = ? AND jira_board_issues.board_id = ?", sourceId, boardId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	boardOriginKey := okgen.NewOriginKeyGenerator(&jiraModels.JiraBoard{}).Generate(sourceId, boardId)
	issueOriginKeyGenerator := okgen.NewOriginKeyGenerator(&jiraModels.JiraIssue{})
	userOriginKeyGenerator := okgen.NewOriginKeyGenerator(&jiraModels.JiraUser{})
	sprintOriginKeyGenerator := okgen.NewOriginKeyGenerator(&jiraModels.JiraSprint{})

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jiraIssue)
		if err != nil {
			return err
		}
		issue := &ticket.Issue{
			DomainEntity: domainlayer.DomainEntity{
				OriginKey: issueOriginKeyGenerator.Generate(jiraIssue.SourceId, jiraIssue.IssueId),
			},
			BoardOriginKey:           boardOriginKey,
			Url:                      jiraIssue.Self,
			Key:                      jiraIssue.Key,
			Summary:                  jiraIssue.Summary,
			EpicKey:                  jiraIssue.EpicKey,
			Type:                     jiraIssue.StdType,
			Status:                   jiraIssue.StdStatus,
			StoryPoint:               jiraIssue.StdStoryPoint,
			OriginalEstimateMinutes:  jiraIssue.OriginalEstimateMinutes,
			AggregateEstimateMinutes: jiraIssue.AggregateEstimateMinutes,
			RemainingEstimateMinutes: jiraIssue.RemainingEstimateMinutes,
			CreatorOriginKey:         userOriginKeyGenerator.Generate(sourceId, jiraIssue.CreatorAccountId),
			ResolutionDate:           jiraIssue.ResolutionDate,
			Priority:                 jiraIssue.PriorityName,
			CreatedDate:              jiraIssue.Created,
			UpdatedDate:              jiraIssue.Updated,
			LeadTimeMinutes:          jiraIssue.LeadTimeMinutes,
			SpentMinutes:             jiraIssue.SpentMinutes,
		}
		if jiraIssue.AssigneeAccountId != "" {
			issue.AssigneeOriginKey = userOriginKeyGenerator.Generate(sourceId, jiraIssue.AssigneeAccountId)
		}
		if jiraIssue.ParentId != 0 {
			issue.ParentOriginKey = issueOriginKeyGenerator.Generate(sourceId, jiraIssue.ParentId)
		}
		if jiraIssue.SprintId != 0 {
			issue.SprintOriginKey = sprintOriginKeyGenerator.Generate(sourceId, jiraIssue.SprintId)
		}

		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(issue).Error
		if err != nil {
			return err
		}
	}
	return nil
}
