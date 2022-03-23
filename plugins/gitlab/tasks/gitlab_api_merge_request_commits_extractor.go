package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ExtractApiMergeRequestsCommitsMeta = core.SubTaskMeta{
	Name:             "extractApiMergeRequestsCommits",
	EntryPoint:       ExtractApiMergeRequestsCommits,
	EnabledByDefault: true,
	Description:      "Extract raw merge requests commit data into tool layer table GitlabMergeRequestCommit and GitlabCommit",
}

func ExtractApiMergeRequestsCommits(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, _ := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_COMMITS_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			// create gitlab commit
			gitlabApiCommit := &GitlabApiCommit{}
			err := json.Unmarshal(row.Data, gitlabApiCommit)
			if err != nil {
				return nil, err
			}
			gitlabCommit, err := ConvertCommit(gitlabApiCommit)
			if err != nil {
				return nil, err
			}

			// get input info
			input := &GitlabInput{}
			err = json.Unmarshal(row.Input, input)
			if err != nil {
				return nil, err
			}

			gitlabMrCommit := &models.GitlabMergeRequestCommit{
				CommitSha:      gitlabApiCommit.GitlabId,
				MergeRequestId: input.GitlabId,
			}

			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 2)

			results = append(results, gitlabCommit)
			results = append(results, gitlabMrCommit)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
