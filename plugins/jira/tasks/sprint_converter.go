package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

func ConvertSprint(sourceId uint64, boardId uint64) error {

	// select all sprints belongs to the board
	cursor, err := lakeModels.Db.Model(&jiraModels.JiraSprint{}).
		Select("jira_sprints.*").
		Joins("left join jira_board_sprints on jira_board_sprints.sprint_id = jira_sprints.sprint_id").
		Where("jira_board_sprints.source_id = ? AND jira_board_sprints.board_id = ?", sourceId, boardId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	boardOriginKey := okgen.NewOriginKeyGenerator(&jiraModels.JiraBoard{}).Generate(sourceId, boardId)
	sprintOriginKeyGenerator := okgen.NewOriginKeyGenerator(&jiraModels.JiraSprint{})
	issueOriginKeyGenerator := okgen.NewOriginKeyGenerator(&jiraModels.JiraIssue{})
	// iterate all rows
	for cursor.Next() {
		var jiraSprint jiraModels.JiraSprint
		err = lakeModels.Db.ScanRows(cursor, &jiraSprint)
		if err != nil {
			return err
		}
		sprint := &ticket.Sprint{
			DomainEntity: domainlayer.DomainEntity{
				OriginKey: sprintOriginKeyGenerator.Generate(jiraSprint.SourceId, jiraSprint.SprintId),
			},
			BoardOriginKey: boardOriginKey,
			Url:            jiraSprint.Self,
			State:          jiraSprint.State,
			Name:           jiraSprint.Name,
			StartDate:      jiraSprint.StartDate,
			EndDate:        jiraSprint.EndDate,
			CompleteDate:   jiraSprint.CompleteDate,
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(sprint).Error
		if err != nil {
			return err
		}
		var sprintIssues []jiraModels.JiraSprintIssue
		err = lakeModels.Db.Find(&sprintIssues, "source_id = ? AND sprint_id = ?", sourceId, jiraSprint.SprintId).Error
		if err != nil {
			return err
		}
		domainSprintIssues := make([]ticket.SprintIssue, 0, len(sprintIssues))
		for _, si := range sprintIssues{
			dsi := ticket.SprintIssue{
				SprintOriginKey: sprint.OriginKey,
				IssueOriginKey:  issueOriginKeyGenerator.Generate(sourceId, si.IssueId),
			}
			domainSprintIssues = append(domainSprintIssues, dsi)
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(&domainSprintIssues, BatchSize).Error
		if err != nil {
			return err
		}
	}
	return nil
}
