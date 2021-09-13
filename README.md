<br />
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

# Dev Lake
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/discord/844603288082186240.svg?style=flat?label=&logo=discord&logoColor=ffffff&color=747df7&labelColor=454545)](https://discord.gg/83rDG6ydVZ)
![badge](https://github.com/merico-dev/lake/actions/workflows/main.yml/badge.svg)

### What is Dev Lake?

Dev Lake is the one-stop solution that _**integrates, analyzes, and visualizes**_ software development data throughout the _**software development life cycle (SDLC)**_ for engineering teams.

<img src="https://user-images.githubusercontent.com/2908155/130271622-827c4ffa-d812-4843-b09d-ea1338b7e6e5.png" width="100%" alt="Dev Lake Grafana Dashboard" />

### Why choose Dev Lake?

1.  Supports various data sources (<a href="https://gitlab.com/" target="_blank">Gitlab</a>, <a href="https://www.atlassian.com/software/jira" target="_blank">Jira</a>) and more are being added all the time
2.  Relevant, customizable data metrics ready to view as visual charts
7.  Easily build and view new charts and dashboards with <a href="https://grafana.com/" target="_blank">Grafana</a>
4.  Easy-to-setup via <a href="https://docs.docker.com/desktop/" target="_blank">Docker</a>
5.  Extensible plugin system to add your own data collectors
6.  Designed to process enterprise-scale data

## Contents

Section | Description | Documentation Link
:------------ | :------------- | :-------------
Data Sources | Links to specific plugin usage & details | [View Section](#data-source-plugins)
User Setup | Steps to run the project as a user | [View Section](#user-setup) 
Developer Setup | How to setup dev environment | [View Section](#dev-setup)
Tests | Commands for running tests | [View Section](#tests)
Grafana | How to visualize the data | [View Section](#grafana)
Build a Plugin | Details on how to make your own | [Link](plugins/README.md) 
Add Plugin Metrics | Guide to adding plugin metrics | [Link](plugins/HOW-TO-ADD-METRICS.md) 
Contributing | How to contribute to this repo | [Link](CONTRIBUTING.md)


## Data Sources We Currently Support<a id="data-source-plugins"></a>

Below is a list of _data source plugins_ used to collect & enrich data from specific sources. Each have a `README.md` file with basic setup, troubleshooting and metrics info.

For more information on building a new _data source plugin_ see [Build a Plugin](plugins/README.md).

Section | Section Info | Docs
------------ | ------------- | -------------
Jira | Metrics, Generating API Token, Find Project/Board ID | [Link](plugins/jira/README.md) 
Gitlab | Metrics, Generating API Token | [Link](plugins/gitlab/README.md) 
Jenkins | Metrics, Generating API Token | [Link](plugins/jenkins/README.md) 


## User setup<a id="user-setup"></a>

**NOTE: If you only plan to run the product, this is the only section you should need**
**NOTE: Commands written `like this` are to be run in your terminal**

### Required Packages to Install<a id="user-setup-requirements"></a>

- [Docker](https://docs.docker.com/get-docker)
- [docker-compose](https://docs.docker.com/compose/install/)

**NOTE:** After installing docker, you may need to run the docker application and restart your terminal

### Commands to run in your terminal<a id="user-setup-commands"></a>

1. Create a directory and download files

   ```sh
   git clone https://github.com/merico-dev/lake.git devlake
   cd devlake
   git checkout go-main
   cp .env.example .env
   ```

2. Open `.env` file with your editor, fill the values in `Jira`, `Gitlab` and `Jenkins` sections with your deployments.

> For more info on how to configure plugins, please refer to the [data source plugins](#data-source-plugins) section

3. Launch `docker-compose`

   ```shell
   make compose
   ```

4. Create a http request to trigger data collect tasks, please replace your [gitlab projectId](plugins/gitlab/README.md#finding-project-id) and [jira boardId](plugins/jira/README.md#find-board-id) in the request body. This can take up to 20 minutes for large projects. (gitlab 10k+ commits or jira 5k+ issues)

   ```shell
   curl -XPOST 'http://localhost:8080/task' \
   -H 'Content-Type: application/json' \
   -d '[
       {
           "plugin": "gitlab",
           "options": {
               "projectId": 8967944
           }
       },
       {
           "plugin": "jira",
           "options": {
               "boardId": 8
           }
       },
       {
           "plugin": "jenkins",
           "options": {}
       }
   ]'
   ```

5. Navigate to grafana dashboard `http://localhost:3002` (username: `admin`, password: `admin`).

### Setup cron job
Commonly, we have requirement to synchorize data periodly. We providered a tool called `lake-cli` to meet that requirement. Check `lake-cli` usage at [here](./cmd/lake-cli/README.md).  

Otherwise, if you just want to use the cron job, please check `docker-compose` version at [here](./devops/sync/README.md)


## Developer Setup<a id="dev-setup"></a>

### Requirements

- <a href="https://docs.docker.com/get-docker" target="_blank">Docker</a>
- <a href="https://golang.org/doc/install" target="_blank">Golang</a>
- Make
  - Mac (Already installed)
  - Windows: [Download](http://gnuwin32.sourceforge.net/packages/make.htm)
  - Ubuntu: `sudo apt-get install build-essential`

### How to setup dev environment
1. Navigate to where you would like to install this project and clone the repository

   ```sh
   git clone https://github.com/merico-dev/lake.git
   cd lake
   ```

2. Install go packages

    ```sh
    make install
    ```

3. Copy sample config files to new local file

    ```sh
    cp .env.example .env
    ```

4. Start the docker containers

    > Make sure the docker application is running before this step

    ```sh
    make compose
    ```

5. Run the project

    ```sh
    make dev
    ```

6. You can now post to /task to create a jira task. This will collect data from Jira

    ```
    curl -XPOST 'localhost:8080/task' \
    -H 'Content-Type: application/json' \
    -d '[{
        "plugin": "jira",
        "options": {
            "boardId": 8
        }
    }]'
    ```

7. Visualize the data in the Grafana Dashboard

    _From here you can see existing data visualized from collected & enriched data_

    - Navigate to http://localhost:3002 (username: `admin`, password: `admin`)
    - You can also create/modify existing/save dashboards to `lake`
    - For more info on working with Grafana in Dev Lake see [Grafana Doc](docs/GRAFANA.md)


## Tests<a id="tests"></a>

Sample tests can be found in `/test/example`

To run the tests: `make test`

## Grafana<a id="grafana"></a>

We use <a href="https://grafana.com/" target="_blank">Grafana</a> as a visualization tool to build charts for the data stored in our database. Using SQL queries we can add panels to build, save, and edit customized dashboards.

All the details on provisioning, and customizing a dashboard can be found in the [Grafana Doc](docs/GRAFANA.md)

## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md)

## Need help?

Message us on <a href="https://discord.com/invite/83rDG6ydVZ" target="_blank">Discord</a>




