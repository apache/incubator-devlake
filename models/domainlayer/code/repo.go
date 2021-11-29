package code

import (
	"github.com/merico-dev/lake/models/domainlayer/base"
)

type Repo struct {
	base.DomainEntity
	Name string `json:"name"`
	Url  string `json:"url"`
}
