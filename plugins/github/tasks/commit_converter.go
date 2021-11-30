package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertCommits() error {
	var githubCommits []githubModels.GithubCommit
	err := lakeModels.Db.Find(&githubCommits).Error
	if err != nil {
		return err
	}
	for _, commit := range githubCommits {
		domainCommit := convertToCommitModel(&commit)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainCommit).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToCommitModel(commit *githubModels.GithubCommit) *code.Commit {
	domainCommit := &code.Commit{
		DomainEntity: domainlayer.DomainEntity{
			OriginKey: okgen.NewOriginKeyGenerator(commit).Generate(commit.Sha),
		},
		Sha:            commit.Sha,
		RepoId:         uint64(commit.RepositoryId),
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
