## Summary

GitHub has a rate limit of 5,000 API calls per hour for their REST API.
As a result, it may take hours to collect commits data from GitHub API for a repo that has 10,000+ commits.
To accelerate the process, DevLake introduces GitExtractor, a new plugin that collects git data by cloning the git repo instead of by calling GitHub APIs.

Starting from v0.10.0, DevLake will collect GitHub data in 2 separate plugins: 

- GitHub plugin (via GitHub API): collect repos, issues, pull requests
- GitExtractor (via cloning repos):  collect commits, refs

Note that GitLab plugin still collects commits via API by default since GitLab has a much higher API rate limit.

This doc details the process of collecting GitHub data in v0.10.0. We're working on simplifying this process in the next releases.

Before start, please make sure all services are started.

## Data Collection Procedure
There're 3 steps in total.

### Step 1 - Add a GitHub/GitLab connection
The following doc shows the steps for GitHub data collection. It's the same for GitLab.

1. Visit config-UI in localhost:4000 , click the GitHub icon

2. Click the default connection 'Github' in the list
    ![image](https://user-images.githubusercontent.com/14050754/163591959-11d83216-057b-429f-bb35-a9d845b3de5a.png)
    
3. Enter connection info
    ![image](https://user-images.githubusercontent.com/14050754/163592015-b3294437-ce39-45d6-adf6-293e620d3942.png)

    > Connection Name [READONLY] Defaults to "Github" and may not be changed
    >
    > Endpoint URL (REST URL, starts with https:// or http://)This should be a valid REST API Endpoint eg. https://api.github.com/
    >
    > URL should end with/
    >
    > Auth Token(s) (Personal Access Token)For help on Creating a personal access token, please see official GitHub Docs on Personal Tokens
    > Provide at least one token for Authentication with the . This field accepts a comma-separated list of values for multiple tokens. The data collection will take longer for GitHub since they have a rate limit of 2k requests per hour. You can accelerate the process by configuring multiple personal access tokens.
    >
    > GitHub Proxy URL [Optional] Enter a valid proxy server address on your Network, e.g. http://your-proxy-server.com:1080
    
4. Click 'Test Connection' and make sure it is working, and then click 'Save Connection'.

5. Add or edit data enrichment rules, we need it because Github has no built-in properties to indicate the type/component... of an issue/pull request. 
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
   
      6. **Incident**: Same as above, with `type` setting to `INCIDENT
   
         
   
6. Click 'Save Settings'

### Step 2 - Create a GitHub/GitLab  pipeline

To collect issues/pull requests from API

1. On `config-ui`, click 'Create Pipeline Run'.
   ![image](https://user-images.githubusercontent.com/14050754/163592542-8b9d86ae-4f16-492c-8f90-12f1e90c5772.png)
2. Toggle on GitHub/Gitlab, type in the repo you want to collect.![image](https://user-images.githubusercontent.com/14050754/163592606-92141c7e-e820-4644-b2c9-49aa44f10871.png)
3. Click Run Pipeline
   It'll jump to the Pipeline Activity page. 
   ![image](https://user-images.githubusercontent.com/14050754/163592677-268e6b77-db3f-4eec-8a0e-ced282f5a361.png)
   Please wait until the progress turns to 100%.
   ![image](https://user-images.githubusercontent.com/14050754/163592709-cce0d502-92e9-4c19-8504-6eb521b76169.png)

### Step 3 - Create a GitExtractor Pipeline
To collect commits, branch/tags by cloning the repo, 

simply 



1. Enable the `GitExtractor` plugin, and enter your `Git URL` and, select the `Repository ID` from dropdown menu.
   ![image-20220416154427292](GitHub-and-GitLab-Quick-Startup-Guide-v0.10.0.assets/image-20220416154427292.png)
2. Go ahead and Run the Pipeline, wait until it is done
3. Click `View Dashboards` on the top left corner
   ![image](https://user-images.githubusercontent.com/61080/163666814-e48ac68d-a0cc-4413-bed7-ba123dd291c8.png)
4. Check and verify if charts were working fine for you

## Create a Blueprint (recurring pipelines)

Now, let's assume that you are happy with what you see on grafana dashboard. Most likely, what you want next is to make DevLake collecting data periodically for you without manual operation, no worry, we've got you covered:

1. Click 'Create Pipeline Run'
  - Toggle on Github
  - Toggle on GitExtractor
  - Toggle on Automate Pipeline
    ![image](https://user-images.githubusercontent.com/14050754/163596590-484e4300-b17e-4119-9818-52463c10b889.png)


2. Click 'Add Blueprint'. Fill in the form and 'Save Blueprint'.
    
    - **NOTE**: That the schedule syntax is standard unix cron syntax, check [Crontab.guru](https://crontab.guru/) to learn more
    - **IMPORANT**: The scheduler is running under `UTC` timezone. If you prefer data collecting happens at 3am NewYork(UTC-04:00) every day, use **Custom Shedule** and set it to `0 7 * * *`
    
    ![image](https://user-images.githubusercontent.com/14050754/163596655-db59e154-405f-4739-89f2-7dceab7341fe.png)
    
3. Click 'Save Blueprint'.
    
4. Click 'Pipeline Blueprints', you can view and edit the new blueprint in the blueprint list.
    
    ![image](https://user-images.githubusercontent.com/14050754/163596773-4fb4237e-e3f2-4aef-993f-8a1499ca30e2.png)





