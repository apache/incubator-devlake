package tasks

import (
	"context"
	"github.com/merico-dev/lake/errors"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertNotes(ctx context.Context) error {
	githubPrComment := &githubModels.GithubPullRequestComment{}
	cursor, err := lakeModels.Db.Model(githubPrComment).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	domainNoteIdGenerator := didgen.NewDomainIdGenerator(githubPrComment)
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubPrComment)
		if err != nil {
			return err
		}
		domainNote := convertToNoteModel(githubPrComment, domainNoteIdGenerator)
		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainNote).Error
		if err != nil {
			return err
		}

	}
	return nil
}
func convertToNoteModel(note *githubModels.GithubPullRequestComment, didGenerator *didgen.DomainIdGenerator) *code.Note {
	domainNote := &code.Note{
		DomainEntity: domainlayer.DomainEntity{
			Id: didGenerator.Generate(note.GithubId),
		},
		PrId:        uint64(note.PullRequestId),
		Author:      note.AuthorUsername,
		Body:        note.Body,
		CreatedDate: note.GithubCreatedAt,
	}
	return domainNote
}
