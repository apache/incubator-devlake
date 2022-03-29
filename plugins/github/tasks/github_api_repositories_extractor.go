package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ExtractApiRepoMeta = core.SubTaskMeta{
	Name:        "extractApiRepo",
	EntryPoint:  ExtractApiRepositories,
	Required:    true,
	Description: "Extract raw Repositories data into tool layer table github_repos",
}

type ApiRepoResponse GithubApiRepo

type GithubApiRepo struct {
	Name        string `json:"name"`
	GithubId    int    `json:"id"`
	HTMLUrl     string `json:"html_url"`
	Language    string `json:"language"`
	Description string `json:"description"`
	Owner       models.GithubUser
	Parent      *GithubApiRepo    `json:"parent"`
	CreatedAt   core.Iso8601Time  `json:"created_at"`
	UpdatedAt   *core.Iso8601Time `json:"updated_at"`
}

func ExtractApiRepositories(taskCtx core.SubTaskContext) error {
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
			Table: RAW_REPOSITORIES_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &ApiRepoResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 1)
			githubRepository := &models.GithubRepo{
				GithubId:    body.GithubId,
				Name:        body.Name,
				HTMLUrl:     body.HTMLUrl,
				Description: body.Description,
				OwnerId:     body.Owner.Id,
				OwnerLogin:  body.Owner.Login,
				Language:    body.Language,
				CreatedDate: body.CreatedAt.ToTime(),
				UpdatedDate: core.Iso8601TimeToTime(body.UpdatedAt),
			}
			data.Repo = githubRepository

			if body.Parent != nil {
				githubRepository.ParentGithubId = body.Parent.GithubId
				githubRepository.ParentHTMLUrl = body.Parent.HTMLUrl
			}
			results = append(results, githubRepository)
			taskCtx.TaskContext().GetData().(*GithubTaskData).Repo = githubRepository
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
