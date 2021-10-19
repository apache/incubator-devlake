# Jira

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

In order to fully use this plugin, you will need to set various config variables in the root of this project in the `.env` file.

### Set JIRA_ENDPOINT

This is the setting that is used as the base of all Jira API calls. You can see this in all of your Jira Urls as the start of the Url.

**Example:** 

If you see `https://mydomain.atlassian.net/secure/RapidBoard.jspa?rapidView=999&projectKey=XXX`, you will need to set `JIRA_ENDPOINT=https://mydomain.atlassian.net` in your `.env` file.

## Generating API token
1. Once logged into Jira, visit the url `https://id.atlassian.com/manage-profile/security/api-tokens`
2. Click the **Create API Token** button, and give it any label name
![image](https://user-images.githubusercontent.com/27032263/129363611-af5077c9-7a27-474a-a685-4ad52366608b.png)
3. Encode with login email using command `echo -n <jira login email>:<jira token> | base64`

NOTE: You can see your project's issue statuses here:

<img width="2035" alt="Screen Shot 2021-09-10 at 4 01 56 PM" src="https://user-images.githubusercontent.com/2908155/133310611-2c5e1254-3456-4e15-9c3c-458fed03c6d3.png">

Or you can make a cUrl request to see the statuses:

```
curl --location --request GET 'https://<YOUR_JIRA_ENDPOINT>/rest/api/2/project/<PROJECT_ID>/statuses' \
--header 'Authorization: Basic <BASE64_ENCODED_TOKEN>' \
--header 'Content-Type: application/json'
```

### Set JIRA_BASIC_AUTH_ENCODED

1. Once logged into Jira, visit the url: <a href="https://id.atlassian.com/manage-profile/security/api-tokens" target="_blank">https://id.atlassian.com/manage-profile/security/api-tokens</a>
2. Click the **Create API Token** button, and give it any label name
![image](https://user-images.githubusercontent.com/27032263/129363611-af5077c9-7a27-474a-a685-4ad52366608b.png)
3. Encode with login email using the command `echo -n <jira login email>:<jira token> | base64`

### Set Jira Custom Fields

This applies to JIRA_ISSUE_STORYPOINT_FIELD and JIRA_ISSUE_EPIC_KEY_FIELD

Custom fields can be applied to Jira stories. We use this to set `JIRA_ISSUE_EPIC_KEY_FIELD` in the `.env` file.

**Example:** `JIRA_ISSUE_EPIC_KEY_FIELD=customfield_10024`

Please follow this guide: [How to find Jira the custom field ID in Jira? · merico-dev/lake Wiki](https://github.com/merico-dev/lake/wiki/How-to-find-the-custom-field-ID-in-Jira)

### Set Issue Type Mapping<a id="issue-type-mapping"></a>

Different companies might use different issue types to represent their Bug/Incident/Requirement,
type mappings allow Devlake to recognize your specific setup with respect to Jira statuses.
Devlake supports three different standard issue types:

 - `Bug`
 - `Incident`
 - `Requirement`

For example, say we were using `Story` to represent our Requirement, what we have to do is setting the following
`Environment Variables` before running Devlake:

**Example:** 

```sh
# JIRA_ISSUE_TYPE_MAPPING=<STANDARD_TYPE>:<YOUR_TYPE_1>,<YOUR_TYPE_2>;....
JIRA_ISSUE_TYPE_MAPPING=Requirement:Story;Incident:CustomerComplaint;Bug:QABug;
```

Type mapping is critical for some metrics, like **Requirement Count**, make sure to map your custom type correctly.

### Set JIRA_ISSUE_STORYPOINT_COEFFICIENT

This is a value you can set to something other than the default of 1 if you want to skew the results of story points.

### Find Board Id

1. Navigate to the Jira board in the browser
2. in the URL bar, get the board id from the parameter `?rapidView=`

**Example:**

`https://<your_jira_endpoint>/secure/RapidBoard.jspa?rapidView=51`

![Screen Shot 2021-08-13 at 10 07 19 AM](https://user-images.githubusercontent.com/27032263/129363083-df0afa18-e147-4612-baf9-d284a8bb7a59.png)

Your board id is used in all REST requests to DevLake. You do not need to set this in your `.env` file. 
## How do I find the custom field ID in Jira?
Using URL
1. Navigate to Administration >> Issues >> Custom Fields .
2. Click the cog and hover over Configure or Screens option.
3. Observe the URL at the bottom left of the browser window. Example: The id for this custom field is 10006.

## How to Trigger Data Collection for This Plugin

**Example:** 

```
curl -XPOST 'localhost:8080/task' \
-H 'Content-Type: application/json' \
-d '[[{
    "plugin": "jira",
    "options": {
        "sourceId": 1,
        "boardId": 8
    }
}]]'
```

## API

### Data Sources Management

#### Data Sources

- Get all data source
```
GET /plugins/jira/sources


[
  {
    "ID": 14,
    "CreatedAt": "2021-10-11T11:49:19.029Z",
    "UpdatedAt": "2021-10-11T11:49:19.029Z",
    "name": "test-jira-source",
    "endpoint": "https://merico.atlassian.net/rest",
    "basicAuthEncoded": "basicAuth",
    "epicKeyField": "epicKeyField",
    "storyPointField": "storyPointField",
    "StoryPointCoefficient": 0.5
  }
]
```
- Create a new data source
```
POST /plugins/jira/sources
{
	"name": "jira data source name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <jira login email>:<jira token> | base64`",
	"epicKeyField": "name of customfield of epic key",
	"storyPointField": "name of customfield of story point",
	"storyPointCoefficient": 1,   // help converting user storypoint to stand storypoint
	"typeMappings": { // optional, send empty object to delete all typeMappings of the data source
		"userType": {
			"standardType": "devlake standard type"
		}
	}
}
```
- Update data source
```
PUT /plugins/jira/sources/:sourceId
{
	"name": "jira data source name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <jira login email>:<jira token> | base64`",
	"epicKeyField": "name of customfield of epic key",
	"storyPointField": "name of customfield of story point",
	"storyPointCoefficient": 1,   // help converting user storypoint to stand storypoint
	"typeMappings": { // optional, send empty object to delete all typeMappings of the data source
		"userType": {
			"standardType": "devlake standard type",
		}
	}
}
```
- Get data source detail
```
GET /plugins/jira/sources/:sourceId


{
	"name": "jira data source name",
	"endpoint": "jira api endpoint, i.e. https://merico.atlassian.net/rest",
	"basicAuthEncoded": "generated by `echo -n <jira login email>:<jira token> | base64`",
	"epicKeyField": "name of customfield of epic key",
	"storyPointField": "name of customfield of story point",
	"storyPointCoefficient": 1,   // help converting user storypoint to stand storypoint
	"typeMappings": { // optional, send empty object to delete all typeMappings of the data source
		"userType": {
			"standardType": "devlake standard type",
		}
	}
}
```
- Delete data source
```
DELETE /plugins/jira/sources/:sourceId
```

#### Type mappings

- Get all type mappings
```
GET /plugins/jira/sources/:sourceId/type-mappings


[
  {
    "jiraSourceId": 16,
    "userType": "userType",
    "standardType": "standardType"
  }
]
```
- Create a new type mapping
```
POST /plugins/jira/sources/:sourceId/type-mappings
{
    "userType": "userType",
    "standardType": "standardType"
}
```
- Update type mapping
```
PUT /plugins/jira/sources/:sourceId/type-mapping/:userType
{
    "standardType": "standardTypeUpdated"
}
```
- Delete type mapping
```
DELETE /plugins/jira/sources/:sourceId/type-mapping/:userType
```
