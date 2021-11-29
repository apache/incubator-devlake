package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/merico-dev/lake/models"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type NotificationService struct {
	EndPoint string
	Secret   string
}

func NewNotificationService(endpoint, secret string) *NotificationService {
	return &NotificationService{
		EndPoint: endpoint,
		Secret:   secret,
	}
}

type PipelineNotification struct {
	PipelineID uint64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	BeganAt    *time.Time
	FinishedAt *time.Time
	Status     string
}

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

	err = models.Db.Save(&notification).Error
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
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	notification.Response = string(respBody)
	return models.Db.Save(&notification).Error
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
