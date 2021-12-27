package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

func ConvertCommits() error {
	var gitlabCommits []gitlabModels.GitlabCommit
	err := lakeModels.Db.Find(&gitlabCommits).Error
	if err != nil {
		return err
	}
	for _, commit := range gitlabCommits {
		domainCommit := convertToCommitModel(&commit)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainCommit).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToCommitModel(commit *gitlabModels.GitlabCommit) *code.Commit {
	domainCommit := &code.Commit{
		Sha:            commit.GitlabId,
		Message:        commit.Message,
		AuthorName:     commit.AuthorName,
		Additions:      commit.Additions,
		Deletions:      commit.Deletions,
		AuthorEmail:    commit.AuthorEmail,
		AuthoredDate:   commit.AuthoredDate,
		CommitterName:  commit.CommitterName,
		CommitterEmail: commit.CommitterEmail,
		CommittedDate:  commit.CommittedDate,
	}
	return domainCommit
}
