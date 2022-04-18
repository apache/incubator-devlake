package archived

type Job struct {
	Name string `gorm:"type:varchar(255)"`
	DomainEntity
}
