package crossdomain

type BoardRepo struct {
	BoardId string `gorm:"primaryKey;type:varchar(255)"`
	RepoId  string `gorm:"primaryKey;type:varchar(255)"`
}
