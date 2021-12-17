package factory

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/user"
)

func CreateUser() (*user.User, error) {
	user := &user.User{
		DomainEntity: domainlayer.DomainEntity{
			Id: "1",
		},
		Name:      "",
		Email:     "",
		AvatarUrl: "",
		Timezone:  "",
	}
	return user, nil
}
