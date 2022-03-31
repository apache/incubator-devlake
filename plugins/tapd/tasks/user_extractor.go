package tasks

//
//import (
//	"encoding/json"
//	"github.com/merico-dev/lake/plugins/core"
//	"github.com/merico-dev/lake/plugins/helper"
//	"github.com/merico-dev/lake/plugins/tapd/models"
//)
//
//var _ core.SubTaskEntryPoint = ExtractUsers
//
//var ExtractUsersMeta = core.SubTaskMeta{
//	Name:             "extractUsers",
//	EntryPoint:       ExtractUsers,
//	EnabledByDefault: true,
//	Description:      "Extract raw workspace data into tool layer table tapd_users",
//}
//
//type TapdUserRes struct {
//	User models.TapdUser
//}
//
//func ExtractUsers(taskCtx core.SubTaskContext) error {
//	data := taskCtx.GetData().(*TapdTaskData)
//	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
//		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
//			Ctx: taskCtx,
//			Params: TapdApiParams{
//				SourceId: data.Source.ID,
//				//CompanyId: data.Options.CompanyId,
//				WorkspaceId: data.Options.WorkspaceId,
//			},
//			Table: RAW_ITERATION_TABLE,
//		},
//		Extract: func(row *helper.RawData) ([]interface{}, error) {
//			var iterRes TapdUserRes
//			err := json.Unmarshal(row.Data, &iterRes)
//			if err != nil {
//				return nil, err
//			}
//			results := make([]interface{}, 0, 1)
//			iterRes.User.SourceId = data.Source.ID
//			results = append(results, &iterRes.User)
//			return results, nil
//		},
//	})
//
//	if err != nil {
//		return err
//	}
//
//	return extractor.Execute()
//}
