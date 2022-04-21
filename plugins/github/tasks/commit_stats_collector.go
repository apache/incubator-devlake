package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/helper"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
)

const RAW_COMMIT_STATS_TABLE = "github_api_commit_stats"

var CollectApiCommitStatsMeta = core.SubTaskMeta{
	Name:             "collectApiCommitStats",
	EntryPoint:       CollectApiCommitStats,
	EnabledByDefault: false,
	Description:      "Collect commitStats data from Github api",
}

func CollectApiCommitStats(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)

	var latestUpdated models.GithubCommitStat
	err := db.Model(&latestUpdated).Joins("left join _tool_github_repo_commits on _tool_github_commit_stats.sha = _tool_github_repo_commits.commit_sha").
		Where("_tool_github_repo_commits.repo_id = ?", data.Repo.GithubId).
		Order("committed_date DESC").Limit(1).Find(&latestUpdated).Error
	if err != nil {
		return fmt.Errorf("failed to get latest github commit record: %w", err)
	}

	cursor, err := db.Model(&models.GithubCommit{}).
		Joins("left join _tool_github_repo_commits on _tool_github_commits.sha = _tool_github_repo_commits.commit_sha").
		Where("_tool_github_repo_commits.repo_id = ? and _tool_github_commits.committed_date >= ?",
			data.Repo.GithubId, latestUpdated.CommittedDate.String()).
		Rows()
	if err != nil {
		return err
	}
	iterator, err := helper.NewCursorIterator(db, cursor, reflect.TypeOf(models.GithubCommit{}))
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
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
		ApiClient: data.ApiClient,
		PageSize:  100,
		Input:     iterator,
		/*
			url may use arbitrary variables from different source in any order, we need GoTemplate to allow more
			flexible for all kinds of possibility.
			Pager contains information for a particular page, calculated by ApiCollector, and will be passed into
			GoTemplate to generate a url for that page.
			We want to do page-fetching in ApiCollector, because the logic are highly similar, by doing so, we can
			avoid duplicate logic for every tasks, and when we have a better idea like improving performance, we can
			do it in one place
		*/
		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/commits/{{ .Input.Sha }}",
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("state", "all")
			query.Set("direction", "asc")
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},

		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return nil, err
			}
			return []json.RawMessage{body}, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
