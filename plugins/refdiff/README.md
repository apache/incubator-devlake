# RefDiff


## Summary

For development workload analysis, we often need to know how many commits have been created between 2 releases. This plugin offers the ability to calculate the commits of difference between 2 Ref(branch/tag), and the result will be stored back into database for further analysis.

## Important Note

You need to run gitextractor before the refdiff plugin. The gitextractor plugin should create records in the `refs` table in your DB before this plugin can be run.

## Configuration

This is a enrichment plugin based on Domain Layer data, no configuration needed

## How to use

In order to trigger the enrichment, you need to insert a new task into your pipeline.

1. Make sure `commits` and `refs` are collected into your database, `refs` table should contain records like following:
```
id                                            ref_type
github:GithubRepo:384111310:refs/tags/0.3.5   TAG
github:GithubRepo:384111310:refs/tags/0.3.6   TAG
github:GithubRepo:384111310:refs/tags/0.5.0   TAG
github:GithubRepo:384111310:refs/tags/v0.0.1  TAG
github:GithubRepo:384111310:refs/tags/v0.2.0  TAG
github:GithubRepo:384111310:refs/tags/v0.3.0  TAG
github:GithubRepo:384111310:refs/tags/v0.4.0  TAG
github:GithubRepo:384111310:refs/tags/v0.6.0  TAG
github:GithubRepo:384111310:refs/tags/v0.6.1  TAG
```
2. If you want to run calculateIssuesDiff, please configure GITHUB_PR_BODY_CLOSE_PATTERN in .env, you can check the example in .env.example(we have a default value, please make sure your pattern is disclosed by single quotes '')
3. If you want to run calculatePrCherryPick, please configure GITHUB_PR_TITLE_PATTERN in .env, you can check the example in .env.example(we have a default value, please make sure your pattern is disclosed by single quotes '')
4. And then, trigger a pipeline like following, you can also define sub tasks, calculateRefDiff will calculate commits between two ref, and creatRefBugStats will create a table to show bug list between two ref:
```
curl -v -XPOST http://localhost:8080/pipelines --data @- <<'JSON'
{
    "name": "test-refdiff",
    "tasks": [
        [
            {
                "plugin": "refdiff",
                "options": {
                    "repoId": "github:GithubRepo:384111310",
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

## Development

This plugin depends on `libgit2`, you need to install version 1.3.0 in order to run and debug this plugin on your local
machine.

### Ubuntu

```
apt install cmake
git clone https://github.com/libgit2/libgit2.git
git checkout v1.3.0
mkdir build
cd build
cmake ..
make
make install
```

### MacOs

```
brew install cmake
git clone https://github.com/libgit2/libgit2.git
cd libgit2
git checkout v1.3.0
mkdir build
cd build
cmake ..
make
make install
```

Troubleshooting (MacOS)

Q: I got an error saying: `pkg-config: exec: "pkg-config": executable file not found in $PATH`

A:

1. Make sure you have pkg-config installed:

  `brew install pkg-config`

2. Make sure your pkg config path covers the installation: 

  `export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/local/lib:/usr/local/lib/pkgconfig`
