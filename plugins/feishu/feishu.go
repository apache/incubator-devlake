package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/faabiosr/cachego/file"
	"github.com/fastwego/feishu"
	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
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

func (plugin Feishu) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	fmt.Println("start feishu plugin execution")

	// 内部应用 tenant_access_token 管理器
	// tenant_access_token manager
	Atm := &feishu.DefaultAccessTokenManager{
		Id:    `cli_a074eb7697f8d00b`,
		Cache: file.New(os.TempDir()),
		GetRefreshRequestFunc: func() *http.Request {
			payload := `{
                "app_id":"` + config.V.GetString("FEISHU_APPID") + `",
                "app_secret":"` + config.V.GetString("FEISHU_APPSCRECT") + `"
            }`
			req, _ := http.NewRequest(http.MethodPost, feishu.ServerUrl+"/open-apis/auth/v3/tenant_access_token/internal/", strings.NewReader(payload))
			return req
		},
	}

	// 创建 飞书 客户端
	// create feishu client
	FeishuClient := feishu.NewClient()

	progress <- 0
	// 调用 AccessToken api 接口
	// request AccessToken api
	tenantAccessToken, err := Atm.GetAccessToken()
	fmt.Println(tenantAccessToken)
	if err != nil {
		return err
	}
	progress <- 0.1

	lakeModels.Db.Delete(models.FeishuMeetingTopUserItem{}, "1=1")
	progress <- 0.2

	endDate := time.Now()
	endDate = endDate.Truncate(24 * time.Hour)
	startDate := endDate.AddDate(0, 0, -1)
	progress <- 0.3
	for i := 0; i < 120; i++ {
		params := url.Values{}
		params.Add(`start_time`, strconv.FormatInt(startDate.Unix(), 10))
		params.Add(`end_time`, strconv.FormatInt(endDate.Unix(), 10))
		params.Add(`limit`, `100`)
		params.Add(`order_by`, `2`)
		request, _ := http.NewRequest(http.MethodGet, feishu.ServerUrl+"/open-apis/vc/v1/reports/get_top_user?"+params.Encode(), nil)
		resp, err := FeishuClient.Do(request, tenantAccessToken)
		fmt.Println(resp)
		if err != nil {
			return err
		}

		var result struct {
			Code int64 `json:"code"`
			Data struct {
				TopUserReport []models.FeishuMeetingTopUserItem `json:"top_user_report"`
			} `json:"data"`
			Msg string `json:"msg"`
		}
		err = json.Unmarshal(resp, &result)
		if err != nil {
			return err
		}

		for index := range result.Data.TopUserReport {
			item := &result.Data.TopUserReport[index]
			item.StartTime = startDate
		}
		lakeModels.Db.Save(result.Data.TopUserReport)

		progress <- 0.3 + 0.01*float32(i)
		endDate = startDate
		startDate = endDate.AddDate(0, 0, -1)

		time.Sleep(100 * time.Millisecond)
	}
	progress <- 1
	return nil
}

func (plugin Feishu) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/feishu"
}

func (plugin Feishu) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
	}
}

var PluginEntry Feishu

// standalone mode for debugging
func main() {
	PluginEntry.Init()
	progress := make(chan float32)
	go func() {
		err := PluginEntry.Execute(
			map[string]interface{}{},
			progress,
			context.Background(),
		)
		if err != nil {
			panic(err)
		}
	}()
	for p := range progress {
		fmt.Println(p)
	}
}
