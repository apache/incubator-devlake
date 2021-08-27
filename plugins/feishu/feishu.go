package main

import (
	"encoding/json"
	"fmt"
	"github.com/faabiosr/cachego/file"
	"github.com/fastwego/feishu"
	"github.com/merico-dev/lake/plugins/feishu/models"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Feishu string

func (plugin Feishu) Description() string {
	return "To collect and enrich data from Feishu"
}

func (plugin Feishu) Execute(options map[string]interface{}, progress chan<- float32) {
	fmt.Println("start feishu plugin execution")

	// 内部应用 tenant_access_token 管理器
	Atm := &feishu.DefaultAccessTokenManager{
		Id:    `cli_a074eb7697f8d00b`,
		Cache: file.New(os.TempDir()),
		GetRefreshRequestFunc: func() *http.Request {
			payload := `{
            "app_id":"` + `cli_a074eb7697f8d00b` + `",
            "app_secret":"` + `TMXU6KhEDdyvHImzSCk8RuXcCyjTu8zY` + `"
        }`
			req, _ := http.NewRequest(http.MethodPost, feishu.ServerUrl+"/open-apis/auth/v3/tenant_access_token/internal/", strings.NewReader(payload))
			return req
		},
	}

	// 创建 飞书 客户端
	FeishuClient := feishu.NewClient()

	progress <- 0
	// 调用 api 接口
	tenantAccessToken, _ := Atm.GetAccessToken()
	progress <- 0.1

	db.Delete(models.MeetingTopUserItem{}, "1=1")

	endDate := time.Now()
	endDate = endDate.Truncate(24 * time.Hour)
	startDate := endDate.AddDate(0, 0, -1)
	for i := 0; i < 30; i++ {
		params := url.Values{}
		params.Add(`start_time`, strconv.FormatInt(startDate.Unix(), 10))
		params.Add(`end_time`, strconv.FormatInt(endDate.Unix(), 10))
		params.Add(`limit`, `100`)
		params.Add(`order_by`, `2`)
		request, _ := http.NewRequest(http.MethodGet, feishu.ServerUrl+"/open-apis/vc/v1/reports/get_top_user?"+params.Encode(), nil)
		resp, err := FeishuClient.Do(request, tenantAccessToken)
		if err != nil {
			panic(err)
		}

		var result struct {
			Code int64 `json:"code"`
			Data struct {
				TopUserReport []models.MeetingTopUserItem `json:"top_user_report"`
			} `json:"data"`
			Msg string `json:"msg"`
		}
		err = json.Unmarshal(resp, &result)
		if err != nil {
			panic(err)
		}

		for index := range result.Data.TopUserReport {
			item := &result.Data.TopUserReport[index]
			item.StartTime = startDate
		}
		db.Save(result.Data.TopUserReport)

		progress <- 0.1 + 0.01*float32(i)
		endDate = startDate
		startDate = endDate.AddDate(0, 0, -1)
	}
	progress <- 1
	close(progress)
}

var PluginEntry Feishu
