# Dev Lake

## Quick start

### Required
- [Install Golang](https://golang.org/doc/install)

### Developer Setup

```shell
git clone https://github.com/merico-dev/lake.git
cd lake
make get
cp .env.example .env
make build
make compose

```
While docker is running, in a new terminal:
```
cd lake
./lake
```

Then you can make a POST request:
```
curl --location --request POST 'localhost:8080/source' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Plugin": "Jira",
    "Options": {}
    
}'
```

### Makefile

This is like the package.json file that runs all our commands

1. You can install make 
Ubuntu: `sudo apt-get install build-essential`
Windows: `http://gnuwin32.sourceforge.net/packages/make.htm`
Mac: Comes pre installed
2. Then you can run make commands like this:
`make hello`

### How to make a commit

We use https://github.com/lintingzhen/commitizen-go to author our commits. 

Then you can run:
`make commit`

### How to run the tests

You can see a sample test in /test/example
You can run the tests with `make test`