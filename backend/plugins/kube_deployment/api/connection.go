/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"

	"github.com/apache/incubator-devlake/core/plugin"
	kubeDeploymentHelper "github.com/apache/incubator-devlake/plugins/kube_deployment/helper"
	"github.com/apache/incubator-devlake/plugins/kube_deployment/models"
	"github.com/apache/incubator-devlake/server/api/shared"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubeDeploymentTestConnResponse struct {
	shared.ApiBody
	Connection *models.KubeConn
}

type ReturnObject struct {
	models.KubeConnection `mapstructure:",squash"`
	Credentials           map[string]interface{} `mapstructure:"credentials" json:"credentials"`
}

// TODO Please modify the following code to fit your needs
// @Summary test kube_deployment connection
// @Description Test kube_deployment Connection. endpoint: "https://dev.kube_deployment.com/{organization}/
// @Tags plugins/kube_deployment
// @Param body body models.KubeConn true "json body"
// @Success 200  {object} KubeDeploymentTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// TODO: Modify this
	body := KubeDeploymentTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = nil
	// output
	return &plugin.ApiResourceOutput{Body: body, Status: 200}, nil
	// // decode
	// var err errors.Error
	// var connection models.KubeConn
	// if err = helper.Decode(input.Body, &connection, vld); err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("connection endpoint: %v\n", connection.Endpoint)
	// // test connection
	// apiClient, err := helper.NewApiClient(
	// 	context.TODO(),
	// 	connection.Endpoint,
	// 	map[string]string{
	// 		// "Authorization": fmt.Sprintf("Bearer %v", connection.Token),
	// 	},
	// 	3*time.Second,
	// 	connection.Proxy,
	// 	basicRes,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// res, err := apiClient.Get("", nil, nil)
	// if err != nil {
	// 	return nil, err
	// }
	// // resBody := &models.ApiUserResponse{}
	// // err = helper.UnmarshalResponse(res, resBody)
	// // if err != nil {
	// // 	return nil, err
	// // }

	// if res.StatusCode != http.StatusOK {
	// 	return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code: %d", res.StatusCode))
	// }
	// body := KubeDeploymentTestConnResponse{}
	// body.Success = true
	// body.Message = "success"
	// body.Connection = &connection
	// // output
	// return &plugin.ApiResourceOutput{Body: body, Status: 200}, nil
}

// TODO Please modify the folowing code to adapt to your plugin
// @Summary create kube_deployment connection
// @Description Create kube_deployment connection
// @Tags plugins/kube_deployment
// @Param body body models.KubeConnection true "json body"
// @Success 200  {object} models.KubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.KubeConnection{}

	credentials := input.Body["credentials"].(map[string]interface{})
	token, _ := json.Marshal(credentials)

	if input.Body == nil {
		return nil, errors.BadInput.New("missing connectionId")
	}
	fmt.Println("token: -->", string(token))
	input.Body["credentials"] = string(token)
	fmt.Println("input.Body[credentials]: -->", input.Body["credentials"])
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}

	returnObject := ReturnObject{
		KubeConnection: *connection,
		Credentials:    credentials,
	}

	return &plugin.ApiResourceOutput{Body: returnObject, Status: http.StatusOK}, nil
}

// TODO Please modify the folowing code to adapt to your plugin
// @Summary patch kube_deployment connection
// @Description Patch kube_deployment connection
// @Tags plugins/kube_deployment
// @Param body body models.KubeConnection true "json body"
// @Success 200  {object} models.KubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KubeConnection{}

	credentialsInput := input.Body["credentials"].(map[string]interface{})
	token, _ := json.Marshal(credentialsInput)

	if input.Body == nil {
		return nil, errors.BadInput.New("missing connectionId")
	}

	input.Body["credentials"] = string(token)

	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}

	var credentials map[string]interface{}
	json.Unmarshal([]byte(connection.Credentials), &credentials)

	returnObject := ReturnObject{
		KubeConnection: *connection,
		Credentials:    credentials,
	}

	return &plugin.ApiResourceOutput{Body: returnObject}, nil
}

// @Summary delete a kube_deployment connection
// @Description Delete a kube_deployment connection
// @Tags plugins/kube_deployment
// @Success 200  {object} models.KubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KubeConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

// @Summary get all kube_deployment connections
// @Description Get all kube_deployment connections
// @Tags plugins/kube_deployment
// @Success 200  {object} []models.KubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.KubeConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	returnObjects := make([]ReturnObject, len(connections))
	for i, connection := range connections {
		var credentials map[string]interface{}
		json.Unmarshal([]byte(connection.Credentials), &credentials)

		returnObjects[i].KubeConnection = connection
		returnObjects[i].Credentials = credentials
	}

	return &plugin.ApiResourceOutput{Body: returnObjects, Status: http.StatusOK}, nil
}

// TODO Please modify the folowing code to adapt to your plugin
// @Summary get kube_deployment connection detail
// @Description Get kube_deployment connection detail
// @Tags plugins/kube_deployment
// @Success 200  {object} models.KubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KubeConnection{}
	err := connectionHelper.First(connection, input.Params)

	var credentials map[string]interface{}
	json.Unmarshal([]byte(connection.Credentials), &credentials)

	returnObject := ReturnObject{
		KubeConnection: *connection,
		Credentials:    credentials,
	}
	return &plugin.ApiResourceOutput{Body: returnObject}, err
}

// @Summary Get kubernetes namespaces
// @Description Get kubernetes namespaces
// @Tags plugins/kube_deployment
// @Success 200  {object} models.KubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections/{connectionId}/namespaces [GET]
func GetNameSpaces(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	fmt.Println("input.Params: -->", input.Params)
	connection := &models.KubeConnection{}
	err := connectionHelper.First(connection, input.Params)

	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get KubeDeployment API client instance")
	}

	var credentials map[string]interface{}
	json.Unmarshal([]byte(connection.Credentials), &credentials)
	KubeAPIClient := kubeDeploymentHelper.NewKubeApiClient(credentials)

	// Get the list of namespaces
	namespaces, _ := KubeAPIClient.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Failed to get namespaces: %v", err)
	}

	var namespaceList = make([]string, len(namespaces.Items))

	for i, namespace := range namespaces.Items {
		fmt.Println(namespace.Name)
		namespaceList[i] = namespace.Name
	}

	return &plugin.ApiResourceOutput{Body: namespaceList, Status: http.StatusOK}, nil
}

// @Summary Get kubernetes namespaces
// @Description Get kubernetes namespaces
// @Tags plugins/kube_deployment
// @Success 200  {object} models.KubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections/{connectionId}/{namespace}/deployments [GET]
func GetDeployments(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	fmt.Println("input.Params: -->", input.Params)

	connectionParam := make(map[string]string)
	connectionParam["connectionId"] = input.Params["connectionId"]

	connection := &models.KubeConnection{}
	err := connectionHelper.First(connection, connectionParam)

	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get KubeDeployment API client instance")
	}

	var credentials map[string]interface{}
	json.Unmarshal([]byte(connection.Credentials), &credentials)
	KubeAPIClient := kubeDeploymentHelper.NewKubeApiClient(credentials)

	namespace := input.Params["namespace"]

	// Get the list of namespaces
	deployments, _ := KubeAPIClient.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get deployments")
	}

	var deploymentList = make([]string, len(deployments.Items))

	for i, deployment := range deployments.Items {
		fmt.Println(deployment.Name)
		deploymentList[i] = deployment.Name
	}

	return &plugin.ApiResourceOutput{Body: deploymentList, Status: http.StatusOK}, nil
}
