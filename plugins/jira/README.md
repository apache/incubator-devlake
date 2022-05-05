# Jira

<div align="center">

| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |

</div>

<br>

## Summary

This plugin collects Jira data through Jira Cloud REST API. It then computes and visualizes various engineering metrics from the Jira data.

<img width="2035" alt="Screen Shot 2021-09-10 at 4 01 55 PM" src="https://user-images.githubusercontent.com/2908155/132926143-7a31d37f-22e1-487d-92a3-cf62e402e5a8.png">

## Project Metrics This Covers

Metric Name | Description
:------------ | :-------------
Requirement Count	| Number of issues with type "Requirement"
Requirement Lead Time	| Lead time of issues with type "Requirement"
Requirement Delivery Rate |	Ratio of delivered requirements to all requirements
Requirement Granularity | Number of story points associated with an issue
Bug Count	| Number of issues with type "Bug"<br><i>bugs are found during testing</i>
Bug Age	| Lead time of issues with type "Bug"<br><i>both new and deleted lines count</i>
Bugs Count per 1k Lines of Code |	Amount of bugs per 1000 lines of code
Incident Count | Number of issues with type "Incident"<br><i>incidents are found when running in production</i>
Incident Age | Lead time of issues with type "Incident"
Incident Count per 1k Lines of Code | Amount of incidents per 1000 lines of code

## Configuration

In order to fully use this plugin, you will need to set various configurations via Dev Lake's `config-ui` service. Open `config-ui` on browser, by default the URL is http://localhost:4000, then go to **Data Integrations / JIRA** page. JIRA plugin currently supports multiple data connections, Here you can **add** new connection to your JIRA connection or **update** the settings if needed.

For each connection, you will need to set up following items:

