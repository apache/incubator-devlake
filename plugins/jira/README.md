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
#JIRA_ISSUE_<YOUR_TYPE>_STATUS_MAPPING=<STANDARD_STATUS>:<YOUR_STATUS>;...
JIRA_ISSUE_STORY_STATUS_MAPPING=Resolved:Done;Rejected:Cancelled
```

Status mapping is critical for metrics like **Lead Time**, the `leadtime`s are calculated only for those **Resolved**
issues.


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
# JIRA_ISSUE_TYPE_MAPPING=<STANDARD_TYPE>:<YOUR_TYPE_1>,<YOUR_TYPE_2>;....
JIRA_ISSUE_TYPE_MAPPING=Requirement:Story
```

Type mapping is critical for some metrics, like **Requirement Count**, make sure to map your custom type correctly.


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

Please follow this guide: [How to find Jira the custom field ID in Jira? · merico-dev/lake Wiki](https://github.com/merico-dev/lake/wiki/How-to-find-Jira-the-custom-field-ID-in-Jira%3F)
