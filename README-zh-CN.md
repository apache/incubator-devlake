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
1. 在一个个性化的、统一的视图中可视化和分析你的整个SDLC过程。
2. 提供统一的标准化的度量体系和分析方法，帮助你分析团队研发效能，提高交付速度和质量。

### Dev Lake 可以完成什么?
1. 收集和关联不同来源的数据（Jira、Gitlab、Github、Jenkins等），打破数据孤岛。
2. 提供行业标准指标来识别工程问题，比如交付周期过长、Bug太多等。
3. 高度定制化的图表、指标和仪表盘，用户根据自己的需求分析数据，发现洞见。

<br>

## 目录
<table>
    <tr>
        <td>目录</td>
        <td>子目录</td>
        <td>描述</td>
        <td>文档链接</td>
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
        <td rowspan="3">贡献</td>
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
Jira | 概述，数据和指标，配置，API | [Link](https://github.com/merico-dev/lake/blob/main/plugins/jira/README-zh-CN.md)
Gitlab | 概述，数据和指标，配置，API | [Link](https://github.com/merico-dev/lake/blob/main/plugins/gitlab/README-zh-CN.md)
Jenkins | 概述，数据和指标，配置，API | [Link](https://github.com/merico-dev/lake/blob/main/plugins/jenkins/README-zh-CN.md)
GitHub | 概述，数据和指标，配置，API | [Link](https://github.com/merico-dev/lake/blob/main/plugins/github/README-zh-CN.md)

<br>

## 安装手册
一共有 3 种方式来安装 Dev Lake。


### 用户安装<a id="user-setup"></a>

**注意：如果你只打算运行 Dev Lake，你只需要阅读这一小节**<br>
**注意：写成 `这样` 的命令需要在你的终端中运行**

### 需要安装的软件包<a id="user-setup-requirements"></a>

- [Docker](https://docs.docker.com/get-docker)
- [docker-compose](https://docs.docker.com/compose/install/)

**NOTE:** 安装完 Docker 后，你可能需要运行 Docker 应用程序并重新启动你的终端

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
   >-> <a href="https://github.com/merico-dev/lake/blob/main/plugins/jira/README-zh-CN.md" target="_blank">Jira</a>
   >-> <a href="https://github.com/merico-dev/lake/blob/main/plugins/gitlab/README-zh-CN.md" target="_blank">GitLab</a>
   >-> <a href="https://github.com/merico-dev/lake/blob/main/plugins/jenkins/README-zh-CN.md" target="_blank">Jenkins</a>
   >-> <a href="https://github.com/merico-dev/lake/blob/main/plugins/github/README-zh-CN.md" target="_blank">GitHub</a>


4. 访问 `localhost:4000/triggers`，触发数据收集

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
   >     ],
   >     [
   >       {
   >         "plugin": "github-domain",
   >         "options": {}
   >       }
   >     ]
   >   ]
   >   ```


5. 完成后，点击 *Go to grafana* (用户名: `admin`, 密码: `admin`)。当数据收集完成后，该按钮将显示在触发收集页面。

### 设置 Cron job
通常情况下，我们有定期同步数据的要求。我们提供了一个叫做 `lake-cli` 的工具来满足这个要求。请在 [这里](./cmd/lake-cli/README.md) 查看 `lake-cli` 的用法。

除此之外，如果你只想使用 Cron job，请在 [这里](./devops/sync/README.md) 查看 `docker-compose` 版本。


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

4. 启动 Docker

    > 确保在此步骤之前 Docker 正在运行。

    ```sh
    make compose
    ```

5. 运行项目

    ```sh
    make dev
    ```

6. 访问 `localhost:4000` 来设置 Dev Lake 的配置文件
   >- 在 "Integration"页面上找到到所需的插件页面
   >- 你需要为你打算使用的插件输入必要的信息
   >- 请参考以下内容，以了解如何配置每个插件的更多细节
   >-> <a href="https://github.com/merico-dev/lake/blob/main/plugins/jira/README-zh-CN.md" target="_blank">Jira</a>
   >-> <a href="https://github.com/merico-dev/lake/blob/main/plugins/gitlab/README-zh-CN.md" target="_blank">GitLab</a>
   >-> <a href="https://github.com/merico-dev/lake/blob/main/plugins/jenkins/README-zh-CN.md" target="_blank">Jenkins</a>
   >-> <a href="https://github.com/merico-dev/lake/blob/main/plugins/github/README-zh-CN.md" target="_blank">GitHub</a>


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
   >     ],
   >     [
   >       {
   >         "plugin": "github-domain",
   >         "options": {}
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

### 云端安装<a id="cloud-setup"></a>
如果你想在云端安装Dev Lake，你可以使用 Tin 来进行. [查看详细信息](https://github.com/merico-dev/lake/wiki/How-to-Set-Up-Dev-Lake-with-Tin-zh-CN)

**声明:**
> 对于使用 Tin 在云端托管 Dev Lake 的用户，设置密码来保护实例下配置信息的安全是至关重要的。Dev Lake作为一个自我托管的产品，部分是为了确保用户对数据有完全的保护和所有权，对于 Tin 托管来说也是如此，这个风险点需要由终端用户来消除。

## 测试<a id="tests"></a>

运行测试:

```sh
make test
```

<br>

## 贡献
本节列出了所有的文件，以帮助你快速为 repo 做出贡献。

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
