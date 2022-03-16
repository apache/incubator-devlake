package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

func ExtractApiCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GitlabTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,

			Params: GitlabApiParams{
				ProjectId: data.Options.ProjectId,
			},
			Table: RAW_COMMIT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			// need to extract 3 kinds of entities here
			results := make([]interface{}, 0, 3)

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

			// create project/commits relationship
			gitlabProjectCommit := &models.GitlabProjectCommit{GitlabProjectId: data.Options.ProjectId}
			gitlabProjectCommit.CommitSha = gitlabCommit.Sha

			// create gitlab user
			gitlabUserAuthor := &models.GitlabUser{}
			gitlabUserAuthor.Email = gitlabCommit.AuthorEmail
			gitlabUserAuthor.Name = gitlabCommit.AuthorName

			results = append(results, gitlabCommit)
			results = append(results, gitlabProjectCommit)
			results = append(results, gitlabUserAuthor)

			// For Commiter Email is not same as AuthorEmail
			if gitlabCommit.CommitterEmail != gitlabUserAuthor.Email {
				gitlabUserCommitter := &models.GitlabUser{}
				gitlabUserCommitter.Email = gitlabCommit.CommitterEmail
				gitlabUserCommitter.Name = gitlabCommit.CommitterName
				results = append(results, gitlabUserCommitter)
			}

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
