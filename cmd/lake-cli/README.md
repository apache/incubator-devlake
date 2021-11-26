## Cmd

lake cli tools

### Install

```shell
# install from source
git clone https://github.com/merico-dev/lake.git
cd lake
go install ./cmd/lake-cli/

# go version lower than 1.17
go get -u github.com/merico-dev/lake/cmd/lake-cli

# go 1.17 or higher version
go install github.com/merico-dev/lake/cmd/lake-cli
```

### Usage

```shell
# trigger lake api
# cron schedule defined at https://pkg.go.dev/github.com/robfig/cron#hdr-Predefined_schedules
$ lake-cli api pipeline -m POST --body "{\"name\":\"sync-hourly\", \"tasks\":[[{\"plugin\":\"jira\", \"options\": {\"boardId\": 8}}]]}" --cron "@hourly"

# trigger lake api, and read request body from file
$ lake-cli api pipeline -m POST --body ./req.json

# create lake plugin (TODO)
$ lake-cli plugin init -o ./plugin --name jenkins
```
