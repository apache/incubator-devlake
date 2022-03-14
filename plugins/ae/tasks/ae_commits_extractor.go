package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/ae/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

type ApiCommitsResponse []AeApiCommit

type AeApiCommit struct {
	HexSha      string `json:"hexsha"`
	AnalysisId  string `json:"analysis_id"`
	AuthorEmail string `json:"author_email"`
	DevEq       int    `json:"dev_eq"`
}

func ExtractCommits(taskCtx core.SubTaskContext) error {
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
			body := &ApiCommitsResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, len(*body))
			for _, apiCommit := range *body {
				results = append(results, models.AECommit{
					HexSha:      apiCommit.HexSha,
					AnalysisId:  apiCommit.AnalysisId,
					AuthorEmail: apiCommit.AuthorEmail,
					DevEq:       apiCommit.DevEq,
					AEProjectId: data.Options.ProjectId,
				})
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractCommitsMeta = core.SubTaskMeta{
	Name:             "extractCommits",
	EntryPoint:       ExtractCommits,
	EnabledByDefault: true,
	Description:      "Extract raw commit data into tool layer table ae_commits",
}
