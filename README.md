<div align="center">
<br />
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

# Dev Lake
<p>
    <b>
     <!Software development workflow analysis for free> 
    </b>
  </p>
  <p>

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/discord/844603288082186240.svg?style=flat?label=&logo=discord&logoColor=ffffff&color=747df7&labelColor=454545)](https://discord.gg/83rDG6ydVZ)
![badge](https://github.com/merico-dev/lake/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/merico-dev/lake)](https://goreportcard.com/report/github.com/merico-dev/lake)

| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |


<div align="left">

<br>

### What is Dev Lake?
Dev Lake is the one-stop solution that _**integrates, analyzes, and visualizes**_ software development data throughout the _**software development life cycle (SDLC)**_ for engineering teams.


<img src="https://user-images.githubusercontent.com/2908155/130271622-827c4ffa-d812-4843-b09d-ea1338b7e6e5.png" width="100%" alt="Dev Lake Grafana Dashboard" />
<p align="center">Dashboard Screenshot</p><br>
<img src="https://user-images.githubusercontent.com/14050754/139076905-48d13e40-51ab-49e4-b537-0fe56960a1c0.png" width="100%" alt="Dev Lake Grafana Dashboard" />
<p align="center">User Flow</p>

### Why Dev Lake?
1. Unifies data from multiple sources (<a href="https://www.atlassian.com/software/jira" target="_blank">Jira</a>, <a href="https://gitlab.com/" target="_blank">GitLab</a>, Github, <a href="https://www.jenkins.io/" target="_blank">Jenkins</a>, etc.) in one place to break data isolation.
2. Provides a series of industry-standard metrics to identify engineering problems. For example, long requirement lead time, high bug rate per 1k lines of code, etc.
3. Highly customizable charts, metrics & dashboards for users to find insights based on their own needs.

### What can be accomplished with Dev Lake?
1. Visualize and analyze your entire SDLC process in one personalized, unified view. 
2. Debug process- and team-level issues, scale successes. 
3. Unify and standardize measures of success and benchmarks. 


<br>

## Contents

<table>
    <tr>
        <td>Section</td>
        <td>Subsection</td>
        <td>Description</td>
        <td>Documentation Linl</td>
    </tr>
    <tr>
        <td>Data Sources</td>
        <td>Data Sources</td>
        <td>Links to specific plugin usage & details</td>
        <td><a href="#data-source-plugins">View Section</a></td>
    </tr>
    <tr>
        <td rowspan="3">Setup Guide</td>
        <td>User Setup</td>
        <td>Set up Dev Lake locally as a user</td>
        <td><a href="#user-setup">View Section</a></td>
    </tr>
    <tr>
        <td>Developer Setup</td>
        <td>Set up development environment locally</td>
        <td><a href="#dev-setup">View Section</a></td>
    </tr>
    <tr>
        <td>Cloud Setup</td>
        <td>Set up DevLake in the cloud with Tin</td>
        <td><a href="#cloud-setup">View Section</a></td>
    </tr>
   <tr>
        <td>Tests</td>
        <td>Tests</td>
        <td>Commands for running tests</td>
        <td><a href="#tests">View Section</a></td>
    </tr>
    <tr>
        <td rowspan="3">Make Contribution</td>
        <td>Build a Plugin</td>
        <td>Details on how to make your own plugin</td>
        <td><a href="#plugin">View Section</a></td>
    </tr>
   <tr>
        <td>Add Plugin Metrics</td>
        <td>Guide to add metrics</td>
        <td><a href="#metrics">View Section</a></td>
    </tr>
    <tr>
        <td>Contribution specs</td>
        <td>How to contribute to this repo</td>
        <td><a href="#contributing">View Section</a></td>
    </tr>
    <tr>
        <td rowspan="4">User Guide, Help, and more</td>
        <td>Grafana</td>
        <td>How to visualize the data</td>
        <td><a href="#grafana">View Section</a></td>
    </tr>
    <tr>
        <td>Need Help</td>
        <td>Message us on Discord</td>
        <td><a href="#help">View Section</a></td>
    </tr>
    <tr>
        <td>FAQ</td>
        <td>Frequently asked questions by users</td>
        <td><a href="#faq">View Section</a></td>
    </tr>
    <tr>
        <td>License</td>
        <td>The project license</td>
        <td><a href="#license">View Section</a></td>
    </tr>
</table>


## Data Sources We Currently Support<a id="data-source-plugins"></a>

Below is a list of _data source plugins_ used to collect & enrich data from specific sources. Each has a `README.md` file with basic setup, troubleshooting, and metrics info.

For more information on building a new _data source plugin_, see [Build a Plugin](plugins/README.md).

Section | Section Info | Docs
------------ | ------------- | -------------
Jira | Summary, Data & Metrics, Configuration, Plugin API | <a href="plugins/jira/README.md" target="_blank">Link</a>
GitLab | Summary, Data & Metrics, Configuration, Plugin API | <a href="plugins/gitlab/README.md" target="_blank">Link</a> 
Jenkins | Summary, Data & Metrics, Configuration, Plugin API | <a href="plugins/jenkins/README.md" target="_blank">Link</a>
GitHub | Summary, Data & Metrics, Configuration, Plugin API | <a href="plugins/github/README.md" target="_blank">Link</a>


## Setup Guide
There're 3 ways to set up Dev Lake.

### User setup<a id="user-setup"></a>

**NOTES:**

- **If you only plan to run the product locally, this is the only section you should need.**
- **Commands written `like this` are to be run in your terminal.**

#### Required Packages to Install<a id="user-setup-requirements"></a>

- [Docker](https://docs.docker.com/get-docker)
- [docker-compose](https://docs.docker.com/compose/install/)

**NOTE:** After installing docker, you may need to run the docker application and restart your terminal

#### Commands to run in your terminal<a id="user-setup-commands"></a>

1. Clone repository:

   ```sh
   git clone https://github.com/merico-dev/lake.git devlake
   cd devlake
   cp .env.example .env
   ```
2. Start Docker on your machine, then run `docker-compose up -d` to start the services.
3. Visit `localhost:4000` to setup devlake.

   > For more info on how to configure plugins, please refer to the <a href="https://github.com/merico-dev/lake#data-source-plugins" target="_blank">data source plugins</a> section

3. Visit `localhost:4000` to setup configuration files.
   >- Navigate to desired plugins pages on the Integrations page
   >- You will need to enter the required information for the plugins you intend to use.
   >- Please reference the following for more details on how to configure each one:
   >-> <a href="plugins/jira/README.md" target="_blank">Jira</a>
   >-> <a href="plugins/gitlab/README.md" target="_blank">GitLab</a>
   >-> <a href="plugins/jenkins/README.md" target="_blank">Jenkins</a> 
   >-> <a href="plugins/github/README.md" target="_blank">GitHub</a>

   >- Submit the form to update the values by clicking on the **Save Connection** button on each form page

   >- `devlake` takes a while to fully boot up. if `config-ui` complaining about api being unreachable, please wait a few seconds and try refreshing the page. 
   >- To collect this repo for a quick preview, you must put your github into `Auth Token(s)` **Data Integrations / Github** page.

4. Visit `localhost:4000/triggers` to trigger data collection.


   > - Please replace your [GitLab projectId](plugins/gitlab/README.md#finding-project-id) and [Jira boardId](plugins/jira/README.md#find-board-id) in the request body. Click the **Trigger Collection** button. Data collection can take up to 20 minutes for large projects. (GitLab 10k+ commits or Jira 5k+ issues)
   > - To collect this repo for a quick preview, you can  use the following JSON
   >   ```json
   >   [
   >     [
   >       {
   >         "Plugin": "github",
   >         "Options": {
   >           "repositoryName": "lake",
   >           "owner": "merico-dev"
   >         }
   >       }
   >     ],
   >     [
   >       {
   >         "plugin": "github-domain",
   >         "options": {}
   >       }
   >     ]
   >   ]
   >   ```


5. Click *View Dashboards* button when done (username: `admin`, password: `admin`). The button will be shown on the Trigger Collection page when data collection has finished.

#### Setup cron job
Commonly, we have the requirement to synchronize data periodically. We provided a tool called `lake-cli` to meet that requirement. Check `lake-cli` usage [here](./cmd/lake-cli/README.md).  

Otherwise, if you just want to use the cron job, please check `docker-compose` version at [here](./devops/sync/README.md)

### Developer Setup<a id="dev-setup"></a>

#### Requirements

- <a href="https://docs.docker.com/get-docker" target="_blank">Docker</a>
- <a href="https://golang.org/doc/install" target="_blank">Golang</a>
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

2. Install Go packages:

    ```sh
    make install
    ```

3. Copy the sample config file to new local file:

    ```sh
    cp .env.example .env
    ```
   Find the line start with `DB_URL` in the file `.env` and replace `mysql:3306` with `127.0.0.1:3306`

4. Start the Docker containers:

    > Make sure the Docker daemon is running before this step.

    ```sh
    make compose
    ```

5. Run the project:

    ```sh
    make dev
    ```

6. Visit `localhost:4000` to setup configuration files.
   >- Navigate to desired plugins pages on the Integrations page
   >- You will need to enter the required information for the plugins you intend to use.
   >- Please reference the following for more details on how to configure each one:
   >-> <a href="plugins/jira/README.md" target="_blank">Jira</a>
   >-> <a href="plugins/gitlab/README.md" target="_blank">GitLab</a>, 
   >-> <a href="plugins/jenkins/README.md" target="_blank">Jenkins</a> 
   >-> <a href="plugins/github/README.md" target="_blank">GitHub</a> 

   >- Submit the form to update the values by clicking on the **Save Connection** button on each form page

7. Visit `localhost:4000/triggers` to trigger data collection.

   > Please replace your [GitLab projectId](plugins/gitlab/README.md#finding-project-id) and [Jira boardId](plugins/jira/README.md#find-board-id) in the request body. Click the **Trigger Collection** button. Data collection can take up to 20 minutes for large projects. (GitLab 10k+ commits or Jira 5k+ issues)

8. Click *View Dashboards* button when done (username: `admin`, password: `admin`). The button is shown in the top left.


### Cloud setup<a id="cloud-setup"></a>
If you want to run Dev Lake in a clound environment, you can set up Dev Lake with Tin. [See detailed setup guide](https://github.com/merico-dev/lake/wiki/How-to-Set-Up-Dev-Lake-with-Tin)

**Disclaimer:** 
> To protect your information, it is critical for users of the Tin hosting to set passwords to protect Dev Lake applications. We built Dev Lake as a self-hosted product, in part to ensure users have total protection and ownership of their data, while the same remains true for the Tin hosting, this risk point can only be eliminated by the end-user.

<br>

## Tests<a id="tests"></a>

To run the tests:

```sh
make test
```
<br>

## Make Contribution
This section list all the documents to help you contribute to the repo.

### Add a Plugin<a id="plugin"></a>

[plugins/README.md](/plugins/README.md)

### Add Plugin Metrics<a id="metrics"></a>

[plugins/HOW-TO-ADD-METRICS.md](/plugins/HOW-TO-ADD-METRICS.md)

### Contributing Spec<a id="contributing"></a>

[CONTRIBUTING.md](CONTRIBUTING.md)

<br>

## User Guide, Help and more
### Grafana<a id="grafana"></a>

We use <a href="https://grafana.com/" target="_blank">Grafana</a> as a visualization tool to build charts for the data stored in our database. Using SQL queries, we can add panels to build, save, and edit customized dashboards.

All the details on provisioning and customizing a dashboard can be found in the [Grafana Doc](docs/GRAFANA.md).



### Need help?<a id="help"></a>

Message us on <a href="https://discord.com/invite/83rDG6ydVZ" target="_blank">Discord</a>


### FAQ<a id="faq"></a>

Q: When I run ``` docker-compose up -d ``` I get this error: "qemu: uncaught target signal 11 (Segmentation fault) - core dumped". How do I fix this?

A: M1 Mac users need to download a specific version of docker on their machine. You can find it <a href="https://docs.docker.com/desktop/mac/apple-silicon/" target="_blank">here</a>.


### License<a id="license"></a>

This project is licensed under Apache License 2.0 - see the [`LICENSE`](LICENSE) file for details.