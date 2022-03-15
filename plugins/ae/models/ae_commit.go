package models

import "github.com/merico-dev/lake/plugins/helper"

type AECommit struct {
	HexSha      string `gorm:"primaryKey"`
	AnalysisId  string
	AuthorEmail string
	DevEq       int
	AEProjectId int

	helper.RawDataOrigin
}
