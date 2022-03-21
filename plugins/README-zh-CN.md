# 听说你想建立一个新的插件...

...好消息是，这很容易!


## 基本写法

```golang
type YourPlugin string

func (plugin YourPlugin) Description() string {
	return "To collect and enrich data from YourPlugin"
}

func (plugin YourPlugin) Execute(options map[string]interface{}, progress chan<- float32) {
	logger.Print("Starting YourPlugin execution...")

  // 检查选项中需要的字段
	projectId, ok := options["projectId"]
	if !ok {
		logger.Print("projectId is required for YourPlugin execution")
		return
	}

  // 开始收集
  if err := tasks.CollectProject(projectId); err != nil {
		logger.Error("Could not collect projects: ", err)
		return
	}
  // 处理错误
  if err != nil {
    logger.Error(err)
  }

  // 导出一个名为 PluginEntry 的变量供 Framework 搜索和加载
  var PluginEntry YourPlugin //nolint
}
```

## 概要

要建立一个新的插件，你将需要做下列事项。你应该选择一个你想看的数据的 API。首先考虑你想看到的指标，然后寻找能够支持这些指标的数据。

## 收集（Collection）

然后你要写一个 `Collection` 来收集数据。你需要阅读一些 API 文档，弄清楚你想在最后的 Grafana 仪表盘中看到哪些指标（配置Grafana是最后一步）。

## 构建一个 `Fetcher` 来执行请求

Plugins/core文件夹包含一个 API 客户端，你可以在自己的插件中实现。它有一些方法，比如Get()。<br>
每个API处理分页的方式不同，所以你可能需要实现一个 "带分页的获取 "方法。有一种方法是使用 "ant" 包作为管理并发任务的方法：https://github.com/panjf2000/ants

你的 collection 方法可能看起来像这样:

```golang
func Collect() error {
	pluginApiClient := CreateApiClient()

	return pluginApiClient.FetchWithPagination("<your_api_url>",
		func(res *http.Response) error {
			pluginApiResponse := &ApiResponse{}
      // 你必须解除对api的响应，才能使用这些结果
			err := helper.UnmarshalResponse(res, pluginApiResponse)
			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
      // 将获取到的数据保存到数据库中
			for _, value := range *pluginApiResponse {
				pluginModel := &models.pluginModel{
					pluginId:       value.pluginId,
					Title:          value.Title,
					Message:        value.Message,
				}

				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&pluginModel).Error

				if err != nil {
					logger.Error("Could not upsert: ", err)
				}
			}

			return nil
		})
}
```

请注意 "upsert" 的使用。这对于只保存修改过的记录是很有用的。

## 数据处理（Enrichment）
  
一旦你通过 API 收集了数据，你可能想通过以下方式来对这些数据做 ETL。比如：

  - 添加你目前没有的字段
  - 计算你可能需要的指标字段
  - 消除你不需要的字段

## 你已经完成了!

祝贺你! 你已经创建了你的第一个插件! 🎖
