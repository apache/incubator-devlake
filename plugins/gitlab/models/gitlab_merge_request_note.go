package models

import (
	"github.com/merico-dev/lake/models"
)

type GitlabMergeRequestNote struct {
	GitlabId        int `gorm:"primary_key"`
	NoteableId      int
	MergeRequestId  int `gorm:"index"`
	NoteableType    string
	AuthorUsername  string
	Body            string
	GitlabCreatedAt string
	Confidential    bool
	Resolvable      bool // Resolvable means a comment is a code review comment
	System          bool // System means the comment is generated automatically

	models.NoPKModel
}
