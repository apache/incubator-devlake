package models

import (
	"github.com/merico-dev/lake/models"
)

type AECommit struct {
	HexSha      string `gorm:"primaryKey"`
	AnalysisId  string
	AuthorEmail string
	DevEq       int

	models.NoPKModel
}
