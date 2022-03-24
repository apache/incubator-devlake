package models

import "github.com/merico-dev/lake/models/common"

type AECommit struct {
	HexSha      string `gorm:"primaryKey"`
	AnalysisId  string
	AuthorEmail string
	DevEq       int
	AEProjectId int
	common.NoPKModel
}
