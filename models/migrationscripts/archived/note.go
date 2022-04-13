package archived

import (
	"time"
)

type Note struct {
	DomainEntity
	PrId        string `gorm:"index;comment:References the pull request for this note;type:varchar(100)"`
	Type        string `gorm:"type:char(100)"`
	Author      string `gorm:"type:char(255)"`
	Body        string
	Resolvable  bool `gorm:"comment:Is or is not a review comment"`
	System      bool `gorm:"comment:Is or is not auto-generated vs. human generated"`
	CreatedDate time.Time
}
