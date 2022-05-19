package tasks

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiCommitStatsMeta = core.SubTaskMeta{
	Name:             "extractApiCommitStats",
	EntryPoint:       ExtractApiCommitStats,
	EnabledByDefault: false,
	Description:      "Extract raw commit stats data into tool layer table github_commit_stats",
}

type ApiSingleCommitResponse struct {
	Sha   string
	Stats struct {
		Additions int
		Deletions int
	}
	Commit struct {
		Committer struct {
			Name  string
			Email string
			Date  helper.Iso8601Time
		}
	}
}

func ExtractApiCommitStats(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraCommits by Board
			*/
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_COMMIT_STATS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &ApiSingleCommitResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			if body.Sha == "" {
				return nil, nil
			}

			db := taskCtx.GetDb()
			commit := &models.GithubCommit{}
			err = db.Model(commit).Where("sha = ?", body.Sha).Limit(1).Find(commit).Error
			if err != nil {
				return nil, err
			}

			commit.Additions = body.Stats.Additions
			commit.Deletions = body.Stats.Deletions

			commitStat := &models.GithubCommitStat{
				Additions:     body.Stats.Additions,
				Deletions:     body.Stats.Deletions,
				CommittedDate: body.Commit.Committer.Date.ToTime(),
				Sha:           body.Sha,
			}

			results := make([]interface{}, 0, 2)

			results = append(results, commit)
			results = append(results, commitStat)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
