package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"strconv"
	"time"
)

var _ core.SubTaskEntryPoint = ExtractWorkspaces

var ExtractWorkspacesMeta = core.SubTaskMeta{
	Name:             "extractWorkspaces",
	EntryPoint:       ExtractWorkspaces,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table tapd_workspaces",
}

type TapdWorkspaceRes struct {
	Workspace models.TapdWorkspaceApiRes
}

func ExtractWorkspaces(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_WORKSPACE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var workspaceRes TapdWorkspaceRes
			err := json.Unmarshal(row.Data, &workspaceRes)
			if err != nil {
				return nil, err
			}
			idInt, err := strconv.Atoi(workspaceRes.Workspace.ID)
			if err != nil {
				return nil, err
			}
			tmp := workspaceRes.Workspace
			start, err := time.Parse(shortForm, tmp.BeginDate)
			if err != nil {
				return nil, err
			}
			end, err := time.Parse(shortForm, tmp.EndDate)
			if err != nil {
				return nil, err
			}
			workSpace := &models.TapdWorkspace{
				SourceId:    data.Source.ID,
				ID:          uint64(idInt),
				Name:        tmp.Name,
				PrettyName:  tmp.PrettyName,
				Category:    tmp.Category,
				Status:      tmp.Status,
				Description: tmp.Description,
				BeginDate:   &start,
				EndDate:     &end,
				ExternalOn:  tmp.ExternalOn,
				Creator:     tmp.Creator,
				Created:     tmp.Created,
			}
			return []interface{}{
				workSpace,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
