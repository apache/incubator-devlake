# Github插件

<div align="center">

| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |

</div>

<br>

## 概述

此插件从`Github`收集数据并通过`Grafana`展示。我们可以为技术领导者回答诸如以下问题：
- 本月是否比以往更高产？
- 我们能多快地响应客户需求？
- 质量是否有提升？

## 指标

以下是几个利用`Github`数据的例子：
- 每个人的平均需求研发时间
- 千行代码Bug数
- 提交数依时间分布

## 截图

![image](https://user-images.githubusercontent.com/27032263/141855099-f218f220-1707-45fa-aced-6742ab4c4286.png)


## 配置

### 数据源连接配置
配置界面需要填入以下字段
- **Connection Name** [`只读`]
    - ⚠️ 默认值为 "**Github**" 请不要改动。
- **Endpoint URL** (REST URL, 以 `https://`或`http://`开头)
    - 应当填入可用的REST API Endpoint。例如 `https://api.github.com/`
    - ⚠️url应当以`/`结尾
- **Auth Token(s)** (Personal Access Token)
    - 如何创建**personal access token**，请参考官方文档[GitHub Docs on Personal Tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
    - 填入至少一个token，可以填入多个token并以英文逗号`,`间隔，填入多个token可以加快数据采集速度

对于使用`Basic Authentication`或者`OAuth`的请求，限制为5000次/小时/token
- https://docs.github.com/en/rest/overview/resources-in-the-rest-api
通过在配置文件中设置多个token可以达到更高的请求速率

注意: 如果使用付费的企业版`Github`可以达到15000次/小时/token。
关于**GitHub REST API**的更多信息请参考官方文档[GitHub Docs on REST](https://docs.github.com/en/rest)

点击**Save Connection**保存配置。


### 数据源配置
目前只有一个**可选**配置*Proxy URL*，如果你需要代理才能访问GitHub才需要配置此项
- **GitHub Proxy URL [`可选`]**
  - 输入可用的代理服务器地址，例如：`http://your-proxy-server.com:1080`

点击**Save Settings**保存配置。


## 示例

```
curl --location --request POST 'localhost:8080/pipelines' \
--header 'Content-Type: application/json' \
--data-raw '
{
    "name": "github 20211126",
    "tasks": [[{
        "plugin": "github",
        "options": {
            "repositoryName": "lake",
            "owner": "merico-dev"
        }
    }]]
}
'
```
