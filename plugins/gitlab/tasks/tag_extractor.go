package tasks

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type GitlabApiTag struct {
	Name      string
	Message   string
	Target    string
	Protected bool
	Release   struct {
		TagName     string
		Description string
	}
}

var ExtractTagMeta = core.SubTaskMeta{
	Name:             "extractApiTag",
	EntryPoint:       ExtractApiTag,
	EnabledByDefault: true,
	Description:      "Extract raw tag data into tool layer table GitlabTag",
}

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

// Convert the API response to our DB model instance
func convertTag(tag *GitlabApiTag) (*models.GitlabTag, error) {
	gitlabTag := &models.GitlabTag{
		Name:               tag.Name,
		Message:            tag.Message,
		Target:             tag.Target,
		Protected:          tag.Protected,
		ReleaseDescription: tag.Release.Description,
	}
	return gitlabTag, nil
}
