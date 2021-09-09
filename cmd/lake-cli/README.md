## Cmd

lake cli tools

### Install

```shell
# go version lower than 1.17
go get -u github.com/merico-dev/lake/cmd/lake-cli

# go 1.17 or higner version
go install github.com/merico-dev/lake/cmd/lake-cli
```

### Usage

```shell
# trigger lake api
$ ./lake-cli api task -m POST --body "[{'plugin':'jira', 'options': {'boardId': 8}}]" --cron "@every 5s"

# create lake plugin (TODO)
$ ./lake-cli plugin init -o ./plugin --name jenkins
```
