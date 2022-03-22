package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jira/models"
)

func ConvertBoard(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("collect board:%d", data.Options.BoardId)
	idGen := didgen.NewDomainIdGenerator(&models.JiraBoard{})
	cursor, err := db.Model(&models.JiraBoard{}).Where("source_id = ? AND board_id = ?", data.Source.ID, data.Options.BoardId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				SourceId: data.Source.ID,
				BoardId:  data.Options.BoardId,
			},
			Table: RAW_BOARD_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraBoard{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			board := inputRow.(*models.JiraBoard)
			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{Id: idGen.Generate(data.Source.ID, data.Options.BoardId)},
				Name:         board.Name,
				Url:          board.Self,
			}
			return []interface{}{
				domainBoard,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
