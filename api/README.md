# Lake Api

## Routes Available

### Create a task

```
curl --location --request POST 'localhost:8080/task' \
--header 'Content-Type: application/json' \
--data-raw '[
    {
        "Plugin": "jira",
        "Options": {
            "boardId": 8
        }
    }
]'
```

### Cancel a task

curl --location --request POST 'localhost:8080/task/cancel?taskName=jira' \
--header 'Content-Type: application/json'