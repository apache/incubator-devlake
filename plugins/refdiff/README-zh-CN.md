# RefDiff 插件


| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |


## 概述

在分析开发工作产生代码量时，常常需要知道两个版本之间产生了多少个 commit。本插件基于数据库中存储的 commits 父子关系信息，提供了计算两个 ref(branch/tag) 之间相差 commits 列表的能力。计算的结果回存于数据库中，方便后续的交叉分析。


## 配置

本插件基于领域层数据进行数据增强，无需额外配置。

## 如何使用

为了触发数据增强，您需要在 Pipeline 中加入一个新的任务

1. 确保 `commits` 表和 `refs` 的数据已经正确收集，`refs` 表应含有类似下面的数据:
```
id                                                  ref_type
github:GithubRepository:384111310:refs/tags/0.3.5   TAG
github:GithubRepository:384111310:refs/tags/0.3.6   TAG
github:GithubRepository:384111310:refs/tags/0.5.0   TAG
github:GithubRepository:384111310:refs/tags/v0.0.1  TAG
github:GithubRepository:384111310:refs/tags/v0.2.0  TAG
github:GithubRepository:384111310:refs/tags/v0.3.0  TAG
github:GithubRepository:384111310:refs/tags/v0.4.0  TAG
github:GithubRepository:384111310:refs/tags/v0.6.0  TAG
github:GithubRepository:384111310:refs/tags/v0.6.1  TAG
```
2. 如果您想要使用calculatePrCherryPick，请在.env文件中配置GITHUB_PR_TITLE_PATTERN，可以在.env.example中查看示例
3. 然后，通过类似下面的命令触发一个 pipeline，在tasks中，可以定义想要执行的任务，calculateRefDiff可以计算新老版本间的差了多少个 commits，creatRefBugStats可以生成新老版本间的issue列表
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
                       { "newRef": "refs/tags/v0.6.0", "oldRef": "refs/tags/0.5.0" },
                       { "newRef": "refs/tags/0.5.0", "oldRef": "refs/tags/0.4.0" }
                    ],
                    "tasks": [
                        "calculateCommitsDiff",
                        "calculateIssuesDiff",
                        "calculatePrCherryPick",
                    ]
                }
            }
        ]
    ]
}
JSON
```
