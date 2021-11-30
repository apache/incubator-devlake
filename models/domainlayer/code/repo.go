package code

import (
	"github.com/merico-dev/lake/models/domainlayer"
)

type Repo struct {
	domainlayer.DomainEntity
	Name string `json:"name"`
	Url  string `json:"url"`
}
