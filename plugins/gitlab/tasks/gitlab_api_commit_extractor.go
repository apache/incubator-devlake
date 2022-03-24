package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ExtractCommitsMeta = core.SubTaskMeta{
	Name:             "extractApiCommits",
	EntryPoint:       ExtractApiCommits,
	EnabledByDefault: true,
	Description:      "Extract raw commit data into tool layer table GitlabCommit,GitlabUser and GitlabProjectCommit",
}

func ExtractApiCommits(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
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

// Convert the API response to our DB model instance
func ConvertCommit(commit *GitlabApiCommit) (*models.GitlabCommit, error) {
	gitlabCommit := &models.GitlabCommit{
		Sha:            commit.GitlabId,
		Title:          commit.Title,
		Message:        commit.Message,
		ShortId:        commit.ShortId,
		AuthorName:     commit.AuthorName,
		AuthorEmail:    commit.AuthorEmail,
		AuthoredDate:   commit.AuthoredDate.ToTime(),
		CommitterName:  commit.CommitterName,
		CommitterEmail: commit.CommitterEmail,
		CommittedDate:  commit.CommittedDate.ToTime(),
		WebUrl:         commit.WebUrl,
		Additions:      commit.Stats.Additions,
		Deletions:      commit.Stats.Deletions,
		Total:          commit.Stats.Total,
	}
	return gitlabCommit, nil
}
