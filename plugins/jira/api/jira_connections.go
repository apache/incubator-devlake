package api

import (
	"github.com/merico-dev/lake/plugins/core"
)

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// Get the source ID from the input
	// Get the source from the DB
	// Use endpoint and auth to make a request to JIRA

	return &core.ApiResourceOutput{Body: true}, nil
}
