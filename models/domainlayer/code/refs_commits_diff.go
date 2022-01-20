package code

type RefsCommitsDiff struct {
	NewRefName   string `gorm:"primaryKey;type:varchar(255)"`
	OldRefName   string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha    string `gorm:"primaryKey;type:char(40)"`
	NewCommitSha string `gorm:"type:char(40)"`
	OldCommitSha string `gorm:"type:char(40)"`
	SortingIndex int
}
