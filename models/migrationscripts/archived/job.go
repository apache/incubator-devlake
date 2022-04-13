package archived

type Job struct {
	Name string `gorm:"type:char(255)"`
	DomainEntity
}
