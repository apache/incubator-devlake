package models

type AECommit struct {
	HexSha      string `gorm:"primaryKey"`
	AnalysisId  string
	AuthorEmail string
	DevEq       int
	AEProjectId int
}
