package tasks

import (
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
	connectionId := data.Connection.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert sprints")
	cursor, err := db.Model(&models.JiraSprint{}).
		Select("_tool_jira_sprints.*").
		Joins("left join _tool_jira_board_sprints on _tool_jira_board_sprints.sprint_id = _tool_jira_sprints.sprint_id").
		Where("_tool_jira_board_sprints.connection_id = ? AND _tool_jira_board_sprints.board_id = ?", connectionId, boardId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	var converter *helper.DataConverter
	domainBoardId := didgen.NewDomainIdGenerator(&models.JiraBoard{}).Generate(connectionId, boardId)
	sprintIdGen := didgen.NewDomainIdGenerator(&models.JiraSprint{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.JiraBoard{})
	converter, err = helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_SPRINT_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraSprint{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			var result []interface{}
			jiraSprint := inputRow.(*models.JiraSprint)
			sprint := &ticket.Sprint{
				DomainEntity:    domainlayer.DomainEntity{Id: sprintIdGen.Generate(connectionId, jiraSprint.SprintId)},
				Url:             jiraSprint.Self,
				Status:          strings.ToUpper(jiraSprint.State),
				Name:            jiraSprint.Name,
				StartedDate:     jiraSprint.StartDate,
				EndedDate:       jiraSprint.EndDate,
				CompletedDate:   jiraSprint.CompleteDate,
				OriginalBoardID: boardIdGen.Generate(connectionId, jiraSprint.OriginBoardID),
			}
			result = append(result, sprint)
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
