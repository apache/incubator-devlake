# Github Pond

## Metrics

Currently the data is only fetched and stored in the DB. Soon we will have charts in Grafana to support this data.

## Github rate limits

"For API requests using Basic Authentication or OAuth, you can make up to 5,000 requests per hour."

- https://docs.github.com/en/rest/overview/resources-in-the-rest-api

If you have a need for more api rate limits, you can set many tokens in the config file and we will use all of your tokens.

NOTE: You can get 15000 requests/hour/token if you pay for github enterprise.

## Configuration

In your .env file, you will need to set up

```
GITHUB_AUTH_TOKENS=XXX,YYY,ZZZ // where each token is a different user's token (optional)
GITHUB_AUTH=XXX
```

