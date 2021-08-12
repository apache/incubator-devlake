
<br />
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

# Dev Lake
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/discord/844603288082186240.svg?style=flat?label=&logo=discord&logoColor=ffffff&color=747df7&labelColor=454545)](https://discord.gg/83rDG6ydVZ)
![badge](https://github.com/merico-dev/lake/actions/workflows/main.yml/badge.svg)

### What is Dev Lake?

Dev Lake is the one-stop solution that _**integrates, analyzes, and visualizes**_ the development data throughout the _**software development life cycle (SDLC)**_ for engineering teams.

<img width="1769" alt="Screen Shot 2021-08-12 at 4 52 09 PM" src="https://user-images.githubusercontent.com/3011407/129260553-968f0993-c88a-424f-9041-52127b309403.png">

### Why choose Dev Lake?

1.  Supports various data sources and quickly growing
2.  Comprehensive dev metrics built-in
3.  Customizable visualizations and dashboard
4.  Easy-to-setup via docker
5.  Extensible plugin system to add your own data collectors
6.  Designed to process enterprise-scale data

## Contents

Section | Description | Link
:------------ | :------------- | :-------------
Requirements | Underlying software used | [View Section](#requirements)
User Setup | Quick and easy setup | [View Section](#user-setup)
Data Source Plugins | Links to specific plugin usage & details | [View Section](#data-source-plugins)
Developer Setup | Steps to get up and running | [View Section](#developer-setup)
Build a Plugin | Details on how to make your own | [Link](src/plugins/README.md)
Add Plugin Metrics | Guide to adding plugin metrics | [Link](src/plugins/HOW-TO-ADD-METRICS.md)
Grafana | How to visualize the data | [Link](docs/GRAFANA.md)
Contributing | How to contribute to this repo | [Link](CONTRIBUTING.md)

## Requirements<a id="requirements"></a>

- [Docker](https://docs.docker.com/get-docker)
- [Node.js](https://nodejs.org/en/download) (Developer setup only)

## Data Source Plugins<a id="data-source-plugins"></a>

Below is a list of _data source plugins_ used to collect & enrich data from specific sources. Each have a `README.md` file with basic setup, troubleshooting and metrics info.

For more information on building a new _data source plugin_ see [Build a Plugin](src/plugins/README.md).

Section | Section Info | Docs
------------ | ------------- | -------------
Jira | Metrics, Generating API Token, Find Project/Board ID | [Link](src/plugins/jira-pond/README.md)
Gitlab | Metrics, Generating API Token | [Link](src/plugins/gitlab-pond/README.md)

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

1. Clone this repository<br>

   ```shell
   git clone https://github.com/merico-dev/lake.git
   cd lake
   ```
2. Install dependencies with<br>

   ```
   npm i
   ```
3. Create a copy of the sample configuration files with

   ```
   cp config/local.sample.js config/local.js
   cp config/plugins.sample.js config/plugins.js
   ```
4. Configure settings for services & plugins by editing the newly created config files. The comments will guide you through the process and look for "Replace" keyword in these config files would help as well. For how to configure plugins, please refer to the [data source plugins](#data-source-plugins) section.

5. Start all third-party services and lake's own services with

   ```
   npm run dev
   ```
6. Create a collection job to collect data. See that the:
      - collection job was published
      - _lake plugin_ collection ran
      - enrichment job was published
      - _lake plugin_ enrichment ran<br><br>

      > This process will run through each lake plugin, collecting data from each<br>

   From Postman (or similar), send a request like (`branch` is optional):

   ```json
   POST http://localhost:3001/

    {
        "jira": {
            "boardId": 8
        },
        "gitlab": {
            "projectId": 8967944,
            "branch": "<your-branch-name>",
        }
    }
   ```

   Or, using curl:

   ```
   curl -X POST "http://localhost:3001/" -H 'content-type: application/json' \
    -d '{"jira":{"boardId": 8}, "gitlab": {"projectId": 8967944}}'
   ```

7. Visualize data in Grafana dashboard

   From here you can see existing data visualized from collected & enriched data

   - Navigate to http://localhost:3002 (username: `admin`, password: `admin`)
   - You can also create/modify existing/save dashboards to `lake`
   - For more info on working with Grafana in Dev Lake see [Grafana Doc](docs/GRAFANA.md)

**Migrations**

-  Revert all current migrations `npx sequelize-cli db:migrate:undo:all`
-  Run migration with `npx sequelize-cli db:migrate`

## Grafana

<img src="https://user-images.githubusercontent.com/3789273/128533901-3107e9bf-c3e3-4320-ba47-879fe2b0ea4d.png" width="450px" />

We use Grafana as a visualization tool to build charts for the data stored in our database. Using SQL queries we can add panels to build, save, and edit customized dashboards.

All the details on provisioning, and customizing a dashboard can be found in the [Grafana Doc](docs/GRAFANA.md)

## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md)
