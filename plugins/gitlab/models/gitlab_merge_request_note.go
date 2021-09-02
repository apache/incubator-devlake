package models

import (
	"github.com/merico-dev/lake/models"
)

type GitlabMergeRequestNote struct {
	GitlabId        int `gorm:"primary_key"`
	NoteableId      int
	MergeRequestId  int
	MergeRequest    GitlabMergeRequest `gorm:"foreignKey:MergeRequestId;references:Iid"`
	NoteableType    string
	AuthorUsername  string
	Body            string
	GitlabCreatedAt string
	Confidential    bool

	models.NoPKModel
}
