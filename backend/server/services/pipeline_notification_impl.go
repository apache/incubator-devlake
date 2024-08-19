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

package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/utils"
	"io"
	"net/http"
	"strings"
)

// DefaultPipelineNotificationService FIXME ...
type DefaultPipelineNotificationService struct {
	EndPoint string
	Secret   string
}

// NewDefaultPipelineNotificationService creates a new DefaultPipelineNotificationService
func NewDefaultPipelineNotificationService(endpoint, secret string) *DefaultPipelineNotificationService {
	return &DefaultPipelineNotificationService{
		EndPoint: endpoint,
		Secret:   secret,
	}
}

// PipelineStatusChanged FIXME ...
func (n *DefaultPipelineNotificationService) PipelineStatusChanged(params PipelineNotificationParam) errors.Error {
	return n.sendNotification(models.NotificationPipelineStatusChanged, params)
}

func (n *DefaultPipelineNotificationService) sendNotification(notificationType models.NotificationType, data interface{}) errors.Error {
	var dataJson, err = json.Marshal(data)
	if err != nil {
		return errors.Convert(err)
	}

	var notification models.Notification
	notification.Data = string(dataJson)
	notification.Type = notificationType
	notification.Endpoint = n.EndPoint
	nonce, err1 := utils.RandLetterBytes(16)
	if err1 != nil {
		return err1
	}
	notification.Nonce = nonce

	err = db.Create(&notification)
	if err != nil {
		return errors.Convert(err)
	}

	sign := n.signature(notification.Data, fmt.Sprintf("%d-%s", notification.ID, nonce))
	url := fmt.Sprintf("%s?nouce=%d-%s&sign=%s", n.EndPoint, notification.ID, nonce, sign)

	resp, err := http.Post(url, "application/json", strings.NewReader(notification.Data))
	if err != nil {
		return errors.Convert(err)
	}

	notification.ResponseCode = resp.StatusCode
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Convert(err)
	}
	notification.Response = string(respBody)
	return db.Update(notification)
}

func (n *DefaultPipelineNotificationService) signature(input, nouce string) string {
	sum := sha256.Sum256([]byte(input + n.Secret + nouce))
	return hex.EncodeToString(sum[:])
}
