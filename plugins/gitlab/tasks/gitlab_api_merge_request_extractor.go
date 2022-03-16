package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

func ExtractApiMergeRequests(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GitlabTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,

			Params: GitlabApiParams{
				ProjectId: data.Options.ProjectId,
			},
			Table: RAW_MERGE_REQUEST_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			mr := &MergeRequestRes{}
			err := json.Unmarshal(row.Data, mr)
			if err != nil {
				return nil, err
			}

			gitlabMergeRequest, err := convertMergeRequest(mr, data.Options.ProjectId)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, len(mr.Reviewers)+1)

			results = append(results, gitlabMergeRequest)

			for _, reviewer := range mr.Reviewers {
				gitlabReviewer := NewReviewer(data.Options.ProjectId, mr.GitlabId, reviewer)
				results = append(results, gitlabReviewer)
			}

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
