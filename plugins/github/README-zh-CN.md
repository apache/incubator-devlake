# Github插件

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

## 获取 Access Token

下面的文档说明了如何获取`Github access token`：

https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token

## Github 限流

对于使用`Basic Authentication`或者`OAuth`的请求，限制为5000次/小时/token

- https://docs.github.com/en/rest/overview/resources-in-the-rest-api

通过在配置文件中设置多个token可以达到更高的请求速率

注意: 如果使用付费的企业版`Github`可以达到15000次/小时/token。

## 配置

在`.evn`文件中需要设置如下配置项

```

GITHUB_AUTH=XXX

or...

GITHUB_AUTH=XXX,YYY,ZZZ // 每个token属于不同的用户(可选)
```

如需使用代理需要设置`.env`文件中的`GITHUB_PROXY`配置项。如果此项没有配置或者配置错误，则不会使用代理。目前只支持`http`和`socks5`两种代理协议。

```
GITHUB_PROXY=http://127.0.0.1:1080
```


## 示例

```
curl --location --request POST 'localhost:8080/task' \
--header 'Content-Type: application/json' \
--data-raw '[[{
    "Plugin": "github",
    "Options": {
        "repositoryName": "lake",
        "owner": "merico-dev"
    }
}]]'
```
