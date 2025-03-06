package models

import (
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

type ArgoConnection struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// TableName implements dal.Tabler.
func (conn *ArgoConnection) TableName() string {
	panic("_tool_argo_connection")
}

func (conn *ArgoConnection) PrepareApiClient(apiClient plugin.ApiClient) errors.Error {
	header := http.Header{}
	header.Set("Authorization", "Bearer "+conn.Token)

	return nil
}
