# Jenkins Domain

## Summary

This plugin converts Jenkins data to [Domain Layer](../domainlayer/README.md) data


## How to trigger the conversion task
```
curl -XPOST 'localhost:8080/task' \
-H 'Content-Type: application/json' \
-d '[[{
    "plugin": "jenkinsdomain",
    "options": {
    }
}]]'
```
