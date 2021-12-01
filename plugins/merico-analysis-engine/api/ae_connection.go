package api

import (
	"github.com/merico-dev/lake/plugins/core"
)

type ApiMeResponse struct {
	Name     string `json:"name"`
	GithubId int    `json:"id"`
	HTMLUrl  string `json:"html_url"`
}

/*
GET /plugins/github/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// TODO: implement test connection
	return &core.ApiResourceOutput{Body: true}, nil
}
