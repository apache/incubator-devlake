package code

type CommitParent struct {
	CommitSha       string `json:"commitSha" gorm:"primaryKey;type:varchar(40);comment:commit hash"`
	ParentCommitSha string `json:"parentCommitSha" gorm:"primaryKey;type:varchar(40);comment:parent commit hash"`
}
