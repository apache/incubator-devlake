## Developer Setup

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