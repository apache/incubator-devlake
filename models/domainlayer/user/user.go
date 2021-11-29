package user

import (
	"github.com/merico-dev/lake/models/domainlayer/base"
)

type User struct {
	base.DomainEntity
	Name      string
	Email     string
	AvatarUrl string
	Timezone  string
}
