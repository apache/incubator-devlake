package code

type CommitParent struct {
	CommitSha       string `json:"commitSha" gorm:"primaryKey;type:char(40);comment:commit hash"`
	ParentCommitSha string `json:"parentCommitSha" gorm:"primaryKey;type:char(40);comment:parent commit hash"`
}
