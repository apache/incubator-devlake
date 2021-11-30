package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertNotes() error {
	var githubPullRequestComments []githubModels.GithubPullRequestComment
	err := lakeModels.Db.Find(&githubPullRequestComments).Error
	if err != nil {
		return err
	}
	for _, note := range githubPullRequestComments {
		domainNote := convertToNoteModel(&note)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainNote).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToNoteModel(note *githubModels.GithubPullRequestComment) *code.Note {
	domainNote := &code.Note{
		DomainEntity: domainlayer.DomainEntity{
			OriginKey: okgen.NewOriginKeyGenerator(note).Generate(note.GithubId),
		},
		PrId:        uint64(note.PullRequestId),
		Author:      note.AuthorUsername,
		Body:        note.Body,
		CreatedDate: note.GithubCreatedAt,
	}
	return domainNote
}
