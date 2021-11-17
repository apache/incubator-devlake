# Jenkins 插件

## 简介

本插件通过 [Remote Access API](https://www.jenkins.io/doc/book/using/remote-access-api/) 来收集 Jenkins 数据。然后根据收集到的原始数据，计算并展示相关的 devops 指标。

![image](https://user-images.githubusercontent.com/61080/141943122-dcb08c35-cb68-4967-9a7c-87b63c2d6988.png)

## 指标

指标名称 | 说明
:------------ | :-------------
构建数量 | 创建的构建数量
构建成功率 | 成功构建的百分比


## 配置

在使用本插件之前，您需要先在 `config-ui` 上面对插件进行配置。

### 通过 `config-ui` 进行配置

为了能访问到 Jenkins 的 API ，请确保完成以下的必填设置项。目前 Jenkins 只支持单一数据源，列表只会显示一个连接，同时其名称是固定不可修改的。多数据源支持会在不久的将来实现。

- Connection Name [只读]
  - ⚠️ D默认为 "Jenkins" 且不能修改。
- Endpoint URL (REST URL, 必须以 `https://` 或  `http://` 开头，`/` 结尾)
  - 必须指向一个有效的 REST API 端点, 比如 `https://jenkins.example.com/`
- Username (E-mail)
  - 该 Jenkins 实例上的有效用户名。
- Password (密码或 API 的 Acess Token)
  - 用户名对应的密码
  - 请参照 Jenkins 的官方文档中关于 "Using Credentials" 的说明
  - 您可以使用  **API Access Token** 代替密码, 可在 Jenkins 的面板中依次打开 `User` -> `Configure` -> `API Token` 进行生成。

完成上述项目设定后，请点击保存按钮更新连接的设置。

## 数据收集及计算

为了触发插件进行数据收集和计算，您需要构造一个 JSON， 通过 `config-ui` 中的 `Triggers` 功能，发送请求触发收集计算任务：

```json
[
  [
    {
      "plugin": "jenkins",
      "options": {}
    }
  ]
]
```
