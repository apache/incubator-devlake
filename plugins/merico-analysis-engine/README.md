# Merico Analysis Engine (AE)

THIS PLUGIN IS ONLY FOR MERICO EMPLOYEES AT THIS TIME. SOON IT WILL BE MADE PUBLIC.

## External Documentation

1. Swagger API Docs - {api host url}/docs#/default/list_projects_projects_get
2. AE Api Server - https://merico.feishu.cn/docs/doccnsJJsZEOZKFI5u7dif2NKWf#
3. Source Id and Source Type Architechture - https://merico.feishu.cn/docs/doccnLuIxqeE96L8SbbW1Tiqdmi#

## Important notes

### Some data looks like it is missing...

The commit data stored in Trino. The files can be deleted by Mino expiration strategy over time if they are too old.

### How do I trigger analysis on my project?

Just add DevLake to the Merico Enterprise Edition and triggered an analysis. You can find this item by searching "ae staging"? You can log in AE staging server and restart an analysis of DevLake. (Login credentials for Merico employees are stored in one password)

### Who controls the api for merico analysis engine?

Jingyang Liang and the Merico AE team

### How do I authenticate?

1. You need to be given an app_id and a secret_key by Merico
2. You need to send requests to the api in a special way. Each time you send a request, you need to encode a string called a sign.
3. You can make a `sign` by md5 encoding your query params and some other things...


- *nonce_str* Should be different between each request. It's ok to be a timestamp
- *app_id* Obtained from AE
- *secret_key* Obtained from AE

```
app_id={app_id}&key={secretKey}&nonce_str={timestamp}&page={page}&per_page={page_size}
```

!!!NOTE!!!: The keys have to be sorted in alphabetical order or it will not work!

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
