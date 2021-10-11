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

Same as status mapping, different companies might use different issue types to represent their Bug/Incident/Requirement,
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

### Set Issue status mapping<a id="issue-status-mapping"></a>

Jira is highly customizable, different companies may use different `status` names to represent whether an issue was
resolved or not. One company may name it "Done" and others might name it "Finished".

In order to collect life-cycle information correctly, you'll have to map your specific status to Devlake's standard
status, Devlake supports two standard status:

 - `Resolved`: issue was ended successfully
 - `Rejected`: issue was ended by termination or cancellation

For example, say we were using `Done` and `Cancelled` to represent the final stages of `Story` issues. We will have to set
the following `Environment Variables` in the `.env` file at the root of this project before running Devlake:

```sh
#JIRA_ISSUE_<YOUR_TYPE>_STATUS_MAPPING=<STANDARD_STATUS>:<YOUR_STATUS>;...
JIRA_ISSUE_BUG_STATUS_MAPPING=Resolved:Done;Rejected:Cancelled
JIRA_ISSUE_INCIDENT_STATUS_MAPPING=Resolved:Done;Rejected:Cancelled
JIRA_ISSUE_STORY_STATUS_MAPPING=Resolved:Done;Rejected:Cancelled
```

Status mapping is critical for metrics like **Lead Time** since the `leadtime` that we store in the database is calculated only for **Resolved** issues.

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
    "name": "test-jira-source",
    "endpoint": "https://merico.atlassian.net/rest",
    "basicAuthEncoded": "basicAuth",
    "epicKeyField": "epicKeyField",
    "storyPointField": "storyPointField",
    "storyPointCoefficient": 0.5
}
```
- Update data source
```
PUT /plugins/jira/sources/:sourceId
{
    "name": "test-jira-source-updated",
    "endpoint": "https://merico.atlassian.net/rest",
    "basicAuthEncoded": "basicAuth",
    "epicKeyField": "epicKeyField",
    "storyPointField": "storyPointField",
    "storyPointCoefficient": 0.8
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

#### Status mappings

- Get all status mappings
```
GET /plugins/jira/sources/:sourceId/type-mappings/:userType/status-mappings


[
  {
    "jiraSourceId": 16,
    "userType": "userType",
    "userStatus": "userStatus",
    "standardStatus": "standardStatus"
  }
]
```
- Create a new status mapping
```
POST /plugins/jira/sources/:sourceId/type-mappings/:userType/status-mappings
{
    "userStatus": "userStatus",
    "standardStatus": "standardStatus"
}
```
- Update status mapping
```
PUT /plugins/jira/sources/:sourceId/type-mapping/:userType/status-mappings/:userStatus
{
    "standardStatus": "standardStatusUpdated"
}
```
- Delete status mapping
```
DELETE /plugins/jira/sources/:sourceId/type-mapping/:userType/status-mappings/:userStatus
```
