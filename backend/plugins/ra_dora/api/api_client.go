package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/ra_dora/models"
)

type ArgoApiClient struct {
	Client  *api.ApiClient
	BaseUrl string
	Token   string
}

func NewApiClient(connection *models.ArgoConnection) (*ArgoApiClient, errors.Error) {
	client, err := api.NewApiClient(context.Background(), connection.Endpoint, nil, 0, connection.Token, nil)
	if err != nil {
		return nil, errors.Convert(err)
	}

	return &ArgoApiClient{
		Client:  client,
		BaseUrl: connection.Endpoint,
		Token:   connection.Token,
	}, nil
}

func (c *ArgoApiClient) Get(path string, query map[string]string, headers http.Header) (*http.Response, errors.Error) {
	url := fmt.Sprintf("%s/%s", c.BaseUrl, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Convert(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	for key, value := range headers {
		req.Header[key] = value
	}

	resp, err := c.Client.Do("GET", url, nil, nil, req.Header)
	if err != nil {
		return nil, errors.Convert(err)
	}

	return resp, nil
}
