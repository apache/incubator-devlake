# Gitlab Pond

## Metrics

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

Note: the following instructions are for *User Setup*. For *Developer Setup*, simply replace `config/docker.sample.js` and `config/docker.js` with `config/local.sample.js` and `config/local.js`.

- Make sure you have a file called `config/docker.js`. You can create one from the sample file: `cp config/docker.sample.js config/docker.js`
- Open this file for editing with your editor of choice or use `vi config/local.js`
- In this file, there is a section for cron.
- Set the projectId fo your own project Id in gitlab.

```
gitlab: {
  projectId: 123
}
```

- Restart lake services for the new configurtion to take effect immediately

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
