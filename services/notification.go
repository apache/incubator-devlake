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
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/models"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NotificationService FIXME ...
type NotificationService struct {
	EndPoint string
	Secret   string
}

// NewNotificationService FIXME ...
func NewNotificationService(endpoint, secret string) *NotificationService {
	return &NotificationService{
		EndPoint: endpoint,
		Secret:   secret,
	}
}

// PipelineNotification FIXME ...
type PipelineNotification struct {
	PipelineID uint64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	BeganAt    *time.Time
	FinishedAt *time.Time
	Status     string
}

// PipelineStatusChanged FIXME ...
func (n *NotificationService) PipelineStatusChanged(params PipelineNotification) error {
	return n.sendNotification(models.NotificationPipelineStatusChanged, params)
}

func (n *NotificationService) sendNotification(notificationType models.NotificationType, data interface{}) error {
	var dataJson, err = json.Marshal(data)
	if err != nil {
		return err
	}

	var notification models.Notification
	notification.Data = string(dataJson)
	notification.Type = notificationType
	notification.Endpoint = n.EndPoint
	nonce := randSeq(16)
	notification.Nonce = nonce

	err = db.Save(&notification).Error
	if err != nil {
		return err
	}

	sign := n.signature(notification.Data, fmt.Sprintf("%d-%s", notification.ID, nonce))
	url := fmt.Sprintf("%s?nouce=%d-%s&sign=%s", n.EndPoint, notification.ID, nonce, sign)

	resp, err := http.Post(url, "application/json", strings.NewReader(notification.Data))
	if err != nil {
		return err
	}

	notification.ResponseCode = resp.StatusCode
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	notification.Response = string(respBody)
	return db.Save(&notification).Error
}

func (n *NotificationService) signature(input, nouce string) string {
	sum := sha256.Sum256([]byte(input + n.Secret + nouce))
	return hex.EncodeToString(sum[:])
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
