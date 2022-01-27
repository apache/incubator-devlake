# RefDiff


| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |


## Summary

For development workload analysis, we often need to know how many commits have been created between 2 releases. This plugin offers the ability to calculate the commits of difference between 2 Ref(branch/tag), and the result will be stored back into database for further analysis.


## Configuration

This is a enrichment plugin based on Domain Layer data, no configuration needed

## How to use

In order to trigger the enrichment, you need to insert a new task into your pipeline.

1. Make sure `commits` and `refs` are collected into your database, `refs` table should contain records like following:
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
2. And then, trigger a pipeline like following:
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
                    ]
                }
            }
        ]
    ]
}
JSON
```

## Install `libgit2`

### Ubuntu

```
apt install cmake
git clone https://github.com/libgit2/libgit2.git
git checkout v1.3.0
make build
cd build
cmake ..
make
make install
ldconfig
```

### MacOs

```
brew install cmake
git clone https://github.com/libgit2/libgit2.git
git checkout v1.3.0
make build
cd build
cmake ..
make
make install
ldconfig
```
