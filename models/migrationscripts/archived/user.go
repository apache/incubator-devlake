package archived

type User struct {
	DomainEntity
	Name      string `gorm:"type:varchar(255)"`
	Email     string `gorm:"type:varchar(255)"`
	AvatarUrl string `gorm:"type:varchar(255)"`
	Timezone  string `gorm:"type:varchar(255)"`
}
