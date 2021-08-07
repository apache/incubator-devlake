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
4.  Easy-to-setup via docker
5.  Extensible plugin system to add your own data collectors
6.  Designed to process enterprise-scale data

## Contents

Section | Description | Link
:------------ | :------------- | :-------------
Requirements | Underlying software used | [View Section](#requirements)
Installation | Getting all the required files | [View Section](#installation)
Setup | Steps to get up and running | [View Section](#setup)
Core Usage | Using core `lake` features | [View Section](#core-usage)
Plugin Usage | Links to specific plugin usage & details | [View Section](#plugin-usage)
Configuration | Local file config settings info | [Link](CONFIGURATION.md)
Contributing | How to contribute to this repo | [Link](CONTRIBUTING.md)

## Requirements<a id="requirements" />

- [Node.js](https://nodejs.org/en/download)
- [Docker](https://docs.docker.com/get-docker)

## How to run this application<a id="howToRun" />

**NOTE: If you only plan to run the product, this is the only section you should need**

1. Clone this repository and `cd` into it
2. Configure settings for services & plugins with `cp config/docker.sample.js config/docker.js` and edit the newly created file
3. Start the service with `npm run compose-prod`
- you can see the logs with `npm run compose-logs`
- you can stop all docker containers with `npm run compose-down-prod`
4. Send a post request to the service
```
curl -X POST "http://localhost:3001/" -H 'content-type: application/json' \
    -H 'x-token: mytoken' \
    -d '{"jira":{"boardId": 29}, "gitlab": {"projectId": 24547305}}'
```
5. Check the console logs for docker-compose to see when the logs stop collecting your data. This can take up to 30 minutes for large projects. (gitlab 10k+ commits or jira 10k+ issues)
6. Navigate to Grafana Dashboard `https://localhost:3002` (Username: `admin`, password: `admin`)

## Installation<a id="installation" />

1. Clone this repository<br>

   ```shell
   git clone https://github.com/merico-dev/lake.git
   ```
2. Install dependencies with<br>

   ```
   npm i
   ```
3. Configure local settings for services & plugins, see [CONFIGURATION.md](CONFIGURATION.md)

## Setup<a id="setup" />

1. From the root directory, run
   ```shell
   npm run docker
   ```

      > NOTE: If you get an error like this:
      > **"Error response from daemon: invalid mount config for type "bind": bind source path does not exist: /tmp/rabbitmq/etc/"**

      > You can fix it by creating the directories in the terminal:

      > `mkdir -p ./rabbitmq/logs/ ./rabbitmq/etc/ ./rabbitmq/data/`

2. In another tab run
   ```shell
   npm run all
   ```
3. Create a collection job to collect data. See that the:
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

   Or, by using `curl`

   ```shell
   # ee
   curl -X POST "http://localhost:3001/" -H 'content-type: application/json' \
       -H 'x-token: mytoken' \
       -d '{"jira":{"boardId": 8}, "gitlab": {"projectId": 8967944}}'

   # small data set for test
   curl -X POST "http://localhost:3001/" -H 'content-type: application/json' \
       -H 'x-token: mytoken' \
       -d '{"jira":{"boardId": 29}, "gitlab": {"projectId": 24547305}}'
   ```

4. Visualize data in Grafana dashboard

   From here you can see existing data visualized from collected & enriched data

   - Navigate to http://localhost:3002 (username: `admin`, password: `admin`)
   - You can also create/modify existing/save dashboards to `lake`
   - For more info: [Provisioning a Dashboard](#grafana-provisioning-a-dashboard)

**Migrations**

-  Revert all current migrations `npx sequelize-cli db:migrate:undo:all`
-  Run migration with `npx sequelize-cli db:migrate`

## Core Usage<a id="core-usage" />

Section | Section Info
------------ | -------------
Collections | Create a Collection Job
Grafana | Logging In
Grafana | Provisioning a Dashboard

### Collections: Create a Collection Job <a id="create-a-collection-job" />
<details><summary><b>Details</b></summary>
<ol>
    <li>From the terminal, execute <code>npm run all</code></li>
    <li>From Postman (or similar), send a request like:</li>
</ol>

```json

POST http://localhost:3001/

{
    "jira": {
        "boardId": 8
    },
    "gitlab": {
        "projectId": 8967944,
        "branch": "<your-branch-name>", // branch is optional, we fetch Gitlab default branch if this arg is absent
    }
}

```

Or, by using `curl`

```shell
# ee
curl -X POST "http://localhost:3001/" -H 'content-type: application/json' \
    -H 'x-token: mytoken' \
    -d '{"jira":{"boardId": 8}, "gitlab": {"projectId": 8967944}}'

# small data set for test
curl -X POST "http://localhost:3001/" -H 'content-type: application/json' \
    -H 'x-token: mytoken' \
    -d '{"jira":{"boardId": 29}, "gitlab": {"projectId": 24547305}}'
```

3. See that the:
    - collection job was published
    - jira collection ran
    - enrichment job was published
    - jira enrichment ran
</details>

### Grafana: Logging In<a id="grafana-logging-in" />
<details><summary><b>Details</b></summary>
Once the app is up and running, visit <code>http://localhost:3002</code> to view the Grafana dashboard.
<br><br>
Default login credentials are:

- Username: `admin`
- Password: `admin`
</details>

### Grafana: Provisioning a Dasboard<a id="grafana-provisioning-a-dashboard" />
<details><summary><b>Details</b></summary>

To save a dashboard in the `lake` repo and load it:

1. Create a dashboard in browser (visit `/dashboard/new`, or use sidebar)
2. Save dashboard (in top right of screen)
3. Go to dashboard settings (in top right of screen)
4. Click on _JSON Model_ in sidebar
5. Copy code into a new `.json` file in `/grafana/dashboards`
</details>

## Plugin Usage<a id="plugin-usage" />

Section | Section Info | Docs
------------ | ------------- | -------------
Jira | Metrics, Generating API Token, Find Project/Board ID | [Link](src/plugins/jira-pond/README.md)
Gitlab | Metrics, Generating API Token | [Link](src/plugins/gitlab-pond/README.md)

## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md)
