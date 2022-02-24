package tasks

import (
	"context"
	"github.com/merico-dev/lake/errors"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/github/models"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertCommits(ctx context.Context, githubRepoId int) error {
	// select all commits belongs to the repo
	cursor, err := lakeModels.Db.Table("github_commits gc").
		Joins(`left join github_repo_commits grc on (
			grc.commit_sha = gc.sha
		)`).
		Select("gc.*").
		Where("grc.repo_id = ?", githubRepoId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	userDidGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})
	repoCommit := &code.RepoCommit{
		RepoId: didgen.NewDomainIdGenerator(&githubModels.GithubRepo{}).Generate(githubRepoId),
	}
	githubCommit := &models.GithubCommit{}
	commit := &code.Commit{}
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubCommit)
		if err != nil {
			return err
		}
		// convert commit
		commit.Sha = githubCommit.Sha
		commit.Message = githubCommit.Message
		commit.Additions = githubCommit.Additions
		commit.Deletions = githubCommit.Deletions
		commit.AuthorId = userDidGen.Generate(githubCommit.AuthorId)
		commit.AuthorName = githubCommit.AuthorName
		commit.AuthorEmail = githubCommit.AuthorEmail
		commit.AuthoredDate = githubCommit.AuthoredDate
		commit.CommitterName = githubCommit.CommitterName
		commit.CommitterEmail = githubCommit.CommitterEmail
		commit.CommittedDate = githubCommit.CommittedDate
		commit.CommitterId = userDidGen.Generate(githubCommit.CommitterId)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(commit).Error
		if err != nil {
			return err
		}
		// convert repo / commits relationship
		repoCommit.CommitSha = githubCommit.Sha
		err = lakeModels.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(repoCommit).Error
		if err != nil {
			return err
		}
	}
	return nil
}
