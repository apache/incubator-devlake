package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertApiCommitsMeta = core.SubTaskMeta{
	Name:             "convertApiCommits",
	EntryPoint:       ConvertApiCommits,
	EnabledByDefault: true,
	Description:      "Update domain layer commit according to GitlabCommit",
}

func ConvertApiCommits(taskCtx core.SubTaskContext) error {

	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)
	db := taskCtx.GetDb()

	// select all commits belongs to the project
	cursor, err := db.Table("_tool_gitlab_commits gc").
		Joins(`left join _tool_gitlab_project_commits gpc on (
			gpc.commit_sha = gc.sha
		)`).
		Select("gc.*").
		Where("gpc.gitlab_project_id = ?", data.Options.ProjectId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	// TODO: adopt batch indate operation
	userDidGen := didgen.NewDomainIdGenerator(&models.GitlabUser{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabCommit{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabCommit := inputRow.(*models.GitlabCommit)

			// convert commit
			commit := &code.Commit{}
			commit.Sha = gitlabCommit.Sha
			commit.Message = gitlabCommit.Message
			commit.Additions = gitlabCommit.Additions
			commit.Deletions = gitlabCommit.Deletions
			commit.AuthorId = userDidGen.Generate(gitlabCommit.AuthorEmail)
			commit.AuthorName = gitlabCommit.AuthorName
			commit.AuthorEmail = gitlabCommit.AuthorEmail
			commit.AuthoredDate = gitlabCommit.AuthoredDate
			commit.CommitterName = gitlabCommit.CommitterName
			commit.CommitterEmail = gitlabCommit.CommitterEmail
			commit.CommittedDate = gitlabCommit.CommittedDate
			commit.CommitterId = userDidGen.Generate(gitlabCommit.AuthorEmail)

			// convert repo / commits relationship
			repoCommit := &code.RepoCommit{
				RepoId:    didgen.NewDomainIdGenerator(&models.GitlabProject{}).Generate(data.Options.ProjectId),
				CommitSha: gitlabCommit.Sha,
			}

			return []interface{}{
				commit,
				repoCommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
