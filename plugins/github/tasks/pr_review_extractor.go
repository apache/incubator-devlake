package tasks

import (
	"encoding/json"
	"strings"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiPullRequestReviewsMeta = core.SubTaskMeta{
	Name:             "extractApiPullRequestReviews",
	EntryPoint:       ExtractApiPullRequestReviews,
	EnabledByDefault: true,
	Description:      "Extract raw PullRequestReviews data into tool layer table github_reviewers",
}

type PullRequestReview struct {
	GithubId int `json:"id"`
	User     struct {
		Id    int
		Login string
	}
	Body        string
	State       string
	SubmittedAt helper.Iso8601Time `json:"submitted_at"`
}

func ExtractApiPullRequestReviews(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_PULL_REQUEST_REVIEW_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			apiPullRequestReview := &PullRequestReview{}
			if strings.HasPrefix(string(row.Data), "{\"message\": \"Not Found\"") {
				return nil, nil
			}
			err := json.Unmarshal(row.Data, apiPullRequestReview)
			if err != nil {
				return nil, err
			}
			pull := &SimplePr{}
			err = json.Unmarshal(row.Input, pull)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 1)

			githubReviewer := &models.GithubReviewer{
				GithubId:      apiPullRequestReview.User.Id,
				Login:         apiPullRequestReview.User.Login,
				PullRequestId: pull.GithubId,
			}
			results = append(results, githubReviewer)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
