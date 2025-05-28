> **Think:**
> 我需要将这个DevLake开发环境部署文档翻译成中文。我将保持原文档的结构，但翻译所有内容到中文。
>
> 文档主要包括以下部分：
> 1. 环境要求
> 2. 如何设置开发环境
> 3. 克隆仓库
> 4. 安装插件依赖
> 5. 配置环境变量
> 6. 启动MySQL和Grafana容器
> 7. 运行devlake和config-ui
>
> 我会逐部分翻译，保持专有名词如Docker、Golang、MySQL、Grafana等不变，同时保持命令行和代码部分不变。

# DevLake开发环境部署指南

## 环境要求
- Docker v19.03.10+
- Golang v1.19+
- GNU Make
    - Mac (预装)
    - Windows: [下载](http://gnuwin32.sourceforge.net/packages/make.htm)
    - Ubuntu: `sudo apt-get install build-essential libssl-dev`

## 如何设置开发环境
以下指南将介绍如何在开发模式下运行DevLake的前端(config-ui)和后端。

### 克隆代码仓库
导航到您想要安装此项目的位置并克隆仓库：

```bash
git clone https://github.com/apache/incubator-devlake.git
cd incubator-devlake
```

### 安装插件依赖

RefDiff 插件：
安装Go包
```bash
cd backend
go get
cd ..
```

### 配置环境文件
复制示例配置文件到本地新文件：

```bash
cp env.example .env
```

更新`.env`文件中的以下变量：

- `DB_URL`: 将`mysql:3306`替换为`127.0.0.1:3306`
- `DISABLED_REMOTE_PLUGINS`: 设置为`True`

### 启动MySQL和Grafana容器

确保在此步骤之前Docker守护进程正在运行。

> Grafana需要重新build镜像，然后在docker-compose.datasources.yml中更改image为`image: grafana:latest`

```bash
docker-compose -f docker-compose-dev.yml up -d mysql grafana
```

### 运行开发模式
在两个单独的终端中以开发模式运行devlake和config-ui：

```bash
# 安装poetry，按照指南操作：https://python-poetry.org/docs/#installation
# 运行devlake，这里只用了q dev插件
DEVLAKE_PLUGINS=q_dev nohup make dev &
# 运行config-ui
make configure-dev
```

常见错误请参见故障排除文档。

Config UI 运行在 localhost:4000