# RefDiff 插件


| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |


## 概述

在分析开发工作产生代码量时，常常需要知道两个版本之间产生了多少个 commit。本插件基于数据库中存储的 commits 父子关系信息，提供了计算两个 ref(branch/tag) 之间相差 commits 列表的能力。计算的结果回存于数据库中，方便后续的交叉分析。


## 配置

本插件基于领域层数据进行数据增强，无需额外配置。

## 如何使用

为了触发数据增强，您需要在 Pipeline 中加入一个新的任务

```
curl -v -XPOST http://localhost:8080/pipelines --data @- <<'JSON'
{
    "name": "test-refdiff",
    "tasks": [
        [
            {
                "plugin": "refdiff",
                "options": {
                    "repoId": "github:GithubRepository:384111310",
                    "pairs": [
                       { "newRef": "v0.6.0", "oldRef": "0.5.0" },
                       { "newRef": "0.5.0", "oldRef": "0.4.0" }
                    ]
                }
            }
        ]
    ]
}
JSON
```
