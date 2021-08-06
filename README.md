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
User Setup | Quick and easy setup | [View Section](#user-setup)
Developer Setup | Steps to get up and running | [View Section](#developer-setup)
Plugins | Links to specific plugin usage & details | [View Section](#plugins)
Configuration | Local file config settings info | [Link](CONFIGURATION.md)
Contributing | How to contribute to this repo | [Link](CONTRIBUTING.md)

## Requirements<a id="requirements" />

- [Node.js](https://nodejs.org/en/download)
- [Docker](https://docs.docker.com/get-docker)

## User Setup<a id="user-setup" />

**NOTE: If you only plan to run the product, this is the only section you should need**

1. Clone this repository and `cd` into it
2. Configure settings for services & plugins with `cp config/docker.sample.js config/docker.js` and edit the newly created file
3. Start the service with `npm run compose-prod`
    > you can see the logs with `npm run compose-logs`

    > you can stop all docker containers with `npm run compose-down-prod`
4. Send a post request to the service
```
curl -X POST "http://localhost:3001/" -H 'content-type: application/json' \
    -d '{"jira":{"boardId": 29}, "gitlab": {"projectId": 24547305}}'
```
5. Check the console logs for docker-compose to see when the logs stop collecting your data. This can take up to 30 minutes for large projects. (gitlab 10k+ commits or jira 10k+ issues)
6. Navigate to Grafana Dashboard `https://localhost:3002` (Username: `admin`, password: `admin`)

## Developer Setup<a id="developer-setup" />

1. Clone this repository<br>

   ```shell
   git clone https://github.com/merico-dev/lake.git
   ```
2. Install dependencies with<br>

   ```
   npm i
   ```
3. Configure local settings for services & plugins, see [CONFIGURATION.md](CONFIGURATION.md)

4. From the root directory, run
   ```shell
   npm run docker
   ```

      > NOTE: If you get an error like this:
      > **"Error response from daemon: invalid mount config for type "bind": bind source path does not exist: /tmp/rabbitmq/etc/"**

      > You can fix it by creating the directories in the terminal:

      > `mkdir -p ./rabbitmq/logs/ ./rabbitmq/etc/ ./rabbitmq/data/`

5. In another tab run
   ```shell
   npm run all
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

   Or, by using `curl`

   ```shell
   # ee
   curl -X POST "http://localhost:3001/" -H 'content-type: application/json' \
       -d '{"jira":{"boardId": 8}, "gitlab": {"projectId": 8967944}}'

   # small data set for test
   curl -X POST "http://localhost:3001/" -H 'content-type: application/json' \
       -d '{"jira":{"boardId": 29}, "gitlab": {"projectId": 24547305}}'
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

## Plugins<a id="plugins" />

Section | Section Info | Docs
------------ | ------------- | -------------
Jira | Metrics, Generating API Token, Find Project/Board ID | [Link](src/plugins/jira-pond/README.md)
Gitlab | Metrics, Generating API Token | [Link](src/plugins/gitlab-pond/README.md)

## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md)
