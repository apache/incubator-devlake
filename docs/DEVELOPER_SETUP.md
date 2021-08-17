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

3. Run the command `npm run config` to setup your configuration files

    > For more info on how to configure plugins, please refer to the [data source plugins](../README.md#data-source-plugins) section

    > To map a custom status for a plugin refer to `/config/plugins.js`<br>
    > Ex: In Jira, if you're using **Rejected** as a **Bug** type, refer to the `statusMappings` sections for issues mapped to **"Bug"**<br>
    > All `statusMappings` contain 2 objects. an open status (_first object_), and a closed status (_second object_)


4. Start all third-party services and lake's own services with

   ```
   npm run dev
   ```

    Your collection jobs should begin collecting data. See that the:

      - collection job was published
      - _lake plugin_ collection ran
      - enrichment job was published
      - _lake plugin_ enrichment ran<br><br>

      > This process will run through each lake plugin, collecting data from each<br>

      > To create a collection job manually, from
      >  Postman (or similar), send a request like below (`branch` is optional):

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

5. Visualize data in Grafana dashboard

    From here you can see existing data visualized from collected & enriched data

    - Navigate to http://localhost:3002 (username: `admin`, password: `admin`)
    - You can also create/modify existing/save dashboards to `lake`
    - For more info on working with Grafana in Dev Lake see [Grafana Doc](docs/GRAFANA.md)

**Migrations**

-  Revert all current migrations `npx sequelize-cli db:migrate:undo:all`
-  Run migration with `npx sequelize-cli db:migrate`
