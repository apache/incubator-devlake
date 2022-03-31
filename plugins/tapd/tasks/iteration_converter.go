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
	logger.Info("collect board:%d", data.Options.WorkspaceId)
	iterIdGen := didgen.NewDomainIdGenerator(&models.TapdIteration{})
	cursor, err := db.Model(&models.TapdIteration{}).Where("source_id = ? AND workspace_id = ?", data.Source.ID, data.Options.WorkspaceId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId:   data.Source.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_ITERATION_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdIteration{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			iter := inputRow.(*models.TapdIteration)
			domainIter := &ticket.Sprint{
				DomainEntity:  domainlayer.DomainEntity{Id: iterIdGen.Generate(data.Source.ID, iter.ID)},
				Url:           fmt.Sprintf("https://www.tapd.cn/%s/prong/iterations/view/%s", iter.WorkspaceID, iter.ID),
				Status:        strings.ToUpper(iter.Status),
				Name:          iter.Name,
				StartedDate:   iter.Startdate,
				EndedDate:     iter.Enddate,
				OriginBoardID: WorkspaceIdGen.Generate(iter.SourceId, iter.WorkspaceID),
				CompletedDate: iter.Completed,
			}

			return []interface{}{
				domainIter,
			}, nil
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
