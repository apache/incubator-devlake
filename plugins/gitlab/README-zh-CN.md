# Gitlab 插件
此文档包括 3 部分：指标，配置，如何获取数据。

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

插件运行前，你需要在 `.env` 文件中进行以下配置：

    ```
    # Gitlab
    GITLAB_ENDPOINT=https://gitlab.com/api/v4/
    GITLAB_AUTH=<your access token>
    ```

### 获取自己的 Gitlab Endpoint <a id="gitlab-api-token"></a>
我们使用的 `GITLAB_ENDPOINT`  是 `https://gitlab.com/api/v4/`，但不同用户的 Endpoint 会有所区别。关于这个问题的更多信息，请参考 <a href="https://docs.gitlab.com/ee/api/" target="_blank">Gitlab官方API文档</a>






## 如何触发此插件进行数据收集

你可以向 `/task` 发起一个POST请求来触发数据收集。由于我们通过 Gitlab project id 来决定收集的数据范畴，因此需要在请求里带上 Gitlab project id<br>
注意：此请求会在收集全部数据时自动触发，你无需单独执行这一请求，也不需要在 `.env` 文件中设置这个。

    ```
    curl --location --request POST 'localhost:8080/task' \
    --header 'Content-Type: application/json' \
    --data-raw '[[{
        "plugin": "gitlab",
        "options": {
            "projectId": <Your gitlab project id>
        }
    }]]'
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
