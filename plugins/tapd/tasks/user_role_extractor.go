package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"strconv"
)

var _ core.SubTaskEntryPoint = ExtractUserRoles

var ExtractUserRolesMeta = core.SubTaskMeta{
	Name:             "extractUserRoles",
	EntryPoint:       ExtractUserRoles,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table tapd_iterations",
}

type TapdUserRoleRes map[string]string

func ExtractUserRoles(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_USER_ROLE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var iterRes TapdUserRoleRes
			err := json.Unmarshal(row.Data, &iterRes)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, len(iterRes))
			for k, v := range iterRes {
				userRole := models.TapdUserRole{
					SourceId:    data.Source.ID,
					WorkspaceId: strconv.FormatUint(data.Options.WorkspaceId, 10),
					ID:          k,
					Name:        v,
				}
				results = append(results, &userRole)

			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
