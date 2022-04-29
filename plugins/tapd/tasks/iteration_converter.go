package tasks

import (
	"fmt"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"reflect"
	"strings"
)

func ConvertIteration(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("collect board:%d", data.Options.WorkspaceID)
	iterIdGen := didgen.NewDomainIdGenerator(&models.TapdIteration{})
	cursor, err := db.Model(&models.TapdIteration{}).Where("connection_id = ? AND workspace_id = ?", data.Connection.ID, data.Options.WorkspaceID).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,

				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_ITERATION_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdIteration{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			iter := inputRow.(*models.TapdIteration)
			domainIter := &ticket.Sprint{
				DomainEntity:    domainlayer.DomainEntity{Id: iterIdGen.Generate(data.Connection.ID, iter.ID)},
				Url:             fmt.Sprintf("https://www.tapd.cn/%d/prong/iterations/view/%d", iter.WorkspaceID, iter.ID),
				Status:          strings.ToUpper(iter.Status),
				Name:            iter.Name,
				StartedDate:     iter.Startdate.ToNullableTime(),
				EndedDate:       iter.Enddate.ToNullableTime(),
				OriginalBoardID: WorkspaceIdGen.Generate(iter.ConnectionId, iter.WorkspaceID),
				CompletedDate:   iter.Completed.ToNullableTime(),
			}
			results := make([]interface{}, 0)
			results = append(results, domainIter)
			boardSprint := &ticket.BoardSprint{
				BoardId:  domainIter.OriginalBoardID,
				SprintId: domainIter.Id,
			}
			results = append(results, boardSprint)
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertIterationMeta = core.SubTaskMeta{
	Name:             "convertIteration",
	EntryPoint:       ConvertIteration,
	EnabledByDefault: true,
	Description:      "convert Tapd iteration",
}
