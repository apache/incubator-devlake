package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/jenkins/models"
	"strconv"
	"time"
)

// this struct should be moved to `gitub_api_common.go`

var ExtractApiBuildsMeta = core.SubTaskMeta{
	Name:             "extractApiBuilds",
	EntryPoint:       ExtractApiBuilds,
	EnabledByDefault: true,
	Description:      "Extract raw builds data into tool layer table jenkins_builds",
}

func ExtractApiBuilds(taskCtx core.SubTaskContext) error {
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			/*
				Table store raw data
			*/
			Table: RAW_BUILD_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &models.ApiBuildResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}

			input := &SimpleJob{}
			err = json.Unmarshal(row.Input, input)
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0, 1)

			build := &models.JenkinsBuild{
				JobName:           input.Name,
				Duration:          body.Duration,
				DisplayName:       body.DisplayName,
				EstimatedDuration: body.EstimatedDuration,
				Number:            body.Number,
				Result:            body.Result,
				Timestamp:         body.Timestamp,
				StartTime:         time.Unix(body.Timestamp/1000, 0),
			}

			vcs := body.ChangeSet.Kind
			if vcs == "git" || vcs == "hg" {
				for _, a := range body.Actions {
					if a.LastBuiltRevision.SHA1 != "" {
						build.CommitSha = a.LastBuiltRevision.SHA1
					}
					if a.MercurialRevisionNumber != "" {
						build.CommitSha = a.MercurialRevisionNumber
					}
				}
			} else if vcs == "svn" {
				build.CommitSha = strconv.Itoa(body.ChangeSet.Revisions[0].Revision)
			}

			results = append(results, build)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
