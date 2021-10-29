package user

import (
	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type User struct {
	base.DomainEntity
	Name      string
	Email     string
	AvatarUrl string
	Timezone  string
}
