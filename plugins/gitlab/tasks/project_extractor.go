package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

var ExtractProjectMeta = core.SubTaskMeta{
	Name:             "extractApiProject",
	EntryPoint:       ExtractApiProject,
	EnabledByDefault: true,
	Description:      "Extract raw project data into tool layer table GitlabProject",
}

func ExtractApiProject(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, _ := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			// create gitlab commit
			gitlabApiProject := &GitlabApiProject{}
			err := json.Unmarshal(row.Data, gitlabApiProject)
			if err != nil {
				return nil, err
			}
			gitlabProject := convertProject(gitlabApiProject)

			results := make([]interface{}, 0, 1)
			results = append(results, gitlabProject)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
