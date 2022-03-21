package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

func ExtractApiMergeRequestsNotes(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, _ := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_NOTES_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			mrNote := &MergeRequestNote{}
			err := json.Unmarshal(row.Data, mrNote)
			if err != nil {
				return nil, err
			}

			gitlabMergeRequestNote, err := convertMergeRequestNote(mrNote)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 1)

			results = append(results, gitlabMergeRequestNote)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
