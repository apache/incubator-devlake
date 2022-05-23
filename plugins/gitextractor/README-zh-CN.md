# Git Repo Extractor插件

## Summary

该插件可以从远端或本地git仓库提取commit和reference信息，并保存到数据库或csv文件。

## 示例

```
curl --location --request POST 'localhost:8080/pipelines' \
--header 'Content-Type: application/json' \
--data-raw '
{
    "name": "git repo extractor",
    "tasks": [
        [
            {
                "Plugin": "gitextractor",
                "Options": {
                    "url": "https://github.com/apache/incubator-devlake.git",
                    "repoId": "github:GithubRepos:384111310"
                }
            }
        ]
    ]
}
'
```
- `url`: git仓库的位置，如果是远端仓库应当以`http`或`https`开头， 如果是本地仓库则应当以`/`开头
- `repoId`: `repos`表的`id`字段.
- `proxy`: 可选, 只支持`http`代理，例如：`http://your-proxy-server.com:1080`.
- `user`: 可选, 通过HTTP/HTTPS协议克隆私有代码库时使用
- `password`: 可选, 通过HTTP/HTTPS协议克隆私有代码库时使用
- `privateKey`: 可选, 通过SSH协议克隆代码库时使用, 值为经过base64编码的`PEM`文件
- `passphrase`: 可选, 私钥的密码

## 独立运行本插件

本插件可以作为独立于DevLake服务的命令行工具使用:

```
go run plugins/gitextractor/main.go -url https://github.com/apache/incubator-devlake.git -id github:GithubRepo:384111310 -db "merico:merico@tcp(127.0.0.1:3306)/lake?charset=utf8mb4&parseTime=True"
```

如果想了解命令行工具的更多选项，比如如何输出收集结果到csv文件，请直接阅读`plugins/gitextractor/main.go`。