- Connection Name: This allow you to distinguish different connections.
- Endpoint URL: The JIRA instance api endpoint, for JIRA Cloud Service, it would be: `https://<mydomain>.atlassian.net/rest`. devlake officially supports JIRA Cloud Service on atlassian.net, may or may not work for JIRA Server Instance.
- Basic Auth Token: First, generate a **JIRA API TOKEN** for your JIRA account on JIRA console (see [Generating API token](#generating-api-token)), then, in `config-ui` click the KEY icon on the right side of the input to generate a full `HTTP BASIC AUTH` token for you.
- Issue Type Mapping: JIRA is highly customizable, each JIRA instance may have a different set of issue types than others. In order to compute and visualize metrics for different instances, you need to map your issue types to standard ones. See [Issue Type Mapping](#issue-type-mapping) for detail.
- Epic Key: unfortunately, epic relationship implementation in JIRA is based on `custom field`, which is vary from instance to instance. Please see [Find Out Custom Fields](#find-out-custom-fields).
- Story Point Field: same as Epic Key.
- Remotelink Commit SHA:A regular expression that matches commit links to determine whether an external link is a link to a commit. Taking gitlab as an example, to match all commits similar to https://gitlab.com/merico-dev/ce/example-repository/-/commit/8ab8fb319930dbd8615830276444b8545fd0ad24, you can directly use the regular expression **/commit/([0-9a-f]{40})$**
### Generating API token
1. Once logged into Jira, visit the url `https://id.atlassian.com/manage-profile/security/api-tokens`
2. Click the **Create API Token** button, and give it any label name
![image](https://user-images.githubusercontent.com/27032263/129363611-af5077c9-7a27-474a-a685-4ad52366608b.png)


### Issue Type Mapping

Devlake supports 3 standard types, all metrics are computed based on these types:

 - `Bug`: Problems found during `test` phase, before they can reach the production environment.
 - `Incident`: Problems went through `test` phash, got deployed into production environment.
 - `Requirement`: Normally, it would be `Story` on your instance if you adopted SCRUM.

You can may map arbitrary **YOUR OWN ISSUE TYPE** to a single **STANDARD ISSUE TYPE**, normally, one would map `Story` to `Requirement`, but you could map both `Story` and `Task` to `Requirement` if that was your case. Those unspecified type would be copied as standard type directly for your convenience, so you don't need to map your `Bug` to standard `Bug`.

Type mapping is critical for some metrics, like **Requirement Count**, make sure to map your custom type correctly.

### Find Out Custom Field

Please follow this guide: [How to find Jira the custom field ID in Jira? · merico-dev/lake Wiki](https://github.com/merico-dev/lake/wiki/How-to-find-the-custom-field-ID-in-Jira)

## Collect Data From JIRA

In order to collect data from JIRA, you have to compose a JSON looks like following one, and send it via `Triggers` page on `config-ui`:
<font color=“red”>Warning: Data collection only supports single-task execution, and the results of concurrent multi-task execution may not meet expectations.</font>

```
[
  [
    {
      "plugin": "jira",
      "options": {
          "connectionId": 1,
          "boardId": 8,
          "since": "2006-01-02T15:04:05Z"
      }
    }
  ]
]
```

- `connectionId`: The `ID` field from **JIRA Integration** page.
- `boardId`: JIRA board id, see [Find Board Id](#find-board-id) for detail.
- `since`: optional, download data since specified date/time only.

### Find Board Id

1. Navigate to the Jira board in the browser
2. in the URL bar, get the board id from the parameter `?rapidView=`

**Example:**

`https://<your_jira_endpoint>/secure/RapidBoard.jspa?rapidView=51`

![Screen Shot 2021-08-13 at 10 07 19 AM](https://user-images.githubusercontent.com/27032263/129363083-df0afa18-e147-4612-baf9-d284a8bb7a59.png)

Your board id is used in all REST requests to DevLake. You do not need to configure this at the data connection level.


## API

### Data Connections Management

#### Data Connections

- Get all data connection
```
GET /plugins/jira/connections


[
  {
    "ID": 14,
    "CreatedAt": "2021-10-11T11:49:19.029Z",
    "UpdatedAt": "2021-10-11T11:49:19.029Z",
    "name": "test-jira-connection",
    "endpoint": "https://merico.atlassian.net/rest",
    "basicAuthEncoded": "basicAuth",
    "epicKeyField": "epicKeyField",
    "storyPointField": "storyPointField",
  }
]
```
- Create a new data connection
```
POST /plugins/jira/connections
{
	"name": "jira data connection name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <jira login email>:<jira token> | base64`",
	"epicKeyField": "name of customfield of epic key",
	"storyPointField": "name of customfield of story point",
	"typeMappings": { // optional, send empty object to delete all typeMappings of the data connection
		"userType": {
			"standardType": "devlake standard type"
		}
	}
}
```
- Update data connection
```
PUT /plugins/jira/connections/:connectionId
{
	"name": "jira data connection name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <jira login email>:<jira token> | base64`",
	"epicKeyField": "name of customfield of epic key",
	"storyPointField": "name of customfield of story point",
	"typeMappings": { // optional, send empty object to delete all typeMappings of the data connection
		"userType": {
			"standardType": "devlake standard type",
		}
	}
}
```
- Get data connection detail
```
GET /plugins/jira/connections/:connectionId


{
	"name": "jira data connection name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <jira login email>:<jira token> | base64`",
	"epicKeyField": "name of customfield of epic key",
	"storyPointField": "name of customfield of story point",
	"typeMappings": { // optional, send empty object to delete all typeMappings of the data connection
		"userType": {
			"standardType": "devlake standard type",
		}
	}
}
```
- Delete data connection
```
DELETE /plugins/jira/connections/:connectionId
```

#### Type mappings

- Get all type mappings
```
GET /plugins/jira/connections/:connectionId/type-mappings


[
  {
    "jiraConnectionId": 16,
    "userType": "userType",
    "standardType": "standardType"
  }
]
```
- Create a new type mapping
```
POST /plugins/jira/connections/:connectionId/type-mappings
{
    "userType": "userType",
    "standardType": "standardType"
}
```
- Update type mapping
```
PUT /plugins/jira/connections/:connectionId/type-mapping/:userType
{
    "standardType": "standardTypeUpdated"
}
```
- Delete type mapping
```
DELETE /plugins/jira/connections/:connectionId/type-mapping/:userType
```
- API forwarding
```
GET /plugins/jira/connections/:connectionId/proxy/rest/*path

For example:
Requests to http://your_devlake_host/plugins/jira/connections/1/proxy/rest/agile/1.0/board/8/sprint
would forward to
https://your_jira_host/rest/agile/1.0/board/8/sprint

{
    "maxResults": 1,
    "startAt": 0,
    "isLast": false,
    "values": [
        {
            "id": 7,
            "self": "https://merico.atlassian.net/rest/agile/1.0/sprint/7",
            "state": "closed",
            "name": "EE Sprint 7",
            "startDate": "2020-06-12T00:38:51.882Z",
            "endDate": "2020-06-26T00:38:00.000Z",
            "completeDate": "2020-06-22T05:59:58.980Z",
            "originBoardId": 8,
            "goal": ""
        }
    ]
}
```
