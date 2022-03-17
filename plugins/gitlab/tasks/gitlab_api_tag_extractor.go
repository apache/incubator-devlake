package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

func ExtractApiTag(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GitlabTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,

			Params: GitlabApiParams{
				ProjectId: data.Options.ProjectId,
			},
			Table: RAW_TAG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			// need to extract 3 kinds of entities here
			results := make([]interface{}, 0, 3)

			// create gitlab commit
			gitlabApiTag := &GitlabApiTag{}
			err := json.Unmarshal(row.Data, gitlabApiTag)
			if err != nil {
				return nil, err
			}
			gitlabTag, err := convertTag(gitlabApiTag)
			if err != nil {
				return nil, err
			}
			results = append(results, gitlabTag)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
