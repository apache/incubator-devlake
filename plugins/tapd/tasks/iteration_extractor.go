package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"strconv"
	"time"
)

var _ core.SubTaskEntryPoint = ExtractIterations

var ExtractIterationsMeta = core.SubTaskMeta{
	Name:             "extractIterations",
	EntryPoint:       ExtractIterations,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table tapd_iterations",
}

type TapdIterationRes struct {
	Iteration models.TapdIterationRes
}

func ExtractIterations(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_ITERATION_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var iterRes TapdIterationRes
			err := json.Unmarshal(row.Data, &iterRes)
			if err != nil {
				return nil, err
			}
			idInt, err := strconv.Atoi(iterRes.Iteration.ID)
			if err != nil {
				return nil, err
			}
			start, err := time.Parse(shortForm, iterRes.Iteration.Startdate)
			if err != nil {
				return nil, err
			}
			end, err := time.Parse(shortForm, iterRes.Iteration.Enddate)
			if err != nil {
				return nil, err
			}
			iteration := models.TapdIteration{
				SourceId:     data.Options.SourceId,
				ID:           uint64(idInt),
				Name:         iterRes.Iteration.Name,
				WorkspaceID:  data.Options.WorkspaceId,
				Startdate:    &start,
				Enddate:      &end,
				Status:       iterRes.Iteration.Status,
				ReleaseID:    iterRes.Iteration.ReleaseID,
				Description:  iterRes.Iteration.Description,
				Creator:      iterRes.Iteration.Creator,
				Releaseowner: iterRes.Iteration.Releaseowner,
				Notice:       iterRes.Iteration.Notice,
				Releasename:  iterRes.Iteration.Releasename,
			}
			if iterRes.Iteration.Completed != "" {
				v, _ := time.Parse(longForm, iterRes.Iteration.Completed)
				iteration.Completed = &v
			}
			if iterRes.Iteration.Created != "" {
				v, _ := time.Parse(longForm, iterRes.Iteration.Created)
				iteration.Created = &v
			}
			if iterRes.Iteration.Modified != "" {
				v, _ := time.Parse(longForm, iterRes.Iteration.Modified)
				iteration.Modified = &v
			}
			if iterRes.Iteration.Launchdate != "" {
				v, _ := time.Parse(longForm, iterRes.Iteration.Launchdate)
				iteration.Launchdate = &v
			}
			return []interface{}{
				&iteration,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
