package models

import (
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

// ArgoConn holds the essential information to connect to the Argo API
type ArgoConn struct {
	api.RestConnection `mapstructure:",squash"`
	api.AccessToken    `mapstructure:",squash"`
}

const ArgoCloudEndPoint string = "https://argo.example.com/api/v1/"
const ArgoApiClientData_UserId string = "UserId"
const ArgoApiClientData_UserName string = "UserName"
const ArgoApiClientData_ApiVersion string = "ApiVersion"

// this function is used to rewrite the same function of AccessToken
func (conn *ArgoConn) SetupAuthentication(request *http.Request) errors.Error {
	return nil
}

func (conn *ArgoConn) Sanitize() ArgoConn {
	conn.Token = utils.SanitizeString(conn.Token)
	return *conn
}

// PrepareApiClient test api and set the IsPrivateToken,version,UserId and so on.
func (conn *ArgoConn) PrepareApiClient(apiClient plugin.ApiClient) errors.Error {
	header1 := http.Header{}
	header1.Set("Authorization", fmt.Sprintf("Bearer %v", conn.Token))
	// test request for access token
	userResBody := &ApiUserResponse{}
	res, err := apiClient.Get("user", nil, header1)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusUnauthorized {
		err = api.UnmarshalResponse(res, userResBody)
		if err != nil {
			return errors.Convert(err)
		}
		if res.StatusCode == http.StatusUnauthorized {
			return errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error while testing connection")
		}
		if res.StatusCode != http.StatusOK {
			return errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection")
		}
		apiClient.SetHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", conn.Token),
		})
	} else {
		header2 := http.Header{}
		header2.Set("Private-Token", conn.Token)
		res, err = apiClient.Get("user", nil, header2)
		if err != nil {
			return errors.Convert(err)
		}
		err = api.UnmarshalResponse(res, userResBody)
		if err != nil {
			return errors.Convert(err)
		}
		if res.StatusCode == http.StatusUnauthorized {
			return errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error while testing connection[PrivateToken]")
		}
		if res.StatusCode != http.StatusOK {
			return errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection[PrivateToken]")
		}
		apiClient.SetHeaders(map[string]string{
			"Private-Token": conn.Token,
		})
	}
	// get argo version
	versionResBody := &ApiVersionResponse{}
	res, err = apiClient.Get("version", nil, nil)
	if err != nil {
		return errors.Convert(err)
	}

	err = api.UnmarshalResponse(res, versionResBody)
	if err != nil {
		return errors.Convert(err)
	}

	// add v for semver compare
	if versionResBody.Version[0] != 'v' {
		versionResBody.Version = "v" + versionResBody.Version
	}

	apiClient.SetData(ArgoApiClientData_UserId, userResBody.Id)
	apiClient.SetData(ArgoApiClientData_UserName, userResBody.Name)
	apiClient.SetData(ArgoApiClientData_ApiVersion, versionResBody.Version)

	return nil
}

var _ plugin.ApiConnection = (*ArgoConnection)(nil)

// ArgoConnection holds ArgoConn plus ID/Name for database storage
type ArgoConnection struct {
	api.BaseConnection `mapstructure:",squash"`
	ArgoConn           `mapstructure:",squash"`
}

// This object conforms to what the frontend currently expects.
type ArgoResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	ArgoConnection
}

type ApiVersionResponse struct {
	Version  string `json:"version"`
	Revision string `json:"revision"`
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string `json:"name"`
}

func (ArgoConnection) TableName() string {
	return "_tool_argo_connections"
}

func (connection ArgoConnection) Sanitize() ArgoConnection {
	connection.ArgoConn = connection.ArgoConn.Sanitize()
	return connection
}

func (connection *ArgoConnection) MergeFromRequest(target *ArgoConnection, body map[string]interface{}) error {
	token := target.Token
	if err := api.DecodeMapStruct(body, target, true); err != nil {
		return err
	}
	modifiedToken := target.Token
	if modifiedToken == "" || modifiedToken == utils.SanitizeString(token) {
		target.Token = token
	}
	return nil
}
