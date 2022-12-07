## Push API

## Summary

This is a generic API service that gives our users the ability to inject data directly to their own database using a
simple, all-purpose endpoint.

## The Endpoint

POST to ```localhost:8080/push/:tableName```

Where "tableName" is the name of the table you wish to insert into
For example, "commits" would be ```/push/commits```

## The JSON body

Include a JSON body that consists of an array of objects you wish to insert.
Please Note: You must know the schema you are inserting into (column names, types, etc.)

```
[
    {
        "id": "gitlab...etc",
        "sha": "osidjfoawehfwh08",
        "additions": 89,
        ...
    }
]
```



