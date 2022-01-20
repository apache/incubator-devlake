# RefDiff


| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |


## Summary

For development workload analysis, we often need to know how many commits have been created between 2 releases. This plugin offers the ability to calculate the commits of difference between 2 Ref(branch/tag), and the result will be stored back into database for further analysis.


## Configuration

This is a enrichment plugin based on Domain Layer data, no configuration needed

## How to use

为了触发数据增强，您需要在 Pipeline 中加入一个新的任务
In order to trigger the enrichment, you need to insert a new task into your pipeline

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
