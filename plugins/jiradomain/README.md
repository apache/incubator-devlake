# Jira Domain

## Summary

This plugin converts Jira data to [Domain Layer](../domainlayer/README.md) data


## How to trigger the conversion task
```
curl -XPOST 'localhost:8080/task' \
-H 'Content-Type: application/json' \
-d '[{
    "plugin": "jiradomain",
    "options": {
    }
}]'
```
