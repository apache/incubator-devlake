# Github Pond

<div align="center">

| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |

</div>

<br>

## Summary

This plugin gathers data from `GitHub` to display information to the user in `Grafana`. We can help tech leaders answer such questions as:

- Is this month more productive than last?
- How fast do we respond to customer requirements?
- Was our quality improved or not?

## Metrics

Here are some examples of what we can use `GitHub` data to show:
- Avg Requirement Lead Time By Assignee
- Bug Count per 1k Lines of Code
- Commit Count over Time

## Screenshot

![image](https://user-images.githubusercontent.com/27032263/141855099-f218f220-1707-45fa-aced-6742ab4c4286.png)


## Configuration

### Provider (Datasource) Connection
The connection aspect of the configuration screen requires the following key fields to connect to the **GitHub API**. As GitHub is a _single-source data provider_ at the moment, the connection name is read-only as there is only one instance to manage. As we continue our development roadmap we may enable _multi-source_ connections for GitHub in the future.

- **Connection Name** [`READONLY`]
  - ⚠️ Defaults to "**Github**" and may not be changed.
- **Endpoint URL** (REST URL, starts with `https://` or `http://`)
  - This should be a valid REST API Endpoint eg. `https://api.github.com/`
  - ⚠️ URL should end with`/`
- **Auth Token(s)** (Personal Access Token)
  - For help on **Creating a personal access token**, please see official [GitHub Docs on Personal Tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
  - Provide at least one token for Authentication with the . This field accepts a comma-separated list of values for multiple tokens. The data collection will take longer for GitHub since they have a **rate limit of 2k requests per hour**. You can accelerate the process by configuring _multiple_ personal access tokens.
    
"For API requests using `Basic Authentication` or `OAuth`, you can make up to 5,000 requests per hour."

- https://docs.github.com/en/rest/overview/resources-in-the-rest-api

If you have a need for more api rate limits, you can set many tokens in the config file and we will use all of your tokens.

NOTE: You can get 15000 requests/hour/token if you pay for `GitHub` enterprise.
    
For an overview of the **GitHub REST API**, please see official [GitHub Docs on REST](https://docs.github.com/en/rest)
    
Click **Save Connection** to update connection settings.
    

### Provider (Datasource) Settings
Manage additional settings and options for the GitHub Datasource Provider. Currently there is only one **optional** setting, *Proxy URL*. If you are behind a corporate firewall or VPN you may need to utilize a proxy server.

**GitHub Proxy URL [ `Optional`]**
Enter a valid proxy server address on your Network, e.g. `http://your-proxy-server.com:1080`

Click **Save Settings** to update additional settings.

## Sample Request

```
curl --location --request POST 'localhost:8080/pipelines' \
--header 'Content-Type: application/json' \
--data-raw '
{
    "name": "github 20211126",
    "tasks": [[{
        "plugin": "github",
        "options": {
            "repositoryName": "lake",
            "owner": "merico-dev"
        }
    }]]
}
'
```
