# Git Repo Extractor

## 概述
This plugin extract commits and references from a remote or local git repository. It then save the data into database or csv files.



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
                    "url": "https://github.com/merico-dev/lake.git",
                    "repoId": "github:GithubRepos:384111310"
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
