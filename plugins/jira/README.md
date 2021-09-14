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

Set the following environment variables in `.env` file before launching.
For what's issue status mapping, see [Issue status mapping](#issue-status-mapping) section.
For what's issue type mapping, see [Issue type mapping](#issue-type-mapping) section.

```sh
######################
# Jira configuration #
######################

# Jira: basics #

JIRA_ENDPOINT=https://merico.atlassian.net/rest
# ex: echo -n <jira login email>:<jira token> | base64
JIRA_BASIC_AUTH_ENCODED=emhl..........................................a0QzQUE=

# Jira: issue type #

# Format:
#   STANDARD_TYPE_1:ORIGIN_TYPE_1,ORIGIN_TYPE_2;STANDARD_TYPE_2:....
JIRA_ISSUE_TYPE_MAPPING=Requirement:Story


# Jira: issue status #

# Format:
#   JIRA_ISSUE_<STANDARD_ISSUE_TYPE>_STATUS_MAPPING=<STANDARD_STATUS_1>:<ORIGIN_STATUS_1>,<ORIGIN_STATUS_2>;<STANDARD_STATUS_2>
JIRA_ISSUE_BUG_STATUS_MAPPING=Resolved:Approved,Verified,Done,Closed;Reject:ByDesign,Irreproducible
JIRA_ISSUE_INCIDENT_STATUS_MAPPING=Resolved:Done,Closed;Reject:ByDesign,Irreproducible
JIRA_ISSUE_STORY_STATUS_MAPPING=Resolved:Verified,Done,Closed;Reject:Abandoned,Cancelled

# Jira: epic issue #

JIRA_ISSUE_EPIC_KEY_FIELD=customfield_10014

# Jira: story point #

JIRA_ISSUE_STORYPOINT_COEFFICIENT=1
JIRA_ISSUE_STORYPOINT_FIELD=customfield_10024
```


## Issue status mapping<a id="issue-status-mapping"></a>
Jira is highly customizable, different company may use different `status name` to represent whether a issue was
resolved or not, one may named it "Done" and others might named it "Finished".
In order to collect life-cycle information correctly, you'll have to map your specific status to Devlake's standard
status, Devlake supports two standard status:

 - `Resolved`: issue was ended successfully
 - `Rejected`: issue was ended by termination or cancellation

Say we were using `Done` and `Cancelled` to represent the final stage of `Story` issues, what we have to do is setting
the following `Environment Variables` before running Devlake:
```sh
JIRA_ISSUE_STORY_STATUS_MAPPING=Resolved:Done;Reject:Cancelled
```


## Issue type mapping<a id="issue-type-mapping"></a>
Same as status mapping, different company might use different issue type to represent their Bug/Incident/Requirement,
type mapping is for Devlake to recognize your specific setup.
Devlake supports three different standard types:

 - `Bug`
 - `Incident`
 - `Requirement`

Say we were using `Story` to represent our Requirement, what we have to do is setting the following
`Environment Variables` before running Devlake:
```sh
JIRA_ISSUE_TYPE_MAPPING=Requirement:Story
```


## Find Board Id
1. Navigate to the Jira board in the browser
2. in the URL bar, get the board id from the parameter `?rapidView=`

**Example:**
`https://<your_jira_url>/secure/RapidBoard.jspa?rapidView=51`

![Screen Shot 2021-08-13 at 10 07 19 AM](https://user-images.githubusercontent.com/27032263/129363083-df0afa18-e147-4612-baf9-d284a8bb7a59.png)

> Use this board ID in your requests, to collect data from this board

## Generating API token
1. Once logged into Jira, visit the url `https://id.atlassian.com/manage-profile/security/api-tokens`
2. Click the **Create API Token** button, and give it any label name
![image](https://user-images.githubusercontent.com/27032263/129363611-af5077c9-7a27-474a-a685-4ad52366608b.png)
3. Encode with login email using command `echo -n <jira login email>:<jira token> | base64`

## How do I find the custom field ID in Jira?
Using URL
1. Navigate to Administration >> Issues >> Custom Fields .
2. Click the cog and hover over Configure or Screens option.
3. Observe the URL at the bottom left of the browser window. Example: The id for this custom field is 10006.


