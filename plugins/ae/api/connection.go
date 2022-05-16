package api

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/ae/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"net/http"
)

type ApiMeResponse struct {
	Name string `json:"name"`
}

/*
GET /plugins/ae/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// TODO: implement test connection
	return &core.ApiResourceOutput{Body: true}, nil
}

/*
PATCH /plugins/ae/connections/:connectionId
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	v := config.GetConfig()
	connection := &models.AeConnection{}
	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	// update from request and save to .env
	err = helper.DecodeStruct(v, connection, input.Body, "env")
	if err != nil {
		return nil, err
	}
	err = config.WriteConfig(v)
	if err != nil {
		return nil, err
	}
	response := models.AeResponse{
		AeConnection: *connection,
		Name:         "Ae",
		ID:           1,
	}
	return &core.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

/*
GET /plugins/ae/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-connection is developed.
	v := config.GetConfig()
	connection := &models.AeConnection{}

	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	response := models.AeResponse{
		AeConnection: *connection,
		Name:         "Ae",
		ID:           1,
	}

	return &core.ApiResourceOutput{Body: []models.AeResponse{response}}, nil
}

/*
GET /plugins/ae/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-connection is developed.)
	v := config.GetConfig()
	connection := &models.AeConnection{}
	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	response := &models.AeResponse{
		AeConnection: *connection,
		Name:         "Ae",
		ID:           1,
	}
	return &core.ApiResourceOutput{Body: response}, nil
}
