# Contributing to Lake

👍🎉 First off, thanks for taking the time to contribute! 🎉👍

The following is a set of guidelines for contributing to Lake. These are mostly guidelines, not rules. 
Use your best judgment, and feel free to propose changes to this document in a pull request.

## How Can I Contribute?

1. Reporting bugs by filling out the required issue template and labeling the new issue as 'bug'.

2. Suggesting enhancements.

If you intend to change the public API, or make any non-trivial changes to the implementation, we recommend filing an issue. 
This lets us reach an agreement on your proposal before you put significant effort into it.

If you’re only fixing a bug, it’s fine to submit a pull request right away but we still recommend to file an issue detailing what you’re fixing. 
This is helpful in case we don’t accept that specific fix but want to keep track of the issue.

## How to run this project with docker-compose

1. Install docker

https://docs.docker.com/get-docker/

2. Run docker-compose

`docker-compose up -d`

3. Stop all containers

`docker-compose down`

## How to run the project in dev mode if you want to work on the backend (Frontend excluded)

1. Clone the repository

```
git clone https://github.com/apache/incubator-devlake
```

2. Init your config file

`cp .env.example .env`

3. Install the correct version of go 

here: https://go.dev/dl/

4. Install all dependencies

```
go version
go get
brew install pkg-config
brew install cmake
git clone https://github.com/libgit2/libgit2.git
cd libgit2
git checkout v1.3.0
mkdir build
cd build
cmake ..
make
make install
cd ../..
```

5. Install and configure mysql

```
brew install mysql
mysql -uroot -p # password is "root" by default
```

Then from the mysql terminal interface:

```
CREATE USER 'merico'@'%' IDENTIFIED BY 'merico';
GRANT ALL PRIVILEGES ON *.* TO 'merico'@'%';
CREATE DATABASE IF NOT EXISTS lake;
```

Note: You can set anything you want for user/pass/dbname in your .env file

6. Build all plugins

`make all`

You should see an output similar to this: 

```
➜  lake git:(main) ✗ make all
Building plugin ae to bin/plugins/ae/ae.so
Building plugin dbt to bin/plugins/dbt/dbt.so
Building plugin feishu to bin/plugins/feishu/feishu.so
Building plugin gitextractor to bin/plugins/gitextractor/gitextractor.so
Building plugin github to bin/plugins/github/github.so
Building plugin gitlab to bin/plugins/gitlab/gitlab.so
Building plugin jenkins to bin/plugins/jenkins/jenkins.so
Building plugin jira to bin/plugins/jira/jira.so
Building plugin refdiff to bin/plugins/refdiff/refdiff.so
Building plugin tapd to bin/plugins/tapd/tapd.so
go build -ldflags "-X 'github.com/apache/incubator-devlake/version.Version=@5ed73bdb'" -o bin/lake
go build -ldflags "-X 'github.com/apache/incubator-devlake/version.Version=@5ed73bdb'" -o bin/lake-worker ./worker/
```

7. Run the backend

`make run`

## Style guides

### Git Commit message

We follow the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) guidelines.

#### Commit tool

We use https://github.com/lintingzhen/commitizen-go to author our commits.

```sh
make commit
```

```
> lake@1.0.0 commit /home/code/merico-dev/lake
> cz

cz-cli@4.2.4, cz-conventional-changelog@3.3.0

? Select the type of change that you're committing: (Use arrow keys)
> feat:     A new feature
  fix:      A bug fix
  docs:     Documentation only changes
  style:    Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
  refactor: A code change that neither fixes a bug nor adds a feature
  perf:     A code change that improves performance
  test:     Adding missing tests or correcting existing tests
(Move up and down to reveal more choices)
? What is the scope of this change (e.g. component or file name): (press enter to skip)
? Write a short, imperative tense description of the change (max 93 chars):
 (23) add commit message tool
? Provide a longer description of the change: (press enter to skip)

? Are there any breaking changes? No
? Does this change affect any open issues? No
[chore/commit_message dc34f57] chore: add commit message tool
 5 files changed, 585 insertions(+), 4 deletions(-)
```
