<div align="center">
<br />
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

# Dev Lake
<p>
    <b>
     <!Software development workflow analysis for free> 
    </b>
  </p>
  <p>

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/discord/844603288082186240.svg?style=flat?label=&logo=discord&logoColor=ffffff&color=747df7&labelColor=454545)](https://discord.gg/83rDG6ydVZ)
![badge](https://github.com/merico-dev/lake/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/merico-dev/lake)](https://goreportcard.com/report/github.com/merico-dev/lake)


| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |

</div>
<br>
<div align="left">

### 什么是 Dev Lake？
Dev Lake 将你所有的 DevOps 数据以实用、个性化、可扩展的视图呈现。通过我们免费的开源产品，从不断增加的开发者工具列表中收集、分析和可视化数据。

Dev Lake 对于希望**更好地了解其开发数据的管理者**来说是最激动人心的，除此以外，它对于任何希望**通过数据驱动提升自身实践的开发者**来说都是有用的。有了Dev Lake，你可以向你的程序提出任何问题，只要连接和查询。


#### 一键体验 Dev Lake

<table>
  <tr>
    <td valign="middle"><a href="#user-setup">在本地运行</a></td>
    <td valign="middle">
      <a valign="middle" href="https://www.teamcode.com/tin/clone?applicationId=259777118600769536">
        <img
          src="https://static01.teamcode.com/badge/teamcode-badge-run-in-cloud-cn.svg"
          width="140px"
          alt="Teamcode" valign="middle"
        />
      </a>
      <a valign="middle"
        href="https://github.com/merico-dev/lake/wiki/How-to-Set-Up-Dev-Lake-with-Tin-zh-CN"><span valign="middle">查看手册</span>
      </a>
    </td>
  </tr>
</table>

<br>

<div align="left">
<img src="https://user-images.githubusercontent.com/2908155/130271622-827c4ffa-d812-4843-b09d-ea1338b7e6e5.png" width="100%" alt="Dev Lake Grafana Dashboard" style="border-radius:15px;" />
<p align="center">数据面板截图</p><br>
<img src="https://user-images.githubusercontent.com/14050754/142356580-40637a30-5578-48ed-8e4a-128cd0738e3e.png" width="100%" alt="User Flow" style="border-radius:15px;"/>
<p align="center">用户使用流程</p><br>



### 为什么选择 Dev Lake？
1. 全面了解软件研发生命周期，挖掘工作流瓶颈
2. 及时回顾团队迭代表现，快速反馈，敏捷调整
3. 快速搭建场景化数据仪表盘，下钻分析洞察问题根因

### Dev Lake 可以完成什么?
1. 归集 DevOps 全流程效能数据
2. 同类工具共用抽象层，输出标准化效能数据
3. 内置20+效能指标与下钻分析能力
4. 支持自定义 SQL 分析及拖拽搭建场景化数据视图
5. 灵活架构与插件设计，支持快速接入新数据源

### 查看 Demo
[点击这里](https://app-259373083972538368-3002.ars.teamcode.com/d/0Rjxknc7z/demo-homepage?orgId=1) 查看 Demo. Demo里呈现的数据来自此仓库。<br>
用户名/密码: test/test


<br>

## 目录
<table>
    <tr>
        <td><b>目录</b></td>
        <td><b>子目录</b></td>
        <td><b>描述</b></td>
        <td><b>文档链接</b></td>
    </tr>
    <tr>
        <td>数据源</td>
        <td>当前支持的数据源</td>
        <td>链接到具体的插件使用和细节</td>
        <td><a href="#data-source-plugins">查看本节</a></td>
    </tr>
    <tr>
        <td rowspan="3">安装手册</td>
        <td>用户安装</td>
        <td>以用户身份运行项目的步骤</td>
        <td><a href="#user-setup">查看本节</a></td>
    </tr>
    <tr>
        <td>开发者安装</td>
        <td>如何设置开发环境</td>
        <td><a href="#dev-setup">查看本节</a></td>
    </tr>
    <tr>
        <td>云端安装</td>
        <td>使用 Tin 进行云端安装</td>
        <td><a href="#cloud-setup">查看本节</a></td>
    </tr>
   <tr>
        <td>测试</td>
        <td>测试</td>
        <td>运行测试的命令</td>
        <td><a href="#tests">查看本节</a></td>
    </tr>
    <tr>
        <td rowspan="4">贡献</td>
        <td>了解 DevLake 的架构</td>
        <td>查看系统架构图</td>
        <td><a href="#architecture">查看本节</a></td>
    </tr>
    <tr>
        <td>添加一个插件</td>
        <td>如何制作自己的插件的详细信息</td>
        <td><a href="#plugin">查看本节</a></td>
    </tr>
   <tr>
        <td>添加新的指标</td>
        <td>如何给插件添加指标</td>
        <td><a href="#metrics">查看本节</a></td>
    </tr>
    <tr>
        <td>代码规范</td>
        <td>如何进行贡献</td>
        <td><a href="#contributing">查看本节</a></td>
    </tr>
    <tr>
        <td rowspan="4">用户使用手册，帮助文档等</td>
        <td>Grafana</td>
        <td>如何将数据进行可视化</td>
        <td><a href="#grafana">查看本节</a></td>
    </tr>
    <tr>
        <td>帮助</td>
        <td>在 Discord 上联系我们</td>
        <td><a href="#help">查看本节</a></td>
    </tr>
    <tr>
        <td>FAQ</td>
        <td>常见问题</td>
        <td><a href="#faq">查看本节</a></td>
    </tr>
    <tr>
        <td>许可证</td>
        <td>Dev Lake 许可证</td>
        <td><a href="#license">查看本节</a></td>
    </tr>
</table>

<br>

## 我们目前支持的数据源<a id="data-source-plugins"></a>

下面是一个 _数据源插件（data source plugins）_ 的列表，用于收集和处理特定来源的数据。每个插件都有一个 `README.md` 文件，包含基本设置、故障排除和指标信息。

关于建立一个新的 _data source plugins_ 的更多信息，请参见[添加一个插件](plugins/README-zh-CN.md)。

目录 | 内容 | 文档
------------ | ------------- | -------------
Jira | 概述，数据和指标，配置，API | [Link](plugins/jira/README-zh-CN.md) 
Gitlab | 概述，数据和指标，配置，API | [Link](plugins/gitlab/README-zh-CN.md) 
Jenkins | 概述，数据和指标，配置，API | [Link](plugins/jenkins/README-zh-CN.md) 
GitHub | 概述，数据和指标，配置，API | [Link](plugins/github/README-zh-CN.md)

<br>

## 安装手册
一共有 3 种方式来安装 Dev Lake：用户安装，开发者安装和云端安装。


### 用户安装<a id="user-setup"></a>

- 如果你只打算运行 Dev Lake，你只需要阅读这一小节<br>
- 写成 `这样` 的命令需要在你的终端中运行

### 需要安装的软件包<a id="user-setup-requirements"></a>

- [Docker](https://docs.docker.com/get-docker)
- [docker-compose](https://docs.docker.com/compose/install/)

注：安装完 Docker 后，你可能需要运行 Docker 应用程序并重新启动你的终端

#### 在你的终端中运行以下命令<a id="user-setup-commands"></a>

1. 克隆仓库。

   ```sh
   git clone https://github.com/merico-dev/lake.git devlake
   cd devlake
   cp .env.example .env
   ```
2. 启动 Docker，然后运行 `docker-compose up -d` 启动服务。

3. 访问 `localhost:4000` 来设置 Dev Lake 的配置文件
   >- 在 "Integration"页面上找到到所需的插件页面
   >- 你需要为你打算使用的插件输入必要的信息
   >- 请参考以下内容，以了解如何配置每个插件的更多细节
   >-> <a href="plugins/jira/README-zh-CN.md" target="_blank">Jira</a>
   >-> <a href="plugins/gitlab/README-zh-CN.md" target="_blank">GitLab</a>
   >-> <a href="plugins/jenkins/README-zh-CN.md" target="_blank">Jenkins</a> 
   >-> <a href="plugins/github/README-zh-CN.md" target="_blank">GitHub</a>
   
   >- 提交表单，通过点击每个表单页面上的**Save Connection**按钮来更新数值。

   >- `devlake`需要一段时间才能完全启动。如果`config-ui`提示 API 无法访问，请等待几秒钟并尝试刷新页面。
   >- 如果想收集一个 Repo 进行快速预览，请在**数据集成/Github**页面提供一个 Github 的个人 Token。

4. 访问 `localhost:4000/triggers`，触发数据收集

> 请参考这篇Wiki [How to trigger data collection](https://github.com/merico-dev/lake/wiki/How-to-use-the-triggers-page)。数据收集可能需要一段时间，取决于你想收集的数据量。

> - 如果要收集这个 repo 以进行，你可以使用以下 JSON
   >   ```json
   >   [
   >     [
   >       {
   >         "Plugin": "github",
   >         "Options": {
   >           "repositoryName": "lake",
   >           "owner": "merico-dev"
   >         }
   >       }
   >     ]
   >   ]
   >   ```


5. 完成后，点击 *Go to grafana* (用户名: `admin`, 密码: `admin`)。当数据收集完成后，该按钮将显示在触发收集页面。

#### 设置 Cron job
为了定期同步数据，我们提供了[`lake-cli`](./cmd/lake-cli/README.md)以方便发送数据收集请求，我们同时提供了[cron job](./devops/sync/README.md)以定期触发 cli 工具。

<br>

****

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

2. 安装 go packages

    ```sh
    make install
    ```

3. 将样本配置文件复制到新的本地文件

    ```sh
    cp .env.example .env
    ```
   在`.env`文件中找到以`DB_URL`开头的那一行，把`mysql:3306`替换为`127.0.0.1:3306`

4. 启动 MySQL 和 Grafana

    > 确保在此步骤之前 Docker 正在运行。

    ```sh
    docker-compose up mysql grafana
    ```

5. 在 2 个终端种分别以开发者模式运行 lake 和 config UI:

    ```sh
    # run lake
    make dev
    # run config UI
    make configure-dev
    ```

6. 访问 config-ui `localhost:4000` 来配置 Dev Lake 数据源
   >- 在 "Integration"页面上找到到所需的插件页面
   >- 你需要为你打算使用的插件输入必要的信息
   >- 请参考以下内容，以了解如何配置每个插件的更多细节
   >-> <a href="plugins/jira/README-zh-CN.md" target="_blank">Jira</a>
   >-> <a href="plugins/gitlab/README-zh-CN.md" target="_blank">GitLab</a>
   >-> <a href="plugins/jenkins/README-zh-CN.md" target="_blank">Jenkins</a> 
   >-> <a href="plugins/github/README-zh-CN.md" target="_blank">GitHub</a>


7. 访问 `localhost:4000/triggers`，触发数据收集

    > 请参考这篇Wiki [How to trigger data collection](https://github.com/merico-dev/lake/wiki/How-to-use-the-triggers-page)。对于大型项目，这可能需要20分钟。 (Gitlab 10k+ commits 或 Jira 5k+ 事务)

    > - 如果要收集这个 repo 以进行，你可以使用以下 JSON
   >   ```json
   >   [
   >     [
   >       {
   >         "Plugin": "github",
   >         "Options": {
   >           "repositoryName": "lake",
   >           "owner": "merico-dev"
   >         }
   >       }
   >     ]
   >   ]
   >   ```


8. 在Grafana仪表板中实现数据的可视化

    _从这里你可以看到丰富的图表，这些图表来自于收集和处理后的数据_

    - 导航到 http://localhost:3002 (用户名: `admin`, 密码: `admin`)
    - 你也可以创建/修改现有的/保存到 `Dev lake` 中的仪表板
    - 关于在Dev Lake中使用Grafana的更多信息，请看 [Grafana 文档](docs/GRAFANA.md)

<br>

****

<br>

### 云端安装<a id="cloud-setup"></a>
如果你想在云端安装Dev Lake，你可以使用 Tin 来进行. [查看详细信息](https://github.com/merico-dev/lake/wiki/How-to-Set-Up-Dev-Lake-with-Tin-zh-CN)

**声明:** 
> 对于使用 Tin 在云端托管 Dev Lake 的用户，设置密码来保护实例下配置信息的安全是至关重要的。Dev Lake作为一个自我托管的产品，部分是为了确保用户对数据有完全的保护和所有权，对于 Tin 托管来说也是如此，这个风险点需要由终端用户来消除。

<br>

## 测试<a id="tests"></a>

运行测试: 

```sh
make test
```

<br>

## 贡献
本节列出了所有的文件，以帮助你快速为 repo 做出贡献。

### 了解 DevLake 的架构<a id="architecture"></a>
![devlake-architecture](https://user-images.githubusercontent.com/14050754/143292041-a4839bf1-ca46-462d-96da-2381c8aa0fed.png)
<p align="center">架构图</p>

### 添加一个插件<a id="plugin"></a>

[plugins/README.md](/plugins/README.md)

### 添加新的指标<a id="metrics"></a>

[plugins/HOW-TO-ADD-METRICS.md](/plugins/HOW-TO-ADD-METRICS.md)

### 代码规范<a id="contributing"></a>

[CONTRIBUTING.md](CONTRIBUTING.md)

<br>


## 用户使用手册，帮助文档及其他
### Grafana<a id="grafana"></a>

我们使用 <a href="https://grafana.com/" target="_blank">Grafana</a> 作为可视化工具，为存储在我们数据库中的数据建立图表。可以使用SQL查询，添加面板来构建、保存和编辑自定义仪表盘。

关于配置和定制仪表盘的所有细节可以在 [Grafana 文档](docs/GRAFANA.md) 中找到。


### 需要帮助?

在 <a href="https://discord.com/invite/83rDG6ydVZ" target="_blank">Discord</a> 上给我们发消息


### FAQ<a id="faq"></a>

问：当我运行```docker-compose up -d ```时，得到这个错误: "qemu: uncaught target signal 11 (Segmentation fault) - core dumped"。如何解决这个问题？

答：Mac M1用户需要在他们的机器上下载一个特定版本的docker。你可以在这里找到它。
https://docs.docker.com/desktop/mac/apple-silicon/


### License<a id="license"></a>

此项目的许可证为 Apache License 2.0 - 查看 [`许可证`](LICENSE) 详情。
