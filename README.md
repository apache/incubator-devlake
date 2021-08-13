<br />
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

# Dev Lake
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/discord/844603288082186240.svg?style=flat?label=&logo=discord&logoColor=ffffff&color=747df7&labelColor=454545)](https://discord.gg/83rDG6ydVZ)
![badge](https://github.com/merico-dev/lake/actions/workflows/main.yml/badge.svg)

### What is Dev Lake?

Dev Lake is the one-stop solution that _**integrates, analyzes, and visualizes**_ the development data throughout the _**software development life cycle (SDLC)**_ for engineering teams.

### Why choose Dev Lake?

1.  Supports various data sources and quickly growing
2.  Comprehensive dev metrics built-in
3.  Customizable visualizations and dashboard
4.  Easy-to-setup via [docker](https://docs.docker.com/desktop/)
5.  Extensible plugin system to add your own data collectors
6.  Designed to process enterprise-scale data
7.  Easily build and view new charts with [Grafana](https://grafana.com/)

## Contents

Section | Description | Documentation Link
:------------ | :------------- | :-------------
Requirements | Underlying software used | [View Section](#requirements)
User Setup | Quick and easy setup | [View Section](#user-setup)
Data Source Plugins | Links to specific plugin usage & details | [View Section](#data-source-plugins)
Developer Setup | Steps to get up and running | [View Section](#developer-setup)
Build a Plugin | Details on how to make your own | [Link](src/plugins/README.md)
Add Plugin Metrics | Guide to adding plugin metrics | [Link](src/plugins/HOW-TO-ADD-METRICS.md)
Grafana | How to visualize the data | [Link](docs/GRAFANA.md)
Contributing | How to contribute to this repo | [Link](CONTRIBUTING.md)

## Required Packages to Install<a id="requirements"></a>

- [Docker](https://docs.docker.com/get-docker)
- [Node.js](https://nodejs.org/en/download) (Developer setup only)

**NOTE:** After installing docker, you may need to run the docker application and restart your terminal

## Data Sources We Currently Support<a id="data-source-plugins"></a>

Below is a list of _data source plugins_ used to collect & enrich data from specific sources. Each have a `README.md` file with basic setup, troubleshooting and metrics info.

For more information on building a new _data source plugin_ see [Build a Plugin](src/plugins/README.md).

Section | Section Info | Docs
------------ | ------------- | -------------
Jira | Metrics, Generating API Token, Find Project/Board ID | [Link](src/plugins/jira-pond/README.md)
Gitlab | Metrics, Generating API Token | [Link](src/plugins/gitlab-pond/README.md)

## Grafana

<img src="https://user-images.githubusercontent.com/3789273/128533901-3107e9bf-c3e3-4320-ba47-879fe2b0ea4d.png" width="450px" />

We use [Grafana](https://grafana.com/) as a visualization tool to build charts for the data stored in our database. Using SQL queries we can add panels to build, save, and edit customized dashboards.

All the details on provisioning, and customizing a dashboard can be found in the [Grafana Doc](docs/GRAFANA.md)

## User Setup<a id="user-setup"></a>

**NOTE: If you only plan to run the product, this is the only section you should need**

1. Clone this repository<br>

   ```shell
   git clone https://github.com/merico-dev/lake.git
   cd lake
   ```
2. Create a copy of the sample configuration files with

   ```
   cp config/docker.sample.js config/docker.js
   cp config/plugins.sample.js config/plugins.js
   ```

3. Configure settings for services & plugins by editing the newly created config files. The comments will guide you through the process and look for "Replace" keyword in these config files would help as well. For how to configure plugins, please refer to the [data source plugins](#data-source-plugins) section.

4. Start the service with `npm start`
    > you can stop all docker containers with `npm run stop`

5. Run `docker-compose logs -f lake` to check the logs and see when lake stops collecting your data. This can take up to 20 minutes for large projects. (gitlab 10k+ commits or jira 5k+ issues)

6. Navigate to Grafana Dashboard `https://localhost:3002` (Username: `admin`, password: `admin`)

## Developer Setup<a id="developer-setup"></a>

[docs/DEVELOPER_SETUP.md](docs/DEVELOPER_SETUP.md)

## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md)

## Need help?

Message us on [Discord](https://discord.com/invite/83rDG6ydVZ)!