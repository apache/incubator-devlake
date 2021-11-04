# Jira 插件

## 概述

此插件通过 Jira Cloud REST API 收集 Jira 数据。然后，它从 Jira 数据中计算出各种工程指标并使之可视化。

<img width="2035" alt="Screen Shot 2021-09-10 at 4 01 55 PM" src="https://user-images.githubusercontent.com/2908155/132926143-7a31d37f-22e1-487d-92a3-cf62e402e5a8.png">

## Project Metrics This Covers

指标名称 | 描述
:------------ | :-------------
需求数 | 类型为 "需求" 的事务的数量
需求交付时间 | 类型为 "需求" 的事务的交付时间，即从创建到完成的时间
需求交付率 | 已交付的需求/所有需求的比率
需求粒度 | 一个"需求"类型事务的标准故事点
故障数量 | 类型为 "故障" 的事务数量<br><i>测试中发现的Bug</i>。
故障修复时间 |类型为 "故障" 的事务的修复时间
测试故障率（代码行） | 每1000行代码产生的 "故障" 数量<br><i>包括新增和删除的行数</i>
测试故障数 | 类型为 "故障" 的事务数量<br><i>Incident在生产中运行时发现的问题</i>。
质量事故数 | "Incident" 类型的事务的准备时间
质量事故率（代码行） | 每1000行代码产生的 Incident 数量<br><i>包括新增和删除的行数</i>

## 配置

插件运行前，你需要在 `.env` 文件中进行以下配置：

### 设置 JIRA_ENDPOINT

这是所有 Jira API 调用的基础设置。你可以在你所有的 Jira Urls 中看到它作为 Url 的开头。

例如：如果你看到 `https://mydomain.atlassian.net/secure/RapidBoard.jspa?rapidView=999&projectKey=XXX`, 你需要在你的 `.env` 文件中设置 `JIRA_ENDPOINT=https://mydomain.atlassian.net`

### 生成 API token

