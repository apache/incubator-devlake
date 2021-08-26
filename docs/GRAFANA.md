# Grafana

<img src="https://user-images.githubusercontent.com/3789273/128533901-3107e9bf-c3e3-4320-ba47-879fe2b0ea4d.png" width="450px" />

When first visiting grafana, you will be provided with a sample dashboard with some basic charts setup from the database

## Contents

Section | Link
:------------ | :-------------
Logging In | [View Section](#logging-in)
Viewing All Dashboards | [View Section](#viewing-all-dashboards)
Customizing a Dashboard | [View Section](#customizing-a-dashboard)
Dashboard Settings | [View Section](#dashboard-settings)
Provisioning a Dashboard | [View Section](#provisioning-a-dashboard)
Troubleshooting DB Connection | [View Section](#troubleshooting-db-connection)

## Logging In<a id="logging-in"></a>

Once the app is up and running, visit `http://localhost:3002` to view the Grafana dashboard.

Default login credentials are:

- Username: `admin`
- Password: `admin`

## Viewing All Dashboards<a id="viewing-all-dashboards"></a>

To see all dashboards created in Grafana visit `/dashboards`

Or, use the sidebar and click on **Manage**:

![Screen Shot 2021-08-06 at 11 27 08 AM](https://user-images.githubusercontent.com/3789273/128534617-1992c080-9385-49d5-b30f-be5c96d5142a.png)


## Customizing a Dashboard<a id="customizing-a-dashboard"></a>

When viewing a dashboard, click the top bar of a panel, and go to **edit**

![Screen Shot 2021-08-06 at 11 35 36 AM](https://user-images.githubusercontent.com/3789273/128535505-a56162e0-72ad-46ac-8a94-70f1c7a910ed.png)

**Edit Dashboard Panel Page:**

![grafana-sections](https://user-images.githubusercontent.com/3789273/128540136-ba36ee2f-a544-4558-8282-84a7cb9df27a.png)

### 1. Preview Area
- **Top Left** is the variable select area (custom dashboard variables, used for switching projects, or grouping data)
- **Top Right** we have a toolbar with some buttons related to the display of the data:
  - View data results in a table
  - Time range selector
  - Refresh data button
- **The Main Area** will display the chart and should update in real time

> Note: Data should refresh automatically, but may require a refresh using the button in some cases

### 2. Query Builder
Here we form the SQL query to pull data into our chart, from our database
- Ensure the **Data Source** is the correct database

  ![Screen Shot 2021-08-06 at 10 14 22 AM](https://user-images.githubusercontent.com/3789273/128545278-be4846e0-852d-4bc8-8994-e99b79831d8c.png)

- Select **Format as Table**, and **Edit SQL** buttons to write/edit queries as SQL

  ![Screen Shot 2021-08-06 at 10 17 52 AM](https://user-images.githubusercontent.com/3789273/128545197-a9ff9cb3-f12d-4331-bf6a-39035043667a.png)

- The **Main Area** is where the queries are written, and in the top right is the **Query Inspector** button (to inspect returned data)

  ![Screen Shot 2021-08-06 at 10 18 23 AM](https://user-images.githubusercontent.com/3789273/128545557-ead5312a-e835-4c59-b9ca-dd5c08f2a38b.png)

### 3. Main Panel Toolbar
In the top right of the window are buttons for:
- Dashboard settings (regarding entire dashboard)
- Save/apply changes (to specific panel)

### 4. Grafana Parameter Sidebar
- Change chart style (bar/line/pie chart etc)
- Edit legends, chart parameters
- Modify chart styling
- Other Grafana specific settings

## Dashboard Settings<a id="dashboard-settings"></a>

When viewing a dashboard click on the settings icon to view dashboard settings. In here there is 2 pages important sections to use:

![Screen Shot 2021-08-06 at 1 51 14 PM](https://user-images.githubusercontent.com/3789273/128555763-4d0370c2-bd4d-4462-ae7e-4b140c4e8c34.png)

- Variables
  - Create variables to use throughout the dashboard panels, that are also built on SQL queries

  ![Screen Shot 2021-08-06 at 2 02 40 PM](https://user-images.githubusercontent.com/3789273/128553157-a8e33042-faba-4db4-97db-02a29036e27c.png)

- JSON Model
  - Copy `json` code here and save it to a new file in `/grafana/dashboards/` with a unique name in the `lake` repo. This will allow us to persist dashboards when we load the app

  ![Screen Shot 2021-08-06 at 2 02 52 PM](https://user-images.githubusercontent.com/3789273/128553176-65a5ae43-742f-4abf-9c60-04722033339e.png)

## Provisioning a Dashboard<a id="provisioning-a-dashboard"></a>

To save a dashboard in the `lake` repo and load it:

1. Create a dashboard in browser (visit `/dashboard/new`, or use sidebar)
2. Save dashboard (in top right of screen)
3. Go to dashboard settings (in top right of screen)
4. Click on _JSON Model_ in sidebar
5. Copy code into a new `.json` file in `/grafana/dashboards`

## Troubleshooting DB Connection<a id="troubleshooting-db-connection"></a>

To ensure we have properly connected our database to the data source in Grafana, check database settings in `./grafana/datasources/datasource.yml`, specifically:
- `database`
- `user`
- `secureJsonData/password`
