package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

func ConvertNotes() error {
	var gitlabMergeRequestNotes []gitlabModels.GitlabMergeRequestNote
	err := lakeModels.Db.Find(&gitlabMergeRequestNotes).Error
	if err != nil {
		return err
	}
	domainIdGeneratorNote := didgen.NewDomainIdGenerator(gitlabModels.GitlabMergeRequestNote{})
	for _, note := range gitlabMergeRequestNotes {
		domainNote := convertToNoteModel(&note, domainIdGeneratorNote)
		err := lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(domainNote).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func convertToNoteModel(note *gitlabModels.GitlabMergeRequestNote, domainIdGeneratorNote *didgen.DomainIdGenerator) *code.Note {
	domainNote := &code.Note{
		DomainEntity: domainlayer.DomainEntity{
			Id: domainIdGeneratorNote.Generate(note.GitlabId),
		},
		PrId:        uint64(note.MergeRequestId),
		Type:        note.NoteableType,
		Author:      note.AuthorUsername,
		Body:        note.Body,
		Resolvable:  note.Resolvable,
		System:      note.System,
		CreatedDate: note.GitlabCreatedAt,
	}
	return domainNote
}
