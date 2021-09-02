package models

type JiraBoardGitlabProject struct {
	JiraBoardId     uint64 `gorm:"primaryKey"`
	GitlabProjectId uint64 `gorm:"primaryKey"`
}
