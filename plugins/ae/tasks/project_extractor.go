package tasks

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/plugins/ae/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type ApiProjectResponse struct {
	Id           int        `json:"id"`
	GitUrl       string     `json:"git_url"`
	Priority     int        `json:"priority"`
	AECreateTime *time.Time `json:"create_time"`
	AEUpdateTime *time.Time `json:"update_time"`
}

func ExtractProject(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*AeTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: AeApiParams{
				ProjectId: data.Options.ProjectId,
			},
			Table: RAW_PROJECT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &ApiProjectResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			aeProject := &models.AEProject{
				Id:           strconv.Itoa(body.Id),
				GitUrl:       body.GitUrl,
				Priority:     body.Priority,
				AECreateTime: body.AECreateTime,
				AEUpdateTime: body.AEUpdateTime,
			}
			return []interface{}{aeProject}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractProjectMeta = core.SubTaskMeta{
	Name:             "extractProject",
	EntryPoint:       ExtractProject,
	EnabledByDefault: true,
	Description:      "Extract raw project data into tool layer table ae_projects",
}
