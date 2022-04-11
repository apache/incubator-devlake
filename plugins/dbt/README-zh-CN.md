# Dbt

<div align="center">

| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |

</div>

<br>

## 概述
dbt（数据构建工具）使分析工程师能够通过简单地编写select语句来转换仓库中的数据。dbt负责将这些select语句转换为表和视图。dbt在ELT（Extract，Load，Transform）过程中起着重要作用。它不提取或加载数据，但它非常擅长转换已经加载到仓库中的数据。

## 用户安装<a id="user-setup"></a>
- 如果您计划使用本产品，你首先需要安装一些环境。

#### 需要安装的软件包<a id="user-setup-requirements"></a>
- [python3.7+](https://www.python.org/downloads/)
- [dbt-mysql](https://pypi.org/project/dbt-mysql/#configuring-your-profile)

#### 在你的终端和项目中执行或创建以下命令<a id="user-setup-commands"></a>
1.pip install dbt mysql
2.dbt init demoapp（demoapp是项目名称）
3.创建SQL转换和数据模型

## 通过Dbt转换数据
请使用原始JSON API，使用**cURL**或**Postman**等图形化API工具手动启动运行并且将以下请求发送到DevLake API端点。

```json
[
  [
    {
      "plugin": "dbt",
      "options": {
          "projectPath": "/Users/abeizn/demoapp",
          "projectName": "demoapp",
          "projectTarget": "dev",
          "selectedModels": ["my_first_dbt_model","my_second_dbt_model"],
          "projectVars": {
            "demokey1": "demovalue1",
            "demokey2": "demovalue2"
        }
      }
    }
  ]
]
```

- `projectPath`：dbt项目的绝对路径。（必选）
- `projectName`：dbt项目的名称。（必选）
- `projectTarget`：这是dbt项目将使用的默认目标分支。（可选）
- `selectedModels`：模型是select语句。模型在中定义在sql文件，通常位于模型目录中。（必选）
selectedModels接受一个或多个参数。每个参数可以是以下参数之一：
1. 包名: 运行项目中的所有模型，例如：example
2. 模型名: 运行特定的模型，例如：my_First_dbt_model
3. 模型目录的完全限定路径。

- `vars`:dbt提供了一种机制变量，用于向模型提供数据进行编译。（可选）
示例：select * from events where event_type = '{{ var("event_type") }}' ，您需要设置参数“{event_type:real_value}”的值。

### 资源：
-了解更多关于dbt的信息[在文档中](https://docs.getdbt.com/docs/introduction)
-查看[对话](https://discourse.getdbt.com/)常见问题和答案