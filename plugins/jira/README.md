# Jira Pond

## Summary

Jira pond is a plugin used by [lake](https://github.com/merico-dev/lake/blob/main/README.md). The main thing it does is make api requests for you to Jira to fetch data and enrich it into a postgres database. Once this is done, you can use [Grafana](https://grafana.com/), hosted by [Lake](https://github.com/merico-dev/lake/blob/main/README.md) with [docker-compose](https://docs.docker.com/compose/install/)

Currently, this is how data flows:

<img width="830" alt="Screen Shot 2021-08-13 at 9 41 36 AM" src="https://user-images.githubusercontent.com/3011407/129358608-0f95beb3-7933-47d8-9775-65337d66fb1b.png">

And this is what you can expect to see for graphs in Grafana

<img width="1783" alt="Screen Shot 2021-08-13 at 9 42 41 AM" src="https://user-images.githubusercontent.com/3011407/129358727-503953a0-6a8d-43e3-b3d4-b379e40741cb.png">

## Project Metrics This Covers

Metric Name | Description
:------------ | :-------------
Requirement Count	| Number of issues with type "Requirement"
Requirement Lead Time	| Lead time of issues with type "Requirement"
Requirement Delivery Rate |	Ratio of delivered requirements to all requirements
Bug Count	| Number of issues with type "Bug"<br><i>bugs are found during testing</i>
Bug Age	| Lead time of issues with type "Bug"<br><i>both new and deleted lines count</i>
Bugs Count per 1k Lines of Code |	Amount of bugs per 1000 lines of code
Incident Count | Number of issues with type "Incident"<br><i>incidents are found when running in production</i>
Incident Age | Lead time of issues with type "Incident"
Incident Count per 1k Lines of Code | Amount of incidents per 1000 lines of code

## ⚠️ (WIP) Configuration

- [ ] Add in specific configuration options and steps with new go setup here

## ⚠️ (WIP)Gathering Data with Jira

- [ ] Add in specific steps to collect data from Jira with new go setup here

## Find Board Id
1. Navigate to the Jira board in the browser
2. in the URL bar, get the board id from the parameter `?rapidView=`

**Example:**
`https://<your_jira_url>/secure/RapidBoard.jspa?rapidView=51`

![Screen Shot 2021-08-13 at 10 07 19 AM](https://user-images.githubusercontent.com/27032263/129363083-df0afa18-e147-4612-baf9-d284a8bb7a59.png)

> Use this board ID in your requests, to collect data from this board

## ⚠️ (WIP) Generating API token
1. Once logged into Jira, visit the url `https://id.atlassian.com/manage-profile/security/api-tokens`
2. Click the **Create API Token** button, and give it any label name

![image](https://user-images.githubusercontent.com/27032263/129363611-af5077c9-7a27-474a-a685-4ad52366608b.png)

- [ ] Copy and save the API token string into `lake` via new go setup

## ⚠️ (WIP) Jira Specific String Configuration

Adjust what is considered "Bug", "Incident" or "Requirement".

- [ ] Add in options for specific configuration in new go setup
