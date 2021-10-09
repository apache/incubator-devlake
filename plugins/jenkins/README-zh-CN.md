# Jenkins 插件
此文档包括 3 部分：指标，配置，如何获取数据。

## 指标
此插件通过收集 Jenkins 的数据来计算以下指标。

指标名称 | 说明
:------------ | :-------------
构建数量 | 创建的构建数量
构建成功率 | 成功构建的百分比


## 配置

插件运行前，你需要在 `.env` 文件中进行以下配置：

```
# Jenkins configuration

JENKINS_ENDPOINT=https://jenkins.merico.cn/
JENKINS_USERNAME=your user name here
JENKINS_PASSWORD=your password or jenkins token here
```

你可以通过以下步骤找到 Jenkins API Token：`User` -> `Configure` -> `API Token`

## 如何触发此插件进行数据收集

你可以向 `/task` 发起一个POST请求来触发数据收集。<br>
注意：此请求会在收集全部数据时自动触发，你无需单独执行这一请求，也不需要在 `.env` 文件中设置这个。

```
curl --location --request POST 'localhost:8080/task' \
  --header 'Content-Type: application/json' \
  --data-raw '[[{
      "plugin": "jenkins",
      "options": {}
  }]]'
```
