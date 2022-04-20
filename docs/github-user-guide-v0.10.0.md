## Summary

GitHub has a rate limit of 2,000 API calls per hour for their REST API.
As a result, it may take hours to collect commits data from GitHub API for a repo that has 10,000+ commits.
To accelerate the process, DevLake introduces GitExtractor, a new plugin that collects git data by cloning the git repo instead of by calling GitHub APIs.

Starting from v0.10.0, DevLake will collect GitHub data in 2 separate plugins: 

- GitHub plugin (via GitHub API): collect repos, issues, pull requests
- GitExtractor (via cloning repos):  collect commits, refs

Note that GitLab plugin still collects commits via API by default since GitLab has a much higher API rate limit.

This doc details the process of collecting GitHub data in v0.10.0. We're working on simplifying this process in the next releases.

Before start, please make sure all services are started.

## GitHub Data Collection Procedure

There're 3 steps.

1. Configure GitHub connection
2. Create a pipeline to run GitHub plugin
3. Create a pipeline to run GitExtractor plugin
4. [Optional] Set up a recurring pipeline to keep data fresh

### Step 1 - Configure GitHub connection

1. Visit `config-ui` at `http://localhost:4000`, click the GitHub icon

2. Click the default connection 'Github' in the list
    ![image](https://user-images.githubusercontent.com/14050754/163591959-11d83216-057b-429f-bb35-a9d845b3de5a.png)
    
3. Configure connection by providing your GitHub API endpoint URL and your personal access token(s).
    ![image](https://user-images.githubusercontent.com/14050754/163592015-b3294437-ce39-45d6-adf6-293e620d3942.png)

    > Endpoint URL: Leave this unchanged if you're using github.com. Otherwise replace it with your own GitHub instance's REST API endpoint URL. This URL should end with '/'.
    >
    > Auth Token(s): Fill in your personal access tokens(s). For how to generate personal access tokens, please see GitHub's [official documentation](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).
    > You can provide multiple tokens to speed up the data collection process, simply concatenating tokens with commas.
    >
    > GitHub Proxy URL: This is optional. Enter a valid proxy server address on your Network, e.g. http://your-proxy-server.com:1080
    
4. Click 'Test Connection' and see it's working, then click 'Save Connection'.

5. [Optional] Help DevLake understand your GitHub data by customizing data enrichment rules shown below.
    ![image](https://user-images.githubusercontent.com/14050754/163592506-1873bdd1-53cb-413b-a528-7bda440d07c5.png)
   
   1. Pull Request Enrichment Options
   
      1. `Type`: PRs with label that matches given Regular Expression, their properties `type` will be set to the value of first sub match. For example, with Type being set to `type/(.*)$`, a PR with label `type/bug`, its `type` would be set to `bug`, with label `type/doc`, it would be `doc`.
      2. `Component`: Same as above, but for `component` property.
   
   2. Issue Enrichment Options
   
      1. `Severity`: Same as above, but for `issue.severity` of course.
   
      2. `Component`: Same as above.
   
      3. `Priority`: Same as above.
   
      4. **Requirement** : Issues with label that matches given Regular Expression, their properties `type` will be set to `REQUIREMENT`. Unlike `PR.type`, submatch does nothing,    because for Issue Management Analysis, people tend to focus on 3 kinds of type (Requiremnt/Bug/Incident), however, the concrete naming varies from repo to repo, time to time, so we decided to standardize them to help analyst making general purpose metric. 
   
      5. **Bug**: Same as above, with `type` setting to `BUG`
   
      6. **Incident**: Same as above, with `type` setting to `INCIDENT`
   
6. Click 'Save Settings'

### Step 2 - Create a pipeline to collect GitHub data

1. Select 'Pipelines > Create Pipeline Run' from `config-ui`

![image](https://user-images.githubusercontent.com/14050754/163592542-8b9d86ae-4f16-492c-8f90-12f1e90c5772.png)

2. Toggle on GitHub plugin, enter the repo you'd like to collect data from.

![image](https://user-images.githubusercontent.com/14050754/163592606-92141c7e-e820-4644-b2c9-49aa44f10871.png)

3. Click 'Run Pipeline'

You'll be redirected to newly created pipeline:

![image](https://user-images.githubusercontent.com/14050754/163592677-268e6b77-db3f-4eec-8a0e-ced282f5a361.png)


See the pipeline finishes (progress 100%):

![image](https://user-images.githubusercontent.com/14050754/163592709-cce0d502-92e9-4c19-8504-6eb521b76169.png)

### Step 3 - Create a pipeline to run GitExtractor plugin

1. Enable the `GitExtractor` plugin, and enter your `Git URL` and, select the `Repository ID` from dropdown menu.

![image](https://user-images.githubusercontent.com/2908155/164125950-37822d7f-6ee3-425d-8523-6f6b6213cb89.png)

2. Click 'Run Pipeline' and wait until it's finished.

3. Click `View Dashboards` on the top left corner of `config-ui`

![image](https://user-images.githubusercontent.com/61080/163666814-e48ac68d-a0cc-4413-bed7-ba123dd291c8.png)

4. See dashboards populated with GitHub data.

### Step 4 - [Optional] Set up a recurring pipeline to keep data fresh

Please see [How to create recurring pipelines](./recurring-pipeline.md) for details.






