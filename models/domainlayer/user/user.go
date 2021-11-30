package user

import (
	"github.com/merico-dev/lake/models/domainlayer"
)

type User struct {
	domainlayer.DomainEntity
	Name      string
	Email     string
	AvatarUrl string
	Timezone  string
}