1. 登录Jira后，访问网址 `https://id.atlassian.com/manage-profile/security/api-tokens`
2. 点击 **Create API Token** 按钮，随便取个标签名
![image](https://user-images.githubusercontent.com/27032263/129363611-af5077c9-7a27-474a-a685-4ad52366608b.png)
3. 使用 `echo -n <jira login email>:<jira token> | base64` 命令对登录的电子邮件进行编码

### 设置事务类型的映射<a id="issue-type-mapping"></a>


不同公司可能使用不同的事务类型来表示他们的 需求/故障/事故，类型映射允许 Devlake 识别你在 Jira 类型方面的具体设置。Devlake 支持 3 种不同的标准状态类型：
 
 - `需求（Requirement）`
 - `故障（Bug）`
 - `事故（Incident）`

例如，假设我们使用 `故事` 和 `任务` 来表示需求，用 `客户投诉` 表示事故，用 `QABug` 表示故障。我们要做的是在运行 Devlake 之前在 `.env` 文件下设置环境变量：

```sh
# JIRA_ISSUE_TYPE_MAPPING=<STANDARD_TYPE>:<YOUR_TYPE_1>,<YOUR_TYPE_2>;....
JIRA_ISSUE_TYPE_MAPPING=Requirement:故事,任务;Incident:客户投诉;Bug:QABug;
```

事务类型映射对于一些指标来说是至关重要的，比如**需求数**，请确保正确映射你的自定义类型。

### 设置事务状态的映射<a id="issue-status-mapping"></a>

Jira 是高度可定制的，不同公司可能使用不同的状态来表示一个事务是否被解决。一个公司可能将这个状态命名为 "Done"，而其他公司可能将其命名为 "Finished"。

为了正确的收集事务的生命周期信息，你必须将自己使用的事务状态映射到 Devlake 的标准状态，Devlake 支持 2 种标准状态：

- `Resolved`: 事务完成或解决
- `Rejected`: 事务终止或取消

例如，假设我们
- 对于 `Bug` 类型的事务：使用 `已修复` 表示 "Resolved"，用 `拒绝`和 `无法复现` 表示 "Rejected"；
- 对于 `Incident` 类型的事务：使用 `已修复` 表示 "Resolved"，用 `拒绝` 表示 "Rejected"；
- 对于 `Story` 类型的事务：使用 `已完成` 表示 "Resolved"，用 `推迟` 表示 "Rejected"；
我们要做的是在运行 Devlake 之前在 `.env` 文件下设置环境变量：

```sh
#JIRA_ISSUE_<YOUR_TYPE>_STATUS_MAPPING=<STANDARD_STATUS>:<YOUR_STATUS>;...
JIRA_ISSUE_BUG_STATUS_MAPPING=Resolved:已修复;Rejected:拒绝,无法复现
JIRA_ISSUE_INCIDENT_STATUS_MAPPING=Resolved:已修复;Rejected:拒绝
JIRA_ISSUE_STORY_STATUS_MAPPING=Resolved:已完成;Rejected:推迟
```

状态映射对于像 `需求前置时间` 和 `故障/事故修复时间` 这样的指标至关重要，因为我们需要通过需求/故障/事故的解决时间来计算。


## 设置 Jira 的自定义字段
此设置适用于配置 `JIRA_ISSUE_EPIC_KEY_FIELD` 和 `JIRA_ISSUE_STORYPOINT_FIELD`

- `JIRA_ISSUE_EPIC_KEY_FIELD` 表示一个 Jira 事务所属的 `史诗` 的 Key，比如 EE-234
- `JIRA_ISSUE_STORYPOINT_FIELD` 表示 Jira 事务的故事点。配置此字段是为了将你本地的

一个完整的设置形如: 
```sh
JIRA_ISSUE_EPIC_KEY_FIELD=customfield_10024
JIRA_ISSUE_STORYPOINT_FIELD=customfield_10026
```

请遵循此指南，[如何查找 Jira 的自定义字段的ID?](https://github.com/merico-dev/lake/wiki/How-to-find-the-custom-field-ID-in-Jira)



### 设置 JIRA_ISSUE_STORYPOINT_COEFFICIENT

如果你并未使用故事点，而是使用 `预计时间` 等字段来表示需求粒度，你可以将上方的`JIRA_ISSUE_STORYPOINT_FIELD`设置为 `预计时间` 对应的自定义字段，然后将 `JIRA_ISSUE_STORYPOINT_COEFFICIENT` 设置为默认值 1 以外的值。

我们使用的标准故事点字段 `std_story_point` 的值 = 预计时间 * JIRA_ISSUE_STORYPOINT_COEFFICIENT

一般情况下，1 个标准故事点为半天，也就是 4 小时，因此当你使用 `预计时间` 作为需求粒度的字段，并以小时为单位，那么可以将 `JIRA_ISSUE_STORYPOINT_COEFFICIENT` 设置为 `0.25`


## 如何触发此插件进行数据收集
 
你可以向 `/task` 发起一个POST请求来触发数据收集。由于我们通过 Jira 来决定收集的数据范畴，因此需要在请求里带上 Jira board id<br>
注意：此请求会在收集全部数据时自动触发，你无需单独执行这一请求，也不需要在 `.env` 文件中设置这个。

```
curl -XPOST 'localhost:8080/task' \
-H 'Content-Type: application/json' \
-d '[[{
    "plugin": "jira",
    "options": {
        "boardId": 8,
        "since": "2006-01-02T15:04:05Z",
        "sourceId": 1
    }
}]]'
```


### 如何获取 Jira Board Id
1. 打开浏览器，进入待导入的 Jira 面板
2. 在 URL 的参数 `?rapidView=` 中获取面板 ID


例如: 对于 `https://<your_jira_endpoint>/secure/RapidBoard.jspa?rapidView=39`，面板的ID是39

![Screen Shot 2021-08-13 at 10 07 19 AM](https://user-images.githubusercontent.com/27032263/129363083-df0afa18-e147-4612-baf9-d284a8bb7a59.png)

