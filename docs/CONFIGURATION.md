# Configuration

## Lake Application Configuration

1. Make a copy of `config/local.sample.js` under the name of `config/local.js`
2. We can use default values for most fields except the Jira section. For how to set up basic authorization with Jira, please see this [section](#jira) below

## Core Configuration

### Grafana Connection For Data Visualization (https://localhost:3002)

To ensure we have properly connected our database to the data source in Grafana, check database settings in `./grafana/datasources/datasource.yml`, specifically:
- `database`
- `user`
- `secureJsonData/password`

## Plugin Level Configuration

### Jira Specific String Configuration

This can be set up in `/config/constants.js`.

```
{
  "jira": {
    "mappings": {
      "Closed": ["Done", "Closed", "已关闭"],
      "Bug": "Bug",
      "Incident": "Incident"
    }
  }
}
```

You can set multiple values to map from your system as well. Just put the values in an array.
In this object, you can set the values of the object to map to your Jira status definitions. IE:

```
"jira": {
  "Closed": ["MyClosedStatusInJira"],
  "Bug": "MyBugStatusInJira",
  "Incident": "MyIncidentStatusinJira"
}
```

### Jira Api Keys

__Jira auth setup__

1. Create an API key on Jira
3. Create a __basic auth header__ from your API key - [Jira Docs](https://developer.atlassian.com/cloud/jira/platform/basic-auth-for-rest-apis/#supply-basic-auth-headers)
3. Copy your __basic auth header__ into the `jira.basicAuth` field in `/config/local.js` file
4. Add your jira hostname to the `jira.host` field in the `/config/local.js` file

## Default Database Connection Information

### Postgres Connection

- DB Name: lake
- Hostname: localhost
- Port: 5432
- Username: postgres
- Password: postgres

### MongoDB Connection

- DB Name: test
- Hostname: localhost
- Port: 27017
- Username: (none required)
- Password: (none required)

### RabbitMQ Connection

- Vhost Name: rabbitmq
- Hostname: localhost
- Port: 5672
- Username: guest
- Password: guest

<br>

---

## Other Docs

Section | Description | Link
:------------ | :------------- | :-------------
Requirements | Underlying software used | [Link](../README.md#requirements)
User Setup | Quick and easy setup | [Link](../README.md#user-setup)
Developer Setup | Steps to get up and running | [Link](../README.md#developer-setup)
Plugins | Links to specific plugin usage & details | [Link](../README.md#plugins)
Build a Plugin | Details on how to make your own | [Link](../src/plugins/README.md)
Add Plugin Metrics | Guide to adding plugin metrics | [Link](../src/plugins/HOW-TO-ADD-METRICS.md)
Grafana | How to visualize the data | [Link](GRAFANA.md)
Contributing | How to contribute to this repo | [Link](../CONTRIBUTING.md)
