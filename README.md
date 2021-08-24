# Dev Lake

## Quick start

### Required
- [Install Golang](https://golang.org/doc/install)

### Developer Setup
```shell
# clone lake repository
$ git clone https://github.com/merico-dev/lake.git
# enter lake directory
$ cd lake
# get packages
$ go get
# build 
$ go build
# copy .env.example to .env
$ cp .env.example .env
$ ./lake
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
