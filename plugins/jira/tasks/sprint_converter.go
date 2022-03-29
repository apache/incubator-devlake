package tasks

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"strings"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

func ConvertSprints(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	sourceId := data.Source.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert sprints")
	cursor, err := db.Model(&models.JiraSprint{}).
		Select("jira_sprints.*").
		Joins("left join jira_board_sprints on jira_board_sprints.sprint_id = jira_sprints.sprint_id").
		Where("jira_board_sprints.source_id = ? AND jira_board_sprints.board_id = ?", sourceId, boardId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	var converter *helper.DataConverter
	domainBoardId := didgen.NewDomainIdGenerator(&models.JiraBoard{}).Generate(sourceId, boardId)
	sprintIdGen := didgen.NewDomainIdGenerator(&models.JiraSprint{})
	issueIdGen := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.JiraBoard{})
	converter, err = helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				SourceId: data.Source.ID,
				BoardId:  data.Options.BoardId,
			},
			Table: RAW_SPRINT_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraSprint{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			var result []interface{}
			jiraSprint := inputRow.(*models.JiraSprint)
			sprint := &ticket.Sprint{
				DomainEntity:  domainlayer.DomainEntity{Id: sprintIdGen.Generate(sourceId, jiraSprint.SprintId)},
				Url:           jiraSprint.Self,
				Status:        strings.ToUpper(jiraSprint.State),
				Name:          jiraSprint.Name,
				StartedDate:   jiraSprint.StartDate,
				EndedDate:     jiraSprint.EndDate,
				CompletedDate: jiraSprint.CompleteDate,
				OriginBoardID: boardIdGen.Generate(sourceId, jiraSprint.OriginBoardID),
			}
			result = append(result, sprint)
			var sprintIssues []models.JiraSprintIssue
			err = db.Find(&sprintIssues, "source_id = ? AND sprint_id = ?", sourceId, jiraSprint.SprintId).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}
			domainSprintIssues := make([]ticket.SprintIssue, 0, len(sprintIssues))
			for _, si := range sprintIssues {
				dsi := ticket.SprintIssue{
					SprintId:  sprint.Id,
					IssueId:   issueIdGen.Generate(sourceId, si.IssueId),
					AddedDate: si.IssueCreatedDate,
				}
				if dsi.AddedDate != nil {
					dsi.AddedStage = getStage(*dsi.AddedDate, sprint.StartedDate, sprint.CompletedDate)
				}
				if si.ResolutionDate != nil {
					dsi.ResolvedStage = getStage(*si.ResolutionDate, sprint.StartedDate, sprint.CompletedDate)
				}
				domainSprintIssues = append(domainSprintIssues, dsi)
			}
			if len(domainSprintIssues) > 0 {
				err = db.Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns([]string{"resolved_stage"})}).Create(domainSprintIssues).Error
				if err != nil {
					return nil, err
				}

			}
			boardSprint := &ticket.BoardSprint{
				BoardId:  domainBoardId,
				SprintId: sprint.Id,
			}
			result = append(result, boardSprint)
			return result, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
