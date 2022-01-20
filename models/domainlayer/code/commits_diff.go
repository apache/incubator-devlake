package code

type CommitsDiff struct {
	NewCommitSha string `gorm:"primaryKey;type:char(40)"`
	OldCommitSha string `gorm:"primaryKey;type:char(40)"`
	CommitSha    string `gorm:"primaryKey;type:char(40)"`
	SortingIndex int
}
