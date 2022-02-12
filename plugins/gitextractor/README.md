# Git Repo Extractor

## Summary
This plugin extract commits and references from a remote or local git repository. It then saves the data into the database or csv files.

## Steps to make this plugin work

1. Use the Git repo extractor to retrieve commit-and-branch-related data from your repo
2. Use the GitHub plugin to retrieve Github-issue-and-pr-related data from your repo. NOTE: you can run only one the issue collection stage as described in the Github Plugin README.
3. Use the RefDiff plugin to calculate version diff, which will be stored in refs_commits_diffs.

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
- `user`: optional, for cloning private repository using HTTP/HTTPS
- `password`: optional, for cloning private repository using HTTP/HTTPS
- `privateKey`: optional, for SSH cloning, base64 encoded `PEM` file
- `passphrase`: optional, passphrase for the private key
