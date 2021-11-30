# Gitlab 插件

<div align="center">

| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |

</div>

<br>

## 指标
此插件通过收集 Gitlab 的数据来计算以下指标。

指标名称 | 描述
:------------ | :-------------
代码评审次数 | PR/MR创建的数量
代码评审通过率 | PR/MR被合并的比率
代码评审人数 | 评审PR/MR的人数
代码评审时长 | 从PR/MR创建到被合并的时间
代码提交人数 | 提交了Commit的人数
代码提交次数 | 提交Commit的次数
新增代码行数 | 累积新增的代码行数
删除代码行数 | 累计删除的代码行数
代码评审轮数 | PR/MR创建到被合并期间，经过了多少轮的评审


## 配置

### 数据源连接配置
配置界面需要填入以下字段
- **Connection Name** [`只读`]
    - ⚠️ 默认值为 "**Gitlab**" 请不要改动。
- **Endpoint URL** (REST URL, 以 `https://`或`http://`开头)
    - 应当填入可用的REST API Endpoint。例如 `https://gitlab.com/api/v4/`
    - ⚠️url应当以`/`结尾
- **Personal Access Token** (HTTP Basic Auth)
    - 登录你的Gitlab并创建**Personal Access Token**，token长度必须是20个字符。请把生成的token安全保存离开页面后将无法看到。

    1. 右上角选择**avatar**。
    2. 选择**Edit profile**。
    3. 在左侧边栏选择**Access Tokens**。
    4. 输入**name**并且为此token选择**expiry date**。
    5. 选择你所需的**scopes**。
    6. 选择**Create personal access token**。
如何创建**personal access token**，请参考官方文档[GitLab Docs on Personal Tokens](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)

关于**GitLab REST API**的更多信息请参考官方文档[GitLab Docs on REST](https://docs.gitlab.com/ee/development/documentation/restful_api_styleguide.html#restful-api)

点击**Save Connection**保存配置。

### 数据源配置
当前只有一个**可选**配置，它可以让你将JIRA Boards和GitLab Projects关联起来。

- **JIRA Board Mappings [ `可选`]**
  **Map JIRA Boards to GitLab**。请以以下格式输入映射规则
```
# 映射JIRA Board ID 8 ==> Gitlab Projects 8967944,8967945
<JIRA_BOARD>:<GITLAB_PROJECT_ID>,...; 例如 8:8967944,8967945;9:8967946,8967947
```
点击**Save Settings**保存配置。


## 收集数据

你可以向 `/pipelines` 发起一个POST请求来触发数据收集。

```
curl --location --request POST 'localhost:8080/pipelines' \
--header 'Content-Type: application/json' \
--data-raw '
{
    "name": "gitlab 20211126",
    "tasks": [[{
        "plugin": "gitlab",
        "options": {
            "projectId": <Your gitlab project id>
        }
    }]]
}
'
```

## 如何获取 Gitlab Project ID

要获得一个特定的 Gitlab 仓库的项目ID：
- 访问 Gitlab 的仓库页面
- 找到标题下面的项目ID

  ![Screen Shot 2021-08-06 at 4 32 53 PM](https://user-images.githubusercontent.com/3789273/128568416-a47b2763-51d8-4a6a-8a8b-396512bffb03.png)

- 将此项目ID复制在上方的请求示例中，从这个项目收集数据

### 创建一个 Gitlab API Token <a id="gitlab-api-token"></a>

1. 登录 Gitlab 后，访问 `https://gitlab.com/-/profile/personal_access_tokens`
2. Token 可以设置任意名称，不要设置过期日期。在设置范围时，去掉“写入”权限

   ![Screen Shot 2021-08-06 at 4 44 01 PM](https://user-images.githubusercontent.com/3789273/128569148-96f50d4e-5b3b-4110-af69-a68f8d64350a.png)

3. 点击 **Create Personal Access Token** 按钮
4. 通过 config-ui 或者 直接将 API Token 复制并保存到 `.env` 文件中
