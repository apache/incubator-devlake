package models

import (
	"github.com/merico-dev/lake/models"
)

type GitlabMergeRequestNote struct {
	GitlabId        int `gorm:"primary_key"`
	NoteableId      int
	NoteableIid     int
	NoteableType    string
	AuthorUsername  string
	Body            string
	GitlabCreatedAt string
	Confidential    bool

	models.NoPKModel
}
