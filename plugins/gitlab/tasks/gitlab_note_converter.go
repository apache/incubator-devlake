package tasks

import (
	"context"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

func ConvertNotes(ctx context.Context, projectId int) error {
	gitlabMergeRequestNote := &gitlabModels.GitlabMergeRequestNote{}
	cursor, err := lakeModels.Db.Model(gitlabMergeRequestNote).
		Joins("left join gitlab_merge_requests on gitlab_merge_requests.gitlab_id = gitlab_merge_request_notes.merge_request_id").
		Where("gitlab_merge_requests.project_id = ?", projectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	domainIdGeneratorNote := didgen.NewDomainIdGenerator(&gitlabModels.GitlabMergeRequestNote{})
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, gitlabMergeRequestNote)
		if err != nil {
			return err
		}
		domainNote := convertToNoteModel(gitlabMergeRequestNote, domainIdGeneratorNote)
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
