# AE (Merico Analysis Engine) Domain

## Summary

This plugin converts Merico AE data to [Domain Layer](../domainlayer/README.md) data

Specifically, it add a field called `dev_eq` to commits when run.


## How to trigger the conversion task
```
curl -XPOST 'localhost:8080/task' \
-H 'Content-Type: application/json' \
-d '[[{
    "plugin": "merico-analysis-engine-domain",
    "options": {}
}]]'
```
