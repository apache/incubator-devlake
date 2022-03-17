package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

func ExtractApiTag(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, _ := CreateRawDataSubTaskArgs(taskCtx, RAW_TAG_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			// need to extract 1 kinds of entities here
			results := make([]interface{}, 0, 1)

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
