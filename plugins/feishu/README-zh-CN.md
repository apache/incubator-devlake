# Feishu 插件

<div align="center">

| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |

</div>

<br>

## 简介

本插件通过 [Feishu Openapi](https://open.feishu.cn/document/home/user-identity-introduction/introduction) 来收集 Feishu 数据。

## 配置

在使用本插件之前，您需要先找到飞书管理员获取app_id和app_secret（请参照 Feishu 的官方文档中[相关说明](https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/tenant_access_token_internal)
），然后在 `.env` 上面对插件进行配置。

### 编辑.env文件

为了能访问到 Feishu 的 API ，请确保完成以下的必填设置项。目前 Feishu 只支持单一数据源，列表只会显示一个连接，同时其名称是固定不可修改的。多数据源支持会在不久的将来实现。

FEISHU_APPID=app_id

FEISHU_APPSCRECT=app_secret

## 数据收集及计算

为了触发插件进行数据收集和计算，您需要构造一个 JSON， 通过 `config-ui` 中的 `Triggers` 功能，发送请求触发收集计算任务：
numOfDaysToCollect: 收集的天数
rateLimitPerSecond: 每秒发送请求的数量（最大值为8）

```json
[
  [
    {
      "plugin": "feishu",
      "options": {
        "numOfDaysToCollect" : 80,
        "rateLimitPerSecond" : 5
      }
    }
  ]
]
```
