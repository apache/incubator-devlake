# Merico Analysis Engine (AE)

## External Documentation

1. Swagger API Docs - http://34.214.122.134:30012/docs#/default/list_projects_projects_get
2. AE Api Server - https://merico.feishu.cn/docs/doccnsJJsZEOZKFI5u7dif2NKWf#
3. Source Id and Source Type Architechture - https://merico.feishu.cn/docs/doccnLuIxqeE96L8SbbW1Tiqdmi#

## Data Gathered

*Projects*

```
[
  {
    "id": 0,
    "git_url": "string",
    "priority": 0,
    "create_time": "2021-11-23T17:28:10.286Z",
    "update_time": "2021-11-23T17:28:10.286Z"
  }
]
```

*Commits*

```
[
  {
    "hexsha": "string",
    "analysis_id": "string",
    "author_email": "string",
    "dev_eq": 0
  }
]
```

The most valuable data here is the dev_eq. This is a Merico owned measurement of code value

## Configuration

You will need to have two tokens in order to run this plugin.

These can be set in your .env file as

```
AE_SIGN=XXX
AE_NONCE=XXX
```

TBD: How do non merico users get these keys?

## Gathering Data with AE

To collect data on a single project, you can make a POST request to `/task`

    ```
    curl --location --request POST 'localhost:8080/task' \
    --header 'Content-Type: application/json' \
    --data-raw '[[{
        "plugin": "ae",
        "options": {
            "projectId": <Your project id>
        }
    }]]'
    ```
