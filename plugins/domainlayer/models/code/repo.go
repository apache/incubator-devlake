package code

import (
	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type Repo struct {
	base.DomainEntity
	Name string `json:"name"`
	Url  string `json:"url"`
}
