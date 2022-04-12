<div align="center">
<br />
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

# DevLake

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
![badge](https://github.com/merico-dev/lake/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/merico-dev/lake)](https://goreportcard.com/report/github.com/merico-dev/lake)
[![Slack](https://img.shields.io/badge/slack-join_chat-success.svg?logo=slack)](https://join.slack.com/t/devlake-io/shared_invite/zt-17b6vuvps-x98pqseoUagM7EAmKC82xQ)

| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |
</div>
<br>
<div align="left">

### What is DevLake?
DevLake brings your DevOps data into one practical, customized, extensible view. Ingest, analyze, and visualize data from an ever-growing list of developer tools, with our open source product.

DevLake is designed for developer teams looking to make better sense of their development process and to bring a more data-driven approach to their own practices. You can ask DevLake many questions regarding your development process. Just connect and query.

### [See demo](https://app-259373083972538368-3002.ars.teamcode.com/d/0Rjxknc7z/demo-homepage?orgId=1)
Username/password:test/test. The demo is based on the data from this repo, merico-dev/lake.

#### Get started with just a few clicks
<table>
  <tr>
    <td valign="middle"><a href="#user-setup">Run DevLake</a></td>
  </tr>
</table>


<br>


<div align="left">
<img src="https://user-images.githubusercontent.com/14050754/145056261-ceaf7044-f5c5-420f-80ca-54e56eb8e2a7.png" width="100%" alt="User Flow" style="border-radius:15px;"/>
<p align="center">User Flow</p>



### What can be accomplished with DevLake?
1. Collect DevOps data across the entire SDLC process and connect data silos
2. A standard <a href="https://github.com/merico-dev/lake/wiki/DataModel.Domain-layer-schema">data model</a> and out-of-the-box <a href="https://github.com/merico-dev/lake/wiki/Metric-Cheatsheet">metrics</a> for software engineering
3. Flexible <a href="https://github.com/merico-dev/lake/blob/main/ARCHITECTURE.md">framework</a> for data collection and ETL, support customized analysis


<br>

## User setup<a id="user-setup"></a>

- If you only plan to run the product locally, this is the **ONLY** section you should need.
- If you want to run in a cloud environment, click <a valign="middle" href="https://www.teamcode.com/tin/clone?applicationId=259777118600769536">
        <img
          src="https://static01.teamcode.com/badge/teamcode-badge-run-in-cloud-en.svg"
          width="120px"
          alt="Teamcode" valign="middle"
        />
      </a> to set up. This is the detailed [guide](https://github.com/merico-dev/lake/wiki/How-to-Set-Up-Dev-Lake-with-Tin).
- Commands written `like this` are to be run in your terminal.

#### Required Packages to Install<a id="user-setup-requirements"></a>

- [Docker](https://docs.docker.com/get-docker)
- [docker-compose](https://docs.docker.com/compose/install/)

NOTE: After installing docker, you may need to run the docker application and restart your terminal

#### Commands to run in your terminal<a id="user-setup-commands"></a>

**IMPORTANT: DevLake doesn't support Database Schema Migration yet,  upgrading an existing instance is likely to break, we recommend that you deploy a new instance instead.**

1. Download `docker-compose.yml` and `env.example` from [latest release page](https://github.com/merico-dev/lake/releases/latest) into a folder.
2. Rename `env.example` to `.env`. For Mac/Linux users, please run `mv env.example .env` in the terminal.
3. Start Docker on your machine, then run `docker-compose up -d` to start the services.
4. Visit `localhost:4000` to set up configuration files.
   >- Navigate to desired plugins on the Integrations page
   >- Please reference the following for more details on how to configure each one:<br>
      > <a href="plugins/jira/README.md" target="_blank">Jira</a><br>
      > <a href="plugins/gitlab/README.md" target="_blank">GitLab</a><br>
      > <a href="plugins/jenkins/README.md" target="_blank">Jenkins</a><br>
      > <a href="plugins/github/README.md" target="_blank">GitHub</a><br>
   >- Submit the form to update the values by clicking on the **Save Connection** button on each form page
   >- `devlake` takes a while to fully boot up. if `config-ui` complaining about api being unreachable, please wait a few seconds and try refreshing the page.


5. Visit `localhost:4000/pipelines/create` to RUN a Pipeline and trigger data collection.


   Pipelines Runs can be initiated by the new "Create Run" Interface. Simply enable the **Data Source Providers** you wish to run collection for, and specify the data you want to collect, for instance, **Project ID** for Gitlab and **Repository Name** for GitHub.

   Once a valid pipeline configuration has been created, press **Create Run** to start/run the pipeline.
   After the pipeline starts, you will be automatically redirected to the **Pipeline Activity** screen to monitor collection activity.

   **Pipelines** is accessible from the main menu of the config-ui for easy access.

   - Manage All Pipelines: `http://localhost:4000/pipelines`
   - Create Pipeline RUN: `http://localhost:4000/pipelines/create`
   - Track Pipeline Activity: `http://localhost:4000/pipelines/activity/[RUN_ID]`

   For advanced use cases and complex pipelines, please use the Raw JSON API to manually initiate a run using **cURL** or graphical API tool such as **Postman**. `POST` the following request to the DevLake API Endpoint.

    ```json
    [
        [
            {
                "plugin": "github",
                "options": {
                    "repo": "lake",
                    "owner": "merico-dev"
                }
            }
        ]
    ]
    ```

   Please refer to this wiki [How to trigger data collection](https://github.com/merico-dev/lake/wiki/How-to-use-the-triggers-page).

6. Click *View Dashboards* button in the top left when done, or visit `localhost:3002` (username: `admin`, password: `admin`).

   We use <a href="https://grafana.com/" target="_blank">Grafana</a> as a visualization tool to build charts for the <a href="https://github.com/merico-dev/lake/wiki/DataModel.Domain-layer-schema">data stored in our database</a>. Using SQL queries, we can add panels to build, save, and edit customized dashboards.

   All the details on provisioning and customizing a dashboard can be found in the [Grafana Doc](docs/GRAFANA.md).

#### Setup cron job

To synchronize data periodically, we provide [`lake-cli`](./cmd/lake-cli/README.md) for easily sending data collection requests along with [a cron job](./devops/sync/README.md) to periodically trigger the cli tool.


<br>

## Developer Setup<a id="dev-setup"></a>

#### Requirements

- <a href="https://docs.docker.com/get-docker" target="_blank">Docker</a>
- <a href="https://golang.org/doc/install" target="_blank">Golang v1.17+</a>
- Make
  - Mac (Already installed)
  - Windows: [Download](http://gnuwin32.sourceforge.net/packages/make.htm)
  - Ubuntu: `sudo apt-get install build-essential`

#### How to setup dev environment
1. Navigate to where you would like to install this project and clone the repository:

   ```sh
   git clone https://github.com/merico-dev/lake.git
   cd lake
   ```

2. Install dependencies for plugins:

   - [RefDiff](plugins/refdiff#development)

3. Install Go packages

    ```sh
	go get
    ```

4. Copy the sample config file to new local file:

    ```sh
    cp .env.example .env
    ```

5. Update the following variables in the file `.env`:

    * `DB_URL`: Replace `mysql:3306` with `127.0.0.1:3306`

5. Start the MySQL and Grafana containers:

    > Make sure the Docker daemon is running before this step.

    ```sh
    docker-compose up -d mysql grafana
    ```

6. Run lake and config UI in dev mode in two seperate terminals:

    ```sh
    # run lake
    make dev
    # run config UI
    make configure-dev
    ```

7. Visit config UI at `localhost:4000` to configure data sources.
   >- Navigate to desired plugins pages on the Integrations page
   >- You will need to enter the required information for the plugins you intend to use.
   >- Please reference the following for more details on how to configure each one:
   >-> <a href="plugins/jira/README.md" target="_blank">Jira</a>
   >-> <a href="plugins/gitlab/README.md" target="_blank">GitLab</a>,
   >-> <a href="plugins/jenkins/README.md" target="_blank">Jenkins</a>
   >-> <a href="plugins/github/README.md" target="_blank">GitHub</a>

   >- Submit the form to update the values by clicking on the **Save Connection** button on each form page

8. Visit `localhost:4000/pipelines/create` to RUN a Pipeline and trigger data collection.


   Pipelines Runs can be initiated by the new "Create Run" Interface. Simply enable the **Data Source Providers** you wish to run collection for, and specify the data you want to collect, for instance, **Project ID** for Gitlab and **Repository Name** for GitHub.

   Once a valid pipeline configuration has been created, press **Create Run** to start/run the pipeline.
   After the pipeline starts, you will be automatically redirected to the **Pipeline Activity** screen to monitor collection activity.

   **Pipelines** is accessible from the main menu of the config-ui for easy access.

   - Manage All Pipelines: `http://localhost:4000/pipelines`
   - Create Pipeline RUN: `http://localhost:4000/pipelines/create`
   - Track Pipeline Activity: `http://localhost:4000/pipelines/activity/[RUN_ID]`

   For advanced use cases and complex pipelines, please use the Raw JSON API to manually initiate a run using **cURL** or graphical API tool such as **Postman**. `POST` the following request to the DevLake API Endpoint.

    ```json
    [
        [
            {
                "plugin": "github",
                "options": {
                    "repo": "lake",
                    "owner": "merico-dev"
                }
            }
        ]
    ]
    ```

   Please refer to this wiki [How to trigger data collection](https://github.com/merico-dev/lake/wiki/How-to-use-the-triggers-page).


9. Click *View Dashboards* button in the top left when done, or visit `localhost:3002` (username: `admin`, password: `admin`).

   We use <a href="https://grafana.com/" target="_blank">Grafana</a> as a visualization tool to build charts for the <a href="https://github.com/merico-dev/lake/wiki/DataModel.Domain-layer-schema">data stored in our database</a>. Using SQL queries, we can add panels to build, save, and edit customized dashboards.

   All the details on provisioning and customizing a dashboard can be found in the [Grafana Doc](docs/GRAFANA.md).


10. (Optional) To run the tests:

    ```sh
    make test
    ```
<br>


## Project Roadmap
- <a href="https://github.com/merico-dev/lake/wiki/Roadmap-2022" target="_blank">Roadmap 2022</a>: Detailed project roadmaps for 2022.
- DevLake already supported following data sources:
    - <a href="plugins/jira/README.md" target="_blank">Jira(Cloud)</a>
    - <a href="plugins/gitextractor/README.md" target="_blank">Git</a>
    - <a href="plugins/github/README.md" target="_blank">GitHub</a>
    - <a href="plugins/gitlab/README.md" target="_blank">GitLab(Cloud)</a>
    - <a href="plugins/jenkins/README.md" target="_blank">Jenkins</a>
- <a href="https://github.com/merico-dev/lake/wiki/Metric-Cheatsheet" target="_blank">Supported engineering metrics</a>: provide rich perspectives to observe and analyze SDLC.

<br>

## Make Contribution
This section lists all the documents to help you contribute to the repo.

- [Architecture](ARCHITECTURE.md): Architecture of DevLake
- [Data Model](https://github.com/merico-dev/lake/wiki/DataModel.Domain-layer-schema): Domain Layer Schema
- [Add a Plugin](/plugins/README.md): Guide to add a plugin
- [Add metrics](/plugins/HOW-TO-ADD-METRICS.md): Guide to add metrics in a plugin
- [Contribution guidelines](CONTRIBUTING.md): Start from here if you want to make contribution

<br>

## Community

- <a href="https://join.slack.com/t/devlake-io/shared_invite/zt-17b6vuvps-x98pqseoUagM7EAmKC82xQ" target="_blank">Slack</a>: Message us on Slack
- <a href="https://github.com/merico-dev/lake/wiki/FAQ" target="_blank">FAQ</a>: Frequently Asked Questions

<br>

## License<a id="license"></a>

This project is licensed under Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
