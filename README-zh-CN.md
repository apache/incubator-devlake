<div align="center">
<br />
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

# DevLake

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/discord/844603288082186240.svg?style=flat?label=&logo=discord&logoColor=ffffff&color=747df7&labelColor=454545)](https://discord.gg/83rDG6ydVZ)
![badge](https://github.com/merico-dev/lake/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/merico-dev/lake)](https://goreportcard.com/report/github.com/merico-dev/lake)


| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |

</div>
<br>
<div align="left">

### 什么是 DevLake？
DevLake 将你所有 DevOps 工具里的数据以实用、个性化、可扩展的视图呈现。通过 DevLake，从不断增加的工具列表中收集、分析和可视化数据。

DevLake 适用于希望更好地通过数据了解其开发过程的开发团队，以及希望以数据驱动提升自身实践的开发团队。有了 DevLake，你可以向你的开发过程提出任何问题，只要连接数据并查询。


<a href="https://app-259373083972538368-3002.ars.teamcode.com/d/0Rjxknc7z/demo-homepage?orgId=1">查看 demo</a>。用户名/密码： test/test。Demo里呈现的数据来自本仓库 merico-dev/lake。


#### 开始安装 DevLake
<table>
  <tr>
    <td valign="middle"><a href="#user-setup">运行 DevLake</a></td>
  </tr>
</table>



<br>

<div align="left">
<img src="https://user-images.githubusercontent.com/14050754/142356580-40637a30-5578-48ed-8e4a-128cd0738e3e.png" width="100%" alt="User Flow" style="border-radius:15px;"/>
<p align="center">用户使用流程</p><br>



### DevLake 可以完成什么?
1. 归集 DevOps 全流程效能数据，连接数据孤岛
2. 标准的<a href="https://github.com/merico-dev/lake/wiki/DataModel.Domain-layer-schema" target="_blank">研发数据模型</a>和开箱即用的<a href="https://github.com/merico-dev/lake/wiki/Metric-Cheatsheet" target="_blank">效能指标</a>
3. 灵活的数据收集、ETL的<a href="https://github.com/merico-dev/lake/blob/main/ARCHITECTURE.md">框架</a>，支持自定义分析



<br>

## 用户安装<a id="user-setup"></a>

- 如果你只打算运行 DevLake，你只需要阅读这一小节<br>
- 如果你想在云端安装 DevLake，你可以参考[安装手册](https://github.com/merico-dev/lake/wiki/How-to-Set-Up-Dev-Lake-with-Tin-zh-CN)，点击 <a valign="middle" href="https://www.teamcode.com/tin/clone?applicationId=259777118600769536">
        <img
          src="https://static01.teamcode.com/badge/teamcode-badge-run-in-cloud-cn.svg"
          width="120px"
          alt="Teamcode" valign="middle"
        />
      </a> 完成安装
- 写成 `这样` 的命令需要在你的终端中运行

#### 需要安装的软件包<a id="user-setup-requirements"></a>

- [Docker](https://docs.docker.com/get-docker)
- [docker-compose](https://docs.docker.com/compose/install/)

注：安装完 Docker 后，你可能需要运行 Docker 应用程序并重新启动你的终端

#### 在你的终端中运行以下命令<a id="user-setup-commands"></a>

**IMPORTANT（新用户可以忽略）: DevLake暂不支持向前兼容。当 DB Schema 发生变化时，直接更新已有实例可能出错，建议已经安装 DevLake 的用户在升级时，重新部署实例并导入数据。**

1. 在[最新版本列表](https://github.com/merico-dev/lake/releases/latest) 下载 `docker-compose.yml` 和 `env.example`
2. 将 `env.example` 重命名为 `.env`。Mac/Linux 用户请在命令行里运行 `mv env.example .env` 来完成修改
3. 启动 Docker，然后运行 `docker-compose up -d` 启动服务
4. 访问 `localhost:4000` 来设置 DevLake 的配置文件
   >- 在 Integrations 页面上找到你想要导入的数据源
   >- 了解如何配置每个数据源：<br>
      > <a href="plugins/jira/README-zh-CN.md" target="_blank">Jira</a><br>
      > <a href="plugins/gitlab/README-zh-CN.md" target="_blank">GitLab</a><br>
      > <a href="plugins/jenkins/README-zh-CN.md" target="_blank">Jenkins</a><br>
      > <a href="plugins/github/README-zh-CN.md" target="_blank">GitHub</a><br>
   >- 提交表单，通过点击每个表单页面上的**Save Connection**按钮来更新数值。
   >- `devlake`需要一段时间才能完全启动。如果`config-ui`提示 API 无法访问，请等待几秒钟并尝试刷新页面。

5. 访问 `localhost:4000/pipelines/create`，创建 1个Pipeline run，并触发数据收集

   Pipeline Runs 可以通过新的 "Create Run"界面启动。只需启用你希望运行的**数据源**，并指定数据收集的范围，比如Gitlab的项目ID和GitHub的仓库名称。

   一旦创建了有效的 Pipeline Run 配置，按**Create Run**来启动/运行该 Pipeline。
   Pipeline Run 启动后，你会被自动转到**Pipeline Activity**界面，以监控采集活动。

   **Pipelines**可从 config-ui 的主菜单进入。

   - 管理所有Pipeline: `http://localhost:4000/pipelines`。
   - 创建Pipeline Run: `localhost:4000/pipelines/create`。
   - 查看Pipeline Activity: `http://localhost:4000/pipelines/activity/[RUN_ID]`。

   对于复杂度较高的用例，请使用Raw JSON API进行任务配置。使用**cURL**或图形API工具（如**Postman**）手动启动运行。`POST`以下请求到DevLake API端点。

   >   ```json
   >   [
   >     [
   >       {
   >         "Plugin": "github",
   >         "Options": {
   >           "repo": "lake",
   >           "owner": "merico-dev"
   >         }
   >       }
   >     ]
   >   ]
   >   ```
   
   请参考这篇 wiki [How to trigger data collection](https://github.com/merico-dev/lake/wiki/How-to-use-the-triggers-page).

6. 数据收集完成后，点击配置页面左上角的 *View Dashboards* 按钮或者访问 `localhost:3002`，访问 Grafana (用户名: `admin`, 密码: `admin`)

   我们使用 <a href="https://grafana.com/" target="_blank">Grafana</a> 作为可视化工具，为存储在<a href="https://github.com/merico-dev/lake/wiki/DataModel.Domain-layer-schema">我们数据库中的数据</a>建立图表。可以使用SQL查询，添加面板来构建、保存和编辑自定义仪表盘。

   关于配置和定制仪表盘的所有细节可以在 [Grafana 文档](docs/GRAFANA.md) 中找到。

#### 设置 Cron job
为了定期同步数据，我们提供了[`lake-cli`](./cmd/lake-cli/README.md)以方便发送数据收集请求，我们同时提供了[cron job](./devops/sync/README.md)以定期触发 cli 工具。

<br>

### 开发者安装<a id="dev-setup"></a>

#### 前期准备

- <a href="https://docs.docker.com/get-docker" target="_blank">Docker</a>
- <a href="https://golang.org/doc/install" target="_blank">Golang</a>
- Make
  - Mac (Already installed)
  - Windows: [Download](http://gnuwin32.sourceforge.net/packages/make.htm)
  - Ubuntu: `sudo apt-get install build-essential`

#### 如何设置开发环境
1. 进入你想安装本项目的路径，并克隆资源库

   ```sh
   git clone https://github.com/merico-dev/lake.git
   cd lake
   ```

2. 安装插件依赖

   - [RefDiff](plugins/refdiff#development)

2. 安装 go packages

    ```sh
	go get
    ```

3. 将样本配置文件复制到新的本地文件

    ```sh
    cp .env.example .env
    ```

4. 在`.env`文件中找到以`DB_URL`开头的那一行，把`mysql:3306`替换为`127.0.0.1:3306`

5. 启动 MySQL 和 Grafana

    > 确保在此步骤之前 Docker 正在运行。

    ```sh
    docker-compose up -d mysql grafana
    ```


6. 在 2 个终端种分别以开发者模式运行 lake 和 config UI:

    ```sh
    # run lake
    make dev
    # run config UI
    make configure-dev
    ```

7. 访问 config-ui `localhost:4000` 来配置 DevLake 数据源
   >- 在 "Integration"页面上找到到所需的插件页面
   >- 你需要为你打算使用的插件输入必要的信息
   >- 请参考以下内容，以了解如何配置每个插件的更多细节
   >-> <a href="plugins/jira/README-zh-CN.md" target="_blank">Jira</a>
   >-> <a href="plugins/gitlab/README-zh-CN.md" target="_blank">GitLab</a>
   >-> <a href="plugins/jenkins/README-zh-CN.md" target="_blank">Jenkins</a>
   >-> <a href="plugins/github/README-zh-CN.md" target="_blank">GitHub</a>


8. 访问 `localhost:4000/pipelines/create`，创建 1个Pipeline run，并触发数据收集

   Pipeline Runs 可以通过新的 "Create Run"界面启动。只需启用你希望运行的数据源，并指定数据收集的范围，比如Gitlab的项目ID和GitHub的仓库名称。

   一旦创建了有效的 Pipeline Run 配置，按**Create Run**来启动/运行该 Pipeline。
   Pipeline Run 启动后，你会被自动转到**Pipeline Activity**界面，以监控采集活动。

   **Pipelines**可从 config-ui 的主菜单进入。

   - 管理所有Pipeline: `http://localhost:4000/pipelines`。
   - 创建Pipeline Run: `http://localhost:4000/pipelines/create`。
   - 查看Pipeline Activity: `http://localhost:4000/pipelines/activity/[RUN_ID]`。

   对于复杂度较高的用例，请使用Raw JSON API进行任务配置。使用**cURL**或图形API工具（如**Postman**）手动启动运行。`POST`以下请求到DevLake API端点。

   >   ```json
   >   [
   >     [
   >       {
   >         "Plugin": "github",
   >         "Options": {
   >           "repo": "lake",
   >           "owner": "merico-dev"
   >         }
   >       }
   >     ]
   >   ]
   >   ```

   请参考这篇 wiki [How to trigger data collection](https://github.com/merico-dev/lake/wiki/How-to-use-the-triggers-page).


9. 数据收集完成后，点击配置页面左上角的 *View Dashboards* 按钮或者访问 `localhost:3002`(用户名: `admin`, 密码: `admin`)

   我们使用 <a href="https://grafana.com/" target="_blank">Grafana</a> 作为可视化工具，为存储在<a href="https://github.com/merico-dev/lake/wiki/DataModel.Domain-layer-schema">我们数据库中的数据</a>建立图表。可以使用SQL查询，添加面板来构建、保存和编辑自定义仪表盘。

   关于配置和定制仪表盘的所有细节可以在 [Grafana 文档](docs/GRAFANA.md) 中找到。

10. （可选）运行测试: 

    ```sh
    make test
    ```

<br>

## 项目路线图
- <a href="https://github.com/merico-dev/lake/wiki/Roadmap-2022" target="_blank">2022年路线图</a>: 2022年的目标和路线图
- DevLake 已经支持的数据源：
    - <a href="plugins/jira/README.md" target="_blank">Jira(Cloud)</a>
    - <a href="plugins/gitextractor/README.md" target="_blank">Git</a>
    - <a href="plugins/github/README.md" target="_blank">GitHub</a>
    - <a href="plugins/gitlab/README.md" target="_blank">GitLab(Cloud)</a>
    - <a href="plugins/jenkins/README.md" target="_blank">Jenkins</a>
- <a href="https://github.com/merico-dev/lake/wiki/Metric-Cheatsheet" target="_blank">已经支持的指标</a>: 为观测和分析提供不同的视角

<br>

## 贡献
本节列出了所有与共建 DevLake 相关的文档

- [架构设计](ARCHITECTURE.md): DevLake的架构设计
- [添加一个插件](/plugins/README.md): 如何添加一个新插件
- [添加新的指标](/plugins/HOW-TO-ADD-METRICS.md): 如何在一个插件里添加新的指标
- [贡献规范](CONTRIBUTING.md): 如果你想给 DevLake 贡献代码，请看下这个文档

<br>

## 社区

- <a href="https://discord.com/invite/83rDG6ydVZ" target="_blank">Discord</a>: 在 Discord 上给我们发消息
- <a href="https://github.com/merico-dev/lake/wiki/FAQ" target="_blank">FAQ</a>: 常见问题汇总

<br>

### License<a id="license"></a>

此项目的许可证为 Apache License 2.0 - 查看 [许可证](LICENSE) 详情。
