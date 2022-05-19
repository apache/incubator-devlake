# Git Repo Extractor

## Summary
This plugin extract commits and references from a remote or local git repository. It then saves the data into the database or csv files.

## Steps to make this plugin work

1. Use the Git repo extractor to retrieve commit-and-branch-related data from your repo
2. Use the GitHub plugin to retrieve Github-issue-and-pr-related data from your repo. NOTE: you can run only one the issue collection stage as described in the Github Plugin README.
3. Use the [RefDiff](../refdiff) plugin to calculate version diff, which will be stored in `refs_commits_diffs` table.

## Sample Request

```
curl --location --request POST 'localhost:8080/pipelines' \
--header 'Content-Type: application/json' \
--data-raw '
{
    "name": "git repo extractor",
    "tasks": [
        [
            {
                "Plugin": "gitextractor",
                "Options": {
                    "url": "https://github.com/apache/incubator-devlake.git",
                    "repoId": "github:GithubRepo:384111310"
                }
            }
        ]
    ]
}
'
```
- `url`: the location of the git repository. It should start with `http`/`https` for remote git repository or `/` for a local one.
- `repoId`: column `id` of  `repos`.
- `proxy`: optional, http proxy, e.g. `http://your-proxy-server.com:1080`.
- `user`: optional, for cloning private repository using HTTP/HTTPS
- `password`: optional, for cloning private repository using HTTP/HTTPS
- `privateKey`: optional, for SSH cloning, base64 encoded `PEM` file
- `passphrase`: optional, passphrase for the private key


## Standalone Mode

You call also run this plugin in a standalone mode without any DevLake service running using the following command:

```
go run plugins/gitextractor/main.go -url https://github.com/apache/incubator-devlake.git -id github:GithubRepo:384111310 -db "merico:merico@tcp(127.0.0.1:3306)/lake?charset=utf8mb4&parseTime=True"
```

For more options (e.g., saving to a csv file instead of a db), please read `plugins/gitextractor/main.go`.

## Development

This plugin depends on `libgit2`, you need to install version 1.3.0 in order to run and debug this plugin on your local
machine. [Click here](../refdiff#Development) for a brief guide.
