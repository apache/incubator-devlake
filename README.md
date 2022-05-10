<div align="center">
<br />
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

# DevLake

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
![badge](https://github.com/merico-dev/lake/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/merico-dev/lake)](https://goreportcard.com/report/github.com/merico-dev/lake)
[![Slack](https://img.shields.io/badge/slack-join_chat-success.svg?logo=slack)](https://join.slack.com/t/devlake-io/shared_invite/zt-17b6vuvps-x98pqseoUagM7EAmKC82xQ)

| English | [中文](README-zh-CN.md) |
| --- | --- |
</div>
<br>
<div align="left">

### What is DevLake?
DevLake brings your DevOps data into one practical, customized, extensible view. Ingest, analyze, and visualize data from an ever-growing list of developer tools, with our open source product.

DevLake is designed for developer teams looking to make better sense of their development process and to bring a more data-driven approach to their own practices. You can ask DevLake many questions regarding your development process. Just connect and query.

### See [demo based on this repo](https://grafana-lake.demo.devlake.io/d/0Rjxknc7z/demo-homepage?orgId=1)

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

#### Prerequisites

- [Docker v19.03.10+](https://docs.docker.com/get-docker)
- [docker-compose v2.2.3+](https://docs.docker.com/compose/install/)

#### Launch DevLake

1. Download `docker-compose.yml` and `env.example` from [latest release page](https://github.com/merico-dev/lake/releases/latest) into a folder.
2. Rename `env.example` to `.env`. For Mac/Linux users, please run `mv env.example .env` in the terminal.
3. Run `docker-compose up -d` to launch DevLake.

#### Configure data connections and collect data

1. Visit `config-ui` at `http://localhost:4000` in your browser to configure data connections. **For users who'd like to collect GitHub data, we recommend reading our [GitHub data collection guide](./docs/github-user-guide-v0.10.0.md) which covers the following steps in detail.**
   >- Navigate to desired plugins on the Integrations page
   >- Please reference the following for more details on how to configure each one:<br>
      > <a href="plugins/jira/README.md" target="_blank">Jira</a><br>
      > <a href="plugins/gitlab/README.md" target="_blank">GitLab</a><br>
      > <a href="plugins/jenkins/README.md" target="_blank">Jenkins</a><br>
      > <a href="plugins/github/README.md" target="_blank">GitHub</a><br>
   >- Submit the form to update the values by clicking on the **Save Connection** button on each form page
   >- `devlake` takes a while to fully boot up. if `config-ui` complaining about api being unreachable, please wait a few seconds and try refreshing the page.
2. Create pipelines to trigger data collection in `config-ui`
3. Click *View Dashboards* button in the top left when done, or visit `localhost:3002` (username: `admin`, password: `admin`).

   We use <a href="https://grafana.com/" target="_blank">Grafana</a> as a visualization tool to build charts for the <a href="https://github.com/merico-dev/lake/wiki/DataModel.Domain-layer-schema">data stored in our database</a>. Using SQL queries, we can add panels to build, save, and edit customized dashboards.

   All the details on provisioning and customizing a dashboard can be found in the [Grafana Doc](docs/GRAFANA.md).
4. To synchronize data periodically, users can set up recurring pipelines with DevLake's [pipeline blueprint](./docs/recurring-pipeline.md) for details.

#### Upgrade to a newer version

Support for database schema migration was introduced to DevLake in v0.10.0. From v0.10.0 onwards, users can upgrade their instance smoothly to a newer version. However, versions prior to v0.10.0 do not support upgrading to a newer version with a different database schema. We recommend users deploying a new instance if needed.

#### Deploy to Kubernates

We provide a sample [k8s-deploy.yaml](k8s-deploy.yaml) for users interested in deploying DevLake on a k8s cluster.

[k8s-deploy.yaml](k8s-deploy.yaml) will create a namespace `devlake` on your k8s cluster, and use `nodePort 30004` for `config-ui`,  `nodePort 30002` for `grafana` dashboards. If you would like to use certain version of DevLake, please update the image tag of `grafana`, `devlake` and `config-ui` services to specify versions like `v0.10.1`.

Here's the step-by-step guide:

1. Download [k8s-deploy.yaml](k8s-deploy.ymal) to local machine
2. Some key points:
   - `config-ui` deployment:
     * `GRAFANA_ENDPOINT`: FQDN of grafana service which can be reached from user's browser
     * `DEVLAKE_ENDPOINT`: FQDN of devlake service which can be reached within k8s cluster, normally you don't need to change it unless namespace was changed
     * `ADMIN_USER`/`ADMIN_PASS`: Not required, but highly recommended
   - `devlake-config` config map:
     * `MYSQL_USER`: shared between `mysql` and `grafana` service
     * `MYSQL_PASSWORD`: shared between `mysql` and `grafana` service
     * `MYSQL_DATABASE`: shared between `mysql` and `grafana` service
     * `MYSQL_ROOT_PASSWORD`: set root password for `mysql`  service
   - `devlake` deployment:
     * `DB_URL`: update this value if  `MYSQL_USER`, `MYSQL_PASSWORD` or `MYSQL_DATABASE` were changed
3. The `devlake` deployment store its configuration in `/app/.env`. In our sample yaml, we use `hostPath` volume, so please make sure directory `/var/lib/devlake` exists on your k8s workers, or employ other techniques to persist `/app/.env` file. Please do NOT mount the entire `/app` directory, because plugins are located in `/app/bin` folder.
4. Finally, execute the following command, DevLake should be up and running:
    ```sh
    kubectl apply -f k8s-deploy.yaml
    ```


## Developer Setup<a id="dev-setup"></a>

#### Requirements

- <a href="https://docs.docker.com/get-docker" target="_blank">Docker v19.03.10+</a>
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

6. Start the MySQL and Grafana containers:

    > Make sure the Docker daemon is running before this step.

    ```sh
    docker-compose up -d mysql grafana
    ```

7. Run lake and config UI in dev mode in two seperate terminals:

    ```sh
    # run lake
    make dev
    # run config UI
    make configure-dev
    ```

8. Visit config UI at `localhost:4000` to configure data connections.
   >- Navigate to desired plugins pages on the Integrations page
   >- You will need to enter the required information for the plugins you intend to use.
   >- Please reference the following for more details on how to configure each one:
   >-> <a href="plugins/jira/README.md" target="_blank">Jira</a>
   >-> <a href="plugins/gitlab/README.md" target="_blank">GitLab</a>,
   >-> <a href="plugins/jenkins/README.md" target="_blank">Jenkins</a>
   >-> <a href="plugins/github/README.md" target="_blank">GitHub</a>

   >- Submit the form to update the values by clicking on the **Save Connection** button on each form page

9. Visit `localhost:4000/pipelines/create` to RUN a Pipeline and trigger data collection.


   Pipelines Runs can be initiated by the new "Create Run" Interface. Simply enable the **Data Connection Providers** you wish to run collection for, and specify the data you want to collect, for instance, **Project ID** for Gitlab and **Repository Name** for GitHub.

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

   Please refer to [Pipeline Advanced Mode](docs/create-pipeline-advanced-mode.md) for in-depth explanation.


10. Click *View Dashboards* button in the top left when done, or visit `localhost:3002` (username: `admin`, password: `admin`).

   We use <a href="https://grafana.com/" target="_blank">Grafana</a> as a visualization tool to build charts for the <a href="https://github.com/merico-dev/lake/wiki/DataModel.Domain-layer-schema">data stored in our database</a>. Using SQL queries, we can add panels to build, save, and edit customized dashboards.

   All the details on provisioning and customizing a dashboard can be found in the [Grafana Doc](docs/GRAFANA.md).


11. (Optional) To run the tests:

    ```sh
    make test
    ```
    
12. For DB migrations, please refer to [Migration Doc](docs/MIGRATIONS.md).
<br>


## Temporal Mode

Normally, DevLake would execute pipelines on local machine (we call it `local mode`), it is sufficient most of the time.However, when you have too many pipelines that need to be executed in parallel, it can be problematic, either limited by the horsepower or throughput of a single machine.

`temporal mode` was added to support distributed pipeline execution, you can fire up arbitrary workers on multiple machines to carry out those pipelines in parallel without hitting the single machine limitation.

But, be careful, many API services like JIRA/GITHUB have request rate limit mechanism, collect data in parallel against same API service with same identity would most likely hit the wall.

### How it works

1. DevLake Server and Workers connect to the same temporal server by setting up `TEMPORAL_URL`
2. DevLake Server sends `pipeline` to temporal server, and one of the Workers would pick it up and execute


**IMPORTANT: This feature is in early stage of development, use with cautious**


### Temporal Demo

#### Requirements

- [Docker](https://docs.docker.com/get-docker)
- [docker-compose](https://docs.docker.com/compose/install/)
- [temporalio](https://temporal.io/)

#### How to setup

1. Clone and fire up  [temporalio](https://temporal.io/) services
2. Clone this repo, and fire up DevLake with command `docker-compose -f docker-compose-temporal.yml up -d`


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

## How to Contribute
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
- <div>Wechat Group QR Code
![](wechat_group_qr_code.png)</div>

<br>

## License<a id="license"></a>

This project is licensed under Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
