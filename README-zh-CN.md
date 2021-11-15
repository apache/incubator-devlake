<br />
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

# Dev Lake
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/discord/844603288082186240.svg?style=flat?label=&logo=discord&logoColor=ffffff&color=747df7&labelColor=454545)](https://discord.gg/83rDG6ydVZ)
![badge](https://github.com/merico-dev/lake/actions/workflows/main.yml/badge.svg)



| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |

<br>

### 什么是 Dev Lake？

Dev Lake 是一个研发效能分析平台，它通过对软件开发生命周期（SDLC）中产生的数据进行 _**整合、分析和可视化**_ ，提升研发效能。

<img src="https://user-images.githubusercontent.com/2908155/130271622-827c4ffa-d812-4843-b09d-ea1338b7e6e5.png" width="100%" alt="Dev Lake Grafana Dashboard" />
<p align="center">数据面板截图</p><br>
<img src="https://user-images.githubusercontent.com/14050754/139076905-48d13e40-51ab-49e4-b537-0fe56960a1c0.png" width="100%" alt="Dev Lake Grafana Dashboard" />
<p align="center">用户使用流程</p><br>

### 为什么选择 Dev Lake？
1. 将多个来源的数据（Jira、Gitlab、Jenkins等）统一到一个地方。
2. 可以一起计算来自不同数据源的指标
3. 提供一系列行业标准指标来识别工程问题 
4. 高度可定制，用户可以制作自己的图表、指标和仪表盘

### Dev Lake 可以完成什么?
1. 在一个个性化的、统一的视图中可视化和分析你的整个SDLC过程
2. 调试过程和团队层面的问题，扩大成功的规模
3. 统一和规范成功的衡量标准和基准
<br>

## 内容

目录 | 描述 | 文档链接
:------------ | :------------- | :-------------
数据源 | 链接到具体的插件使用和细节 | [查看本节](#data-source-plugins)
用户设置 | 以用户身份运行项目的步骤 | [查看本节](#user-setup) 
开发者设置 | 如何设置开发环境 | [查看本节](#dev-setup)
测试 | 运行测试的命令 | [查看本节](#tests)
Grafana | 如何将数据进行可视化 | [查看本节](#grafana)
添加一个插件 | 如何制作自己的插件的详细信息 | [链接](plugins/README-zh-CN.md) 
添加新的指标 | 如何给插件添加指标 | [链接](plugins/HOW-TO-ADD-METRICS-zh-CN.md) 
贡献 | 如何进行贡献 | [链接](CONTRIBUTING-zh-CN.md)
FAQ | 常见问题 | [链接](#faq)


## 我们目前支持的数据源<a id="data-source-plugins"></a>

下面是一个 _数据源插件（data source plugins）_ 的列表，用于收集和处理特定来源的数据。每个插件都有一个 `README.md` 文件，包含基本设置、故障排除和指标信息。

关于建立一个新的 _data source plugins_ 的更多信息，请参见[添加一个插件](plugins/README-zh-CN.md)。

目录 | 内容 | 文档
------------ | ------------- | -------------
Jira | 指标，生成 API Token，查找项目/看板ID，配置事务状态和字段名称 | [Link](plugins/jira/README-zh-CN.md) 
Gitlab | 指标，生成 API Token | [Link](plugins/gitlab/README-zh-CN.md) 
Jenkins | 指标，生成 API Token | [Link](plugins/jenkins/README-zh-CN.md) 


## 用户设置<a id="user-setup"></a>

**注意：如果你只打算运行 Dev Lake，你只需要阅读这一小节**<br>
**注意：写成 `这样` 的命令需要在你的终端中运行**

### 需要安装的软件包<a id="user-setup-requirements"></a>

- [Docker](https://docs.docker.com/get-docker)
- [docker-compose](https://docs.docker.com/compose/install/)

**NOTE:** 安装完 Docker 后，你可能需要运行 Docker 应用程序并重新启动你的终端

### 在你的终端中运行的命令<a id="user-setup-commands"></a>

1. 克隆仓库

   ```sh
   git clone https://github.com/merico-dev/lake.git devlake
   cd devlake
   cp .env.example .env
   ```
2. 启动 Docker，然后运行 `docker-compose up config-ui` 来启动配置界面。

> 关于如何配置插件的更多信息，请参考 <a href="https://github.com/merico-dev/lake#data-source-plugins" target="_blank">数据源插件</a> 部分

3. 访问 `localhost:4000` 来设置配置文件

4. 运行 `docker-compose up -d` 来启动其他服务

5. 访问 `localhost:4000/triggers` 以运行插件的收集触发器

> 请替换请求正文中的 [gitlab projectId](plugins/gitlab/README-zh-CN.md#如何获取-gitlab-project-id) 和 [jira boardId](plugins/jira/README-zh-CN.md#如何获取-jira-board-id)。对于大型项目，这可能需要20分钟。 (Gitlab 10k+ commits 或 Jira 5k+ 事务)

6. 完成后，点击 *Go to grafana* (用户名: `admin`, 密码: `admin`)

### 设置 Cron job
通常情况下，我们有定期同步数据的要求。我们提供了一个叫做 `lake-cli` 的工具来满足这个要求。请在 [这里](./cmd/lake-cli/README.md) 查看 `lake-cli` 的用法。 

除此之外，如果你只想使用 Cron job，请在 [这里](./devops/sync/README.md) 查看 `docker-compose` 版本。


## 开发者设置<a id="dev-setup"></a>

### 前期准备

- <a href="https://docs.docker.com/get-docker" target="_blank">Docker</a>
- <a href="https://golang.org/doc/install" target="_blank">Golang</a>
- Make
  - Mac (Already installed)
  - Windows: [Download](http://gnuwin32.sourceforge.net/packages/make.htm)
  - Ubuntu: `sudo apt-get install build-essential`

### 如何设置开发环境
1. 进入你想安装本项目的路径，并克隆资源库

   ```sh
   git clone https://github.com/merico-dev/lake.git
   cd lake
   ```

2. 安装 go packages

    ```sh
    make install
    ```

3. 将样本配置文件复制到新的本地文件

    ```sh
    cp .env.example .env
    ```
   在`.env`文件中找到以`DB_URL`开头的那一行，把`mysql:3306`替换为`127.0.0.1:3306`
4. 启动 Docker

    > 确保在此步骤之前 Docker 正在运行。

    ```sh
    make compose
    ```

5. 运行项目

    ```sh
    make dev
    ```

6. 发送请求到 /task，创建一个 Jira 任务。这将从 Jira 收集数据

    ```
    curl -XPOST 'localhost:8080/task' \
    -H 'Content-Type: application/json' \
    -d '[[{
        "plugin": "jira",
        "options": {
            "boardId": 8
        }
    }]]'
    ```

7. 在Grafana仪表板中实现数据的可视化

    _从这里你可以看到丰富的图表，这些图表来自于收集和处理后的数据_

    - 导航到 http://localhost:3002 (用户名: `admin`, 密码: `admin`)
    - 你也可以创建/修改现有的/保存到 `Dev lake` 中的仪表板
    - 关于在Dev Lake中使用Grafana的更多信息，请看 [Grafana 文档](docs/GRAFANA.md)


## 测试<a id="tests"></a>

运行测试: `make test`

## Grafana<a id="grafana"></a>

我们使用 <a href="https://grafana.com/" target="_blank">Grafana</a> 作为可视化工具，为存储在我们数据库中的数据建立图表。可以使用SQL查询，添加面板来构建、保存和编辑自定义仪表盘。

关于配置和定制仪表盘的所有细节可以在 [Grafana 文档](docs/GRAFANA.md) 中找到。

## 贡献

[CONTRIBUTING.md](CONTRIBUTING.md)

## 需要帮助?

在 <a href="https://discord.com/invite/83rDG6ydVZ" target="_blank">Discord</a> 上给我们发消息


## FAQ<a id="faq"></a>

问：当我运行```docker-compose up -d ```时，得到这个错误: "qemu: uncaught target signal 11 (Segmentation fault) - core dumped"。如何解决这个问题？

答：Mac M1用户需要在他们的机器上下载一个特定版本的docker。你可以在这里找到它。
https://docs.docker.com/desktop/mac/apple-silicon/

