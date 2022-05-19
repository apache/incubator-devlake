package tasks

import (
	"fmt"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"reflect"
	"strings"
	"time"
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
				StartedDate:     (*time.Time)(iter.Startdate),
				EndedDate:       (*time.Time)(iter.Enddate),
				OriginalBoardID: WorkspaceIdGen.Generate(iter.ConnectionId, iter.WorkspaceID),
				CompletedDate:   (*time.Time)(iter.Completed),
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
