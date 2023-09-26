package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

// RemoteScopes list all available scopes (services) for this connection
// @Summary list all available scopes (services) for this connection
// @Description list all available scopes (services) for this connection
// @Tags plugins/pagerduty
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	connection := &models.ClickupConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}

	// create api client
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, connection)
	if err != nil {
		return nil, err
	}

	var res *http.Response
	path := fmt.Sprintf("/v2/team/%s/space", connection.TeamId)
	res, err = apiClient.Get(path, nil, nil)
	if res.StatusCode != 200 {
		return nil, errors.Default.New(fmt.Sprintf("Got %d code for %s", res.StatusCode, path))
	}

	if err != nil {
		return nil, err
	}
	response := &struct{ Spaces []ClickupSpace }{
		Spaces: []ClickupSpace{},
	}
	err = api.UnmarshalResponse(res, response)
	if err != nil {
		return nil, err
	}

	outputBody := &api.RemoteScopesOutput{
		Children:      []api.RemoteScopesChild{},
		NextPageToken: "",
	}
	// append service to output
	for _, service := range response.Spaces {
		child := api.RemoteScopesChild{
			Type:     "scope",
			ParentId: nil,
			Id:       service.Id,
			Name:     service.Name,
			Data: models.ClickUpSpace{
				Id:   service.Id,
				Name: service.Name,
			},
		}
		outputBody.Children = append(outputBody.Children, child)
	}

	return &plugin.ApiResourceOutput{Body: outputBody, Status: http.StatusOK}, nil
}

// GetScopeList get clickUp spaces
// @Summary get clickUp boards
// @Description get clickUp boards
// @Tags plugins/clickUp
// @Param connectionId path int false "connection ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []ScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/clickUp/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.GetScopeList(input)
}

type ScopeReq api.ScopeReq[models.ClickUpSpace]

// PutScope create or update clickUp board
// @Summary create or update clickUp board
// @Description Create or update clickUp board
// @Tags plugins/clickUp
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body ScopeReq true "json"
// @Success 200  {object} []models.ClickUpSpace
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/clickUp/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Put(input)
}

// GetScope get one Gitlab project
// @Summary get one Gitlab project
// @Description get one Gitlab project
// @Tags plugins/clickup
// @Param connectionId path int false "connection ID"
// @Param scopeId path int false "project ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} ScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/clickup/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.GetScope(input, "id")
}

// UpdateScope patch to clickup project
// @Summary patch to clickup project
// @Description patch to clickup project
// @Tags plugins/clickup
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scopeId path int false "project ID"
// @Param scope body models.GitlabProject true "json"
// @Success 200  {object} models.GitlabProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/clickup/connections/{connectionId}/scopes/{scopeId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Update(input, "id")
}
