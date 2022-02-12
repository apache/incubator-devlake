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
                    "url": "https://github.com/merico-dev/lake.git",
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