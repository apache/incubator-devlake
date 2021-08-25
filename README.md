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
Grafana | How to visualize the data | [View Section](#grafana)
Requirements | Underlying software used | [View Section](#requirements)
Setup | Steps to setup the project | [View Section](#setup)
Migrations | Commands for running migrations | [View Section](#migrations)
Tests | Commands for running tests | [View Section](#tests)
⚠️ (WIP) Build a Plugin | Details on how to make your own | [Link](src/plugins/README.md)
⚠️ (WIP) Add Plugin Metrics | Guide to adding plugin metrics | [Link](src/plugins/HOW-TO-ADD-METRICS.md)
Contributing | How to contribute to this repo | [Link](CONTRIBUTING.md)



## ⚠️ (WIP) Data Sources We Currently Support<a id="data-source-plugins"></a>

Below is a list of _data source plugins_ used to collect & enrich data from specific sources. Each have a `README.md` file with basic setup, troubleshooting and metrics info.

For more information on building a new _data source plugin_ see [Build a Plugin](src/plugins/README.md).

Section | Section Info | Docs
------------ | ------------- | -------------
Jira | Metrics, Generating API Token, Find Project/Board ID | [Link](src/plugins/jira-pond/README.md)
Gitlab | Metrics, Generating API Token | [Link](src/plugins/gitlab-pond/README.md)

## Grafana<a id="grafana"></a>

We use <a href="https://grafana.com/" target="_blank">Grafana</a> as a visualization tool to build charts for the data stored in our database. Using SQL queries we can add panels to build, save, and edit customized dashboards.

All the details on provisioning, and customizing a dashboard can be found in the [Grafana Doc](docs/GRAFANA.md)

---

## Requirements<a id="requirements"></a>

- <a href="https://docs.docker.com/get-docker" target="_blank">Docker</a>
- <a href="https://golang.org/doc/install" target="_blank">Golang</a>
- Make
  - Mac (Already installed)
  - Windows: [Download](http://gnuwin32.sourceforge.net/packages/make.htm)
  - Ubuntu: `sudo apt-get install build-essential`

## Setup<a id="setup"></a>

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

4. Build the project

    ```sh
    make build
    ```

5. Start the docker container

    > Make sure the docker application is running before this step

    ```sh
    make compose
    ```

6. While docker is running, in a new terminal (from project folder) run:

    ```sh
    cd lake
    ./lake
    ```
7. Collect & enrich data from selected sources and plugins

    _These plugins can be selected from the above list ([View List](#data-source-plugins)), and options for them will be outlined in their specific document._

    **Example:**

    ```sh
    curl --location --request POST 'localhost:8080/source' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "Plugin": "Jira",
        "Options": {}

    }'
    ```

8. Visualize the data in the Grafana Dashboard

    _From here you can see existing data visualized from collected & enriched data_

    - Navigate to http://localhost:3002 (username: `admin`, password: `admin`)
    - You can also create/modify existing/save dashboards to `lake`
    - For more info on working with Grafana in Dev Lake see [Grafana Doc](docs/GRAFANA.md)


## ⚠️ (WIP) Migrations<a id="migrations"></a>

- Make a migration: `<>`
- Migrate your DB: `<>`
- Undo a migration: `<>`

## Tests<a id="tests"></a>

Sample tests can be found in `/test/example`

To run the tests: `make test`

  ## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md)

## Need help?

Message us on <a href="https://discord.com/invite/83rDG6ydVZ" target="_blank">Discord</a>
