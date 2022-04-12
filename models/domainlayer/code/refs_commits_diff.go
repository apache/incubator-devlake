package code

type RefsCommitsDiff struct {
	NewRefId        string `gorm:"primaryKey;type:varchar(255)"`
	OldRefId        string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha       string `gorm:"primaryKey;type:char(40)"`
	NewRefCommitSha string `gorm:"type:char(40)"`
	OldRefCommitSha string `gorm:"type:char(40)"`
	SortingIndex    int
}
