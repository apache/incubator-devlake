package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
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

	domainBoardId := didgen.NewDomainIdGenerator(&jiraModels.JiraBoard{}).Generate(sourceId, boardId)
	sprintIdGen := didgen.NewDomainIdGenerator(&jiraModels.JiraSprint{})
	issueIdGen := didgen.NewDomainIdGenerator(&jiraModels.JiraIssue{})
	// iterate all rows
	for cursor.Next() {
		var jiraSprint jiraModels.JiraSprint
		err = lakeModels.Db.ScanRows(cursor, &jiraSprint)
		if err != nil {
			return err
		}
		sprint := &ticket.Sprint{
			DomainEntity:domainlayer.DomainEntity{Id: sprintIdGen.Generate(jiraSprint.SourceId, jiraSprint.SprintId)},
			Url:           jiraSprint.Self,
			Status:        jiraSprint.State,
			Name:          jiraSprint.Name,
			StartedDate:   jiraSprint.StartDate,
			EndedDate:     jiraSprint.EndDate,
			CompletedDate: jiraSprint.CompleteDate,
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
		for _, si := range sprintIssues {
			dsi := ticket.SprintIssue{
				SprintId: sprint.Id,
				IssueId:  issueIdGen.Generate(sourceId, si.IssueId),
			}
			if si.ResolutionDate != nil {
				dsi.ResolvedStage = getStage(*si.ResolutionDate, sprint.StartedDate, sprint.CompletedDate)
			}
			domainSprintIssues = append(domainSprintIssues, dsi)
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Select("sprint_id", "issue_id", "resolved_stage").CreateInBatches(&domainSprintIssues, BatchSize).Error
		if err != nil {
			return err
		}
		boardSprint := &ticket.BoardSprint{
			BoardId:  domainBoardId,
			SprintId: sprint.Id,
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(boardSprint).Error
		if err != nil {
			return err
		}
	}
	return nil
}
