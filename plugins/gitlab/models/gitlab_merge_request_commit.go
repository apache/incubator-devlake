package models

import (
	"github.com/merico-dev/lake/models"
)

type GitlabMergeRequestCommit struct {
	MergeRequestId int
	CommitId       string
	models.NoPKModel
}
