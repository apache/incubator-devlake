package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertCommits() error {
	var githubCommits []githubModels.GithubCommit
	err := lakeModels.Db.Find(&githubCommits).Error
	if err != nil {
		return err
	}
	userDidGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})
	for _, commit := range githubCommits {
		domainCommit := convertToCommitModel(&commit)
		domainCommit.AuthorId = userDidGen.Generate(commit.AuthorId)
		domainCommit.CommiterId = userDidGen.Generate(commit.CommitterId)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainCommit).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func convertToCommitModel(commit *githubModels.GithubCommit) *code.Commit {
	domainCommit := &code.Commit{
		Sha:            commit.Sha,
		Message:        commit.Message,
		Additions:      commit.Additions,
		Deletions:      commit.Deletions,
		AuthorName:     commit.AuthorName,
		AuthorEmail:    commit.AuthorEmail,
		AuthoredDate:   commit.AuthoredDate,
		CommitterName:  commit.CommitterName,
		CommitterEmail: commit.CommitterEmail,
		CommittedDate:  commit.CommittedDate,
	}
	return domainCommit
}
