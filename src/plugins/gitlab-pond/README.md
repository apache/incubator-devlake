# Gitlab Pond

## Summary

Gitlab pond is a plugin used by [lake](https://github.com/merico-dev/lake/blob/main/README.md). The main thing it does is make api requests for you to Gitlab to fetch data and enrich it into a postgres database. Once this is done, you can use [Grafana](https://grafana.com/), hosted by [Lake](https://github.com/merico-dev/lake/blob/main/README.md) with [docker-compose](https://docs.docker.com/compose/install/)

Currently, this is how data flows:

<img width="830" alt="Screen Shot 2021-08-13 at 9 41 36 AM" src="https://user-images.githubusercontent.com/3011407/129358608-0f95beb3-7933-47d8-9775-65337d66fb1b.png">

And this is what you can expect to see for graphs in Grafana

<img width="1783" alt="Screen Shot 2021-08-13 at 9 42 41 AM" src="https://user-images.githubusercontent.com/3011407/129358727-503953a0-6a8d-43e3-b3d4-b379e40741cb.png">

## Project Metrics This Covers

Metric Name | Description
:------------ | :-------------
Pull Request Count | Number of Pull/Merge Requests
Pull Request Pass Rate | Ratio of Pull/Merge Review requests to merged
Pull Request Reviewer Count | Number of Pull/Merge Reviewers
Pull Request Review Time | Time from the first Pull/Merge Review comment until merged
Commit Author Count | Number of Contributors
Commit Count | Number of Commits
Added Lines | Accumulated Number of New Lines
Deleted Lines | Accumulated Number of Removed Lines

## Configuration

First, you have to configure `lake/config/plugins.js`.

1. `cp config/plugins.sample.js config/plugins.js`
2. Look though this config file to make sure it is set to your needs. `cat config/plugins.js`

## Gathering Data with Gitlab

Once you have [lake](https://github.com/merico-dev/lake/blob/main/README.md) running, you can fetch information from Github in one two ways:

1. Send a POST request to http://localhost:3001/
```
 {
     "gitlab": {
         "projectId": 8967944,
         "branch": "<your-branch-name>", (Optional, default branch is used)
     }
 }
```
2. You can configure lake to get all of your data automatically.

- Make sure you have a file called `config/local.js`
- `cp config/local.sample.js config/local.js`
- Open this file for editing with your editor of choice or use `vi config/local.js`
- In this file, there is a section for cron.
- Set the projectId fo your own project Id in gitlab.

```
gitlab: {
  projectId: 123
}
```

NOTE: If you don't know how to find the projectId, see the section below :)

## Finding Project Id
To get the project id for a specific Gitlab repository:
- Visit the repository page on gitlab
- Find the project id just below the title

  ![Screen Shot 2021-08-06 at 4 32 53 PM](https://user-images.githubusercontent.com/3789273/128568416-a47b2763-51d8-4a6a-8a8b-396512bffb03.png)

> Use this project id in your requests, to collect data from this project

## Create a Gitlab API Token

1. When logged into Gitlab visit `https://gitlab.com/-/profile/personal_access_tokens`
2. Give the token any name, no expiration date and all scopes (excluding write access)

    ![Screen Shot 2021-08-06 at 4 44 01 PM](https://user-images.githubusercontent.com/3789273/128569148-96f50d4e-5b3b-4110-af69-a68f8d64350a.png)

3. Click the **Create Personal Access Token** button
4. Copy the token into the `lake` plugin config file `config/plugins.js`
