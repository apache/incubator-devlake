# Github Pond

## Summary

This plugin gathers data from GitHub to display information to the user in Grafana. We can help tech leaders answer such questions as:

- Is this month more productive than last?
- How fast do we respond to customer requirements?
- Was our quality improved or not?

## Metrics

Here are some examples of what we can use GitHub data to show:
- Avg Requirement Lead Time By Assignee
- Bug Count per 1k Lines of Code
- Commit Count over Time

## Getting Personal Access Token

Here is a link to a guide to get your personal access token from GitHub:

https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token

## Github rate limits

"For API requests using Basic Authentication or OAuth, you can make up to 5,000 requests per hour."

- https://docs.github.com/en/rest/overview/resources-in-the-rest-api

If you have a need for more api rate limits, you can set many tokens in the config file and we will use all of your tokens.

NOTE: You can get 15000 requests/hour/token if you pay for github enterprise.

## Configuration

In your .env file, you will need to set up

```

GITHUB_AUTH=XXX

or...

GITHUB_AUTH=XXX,YYY,ZZZ // where each token is a different user's token (optional)
```

The proxy server address could be set in the `.env` file with the key `GITHUB_PROXY`. 
If the key is empty or any other invalid url, no proxy applied. Only `http` and `socks5` protocol supported for now.

```
GITHUB_PROXY=http://127.0.0.1:1080
```


## Sample Request

```
curl --location --request POST 'localhost:8080/task' \
--header 'Content-Type: application/json' \
--data-raw '[[{
    "Plugin": "github",
    "Options": {
        "repositoryName": "lake",
        "owner": "merico-dev"
    }
}]]'
```
