## Push API

## Summary

This is an optional generic API service that gives our users the ability to inject data directly to their own database using a simple, all-purpose endpoint.

## To Run

```
cd ./push-api
make build-and-run
```

## The Endpoint

POST to ```localhost:9123/api/:tableName```

Where "tableName" is the name of the table you wish to insert into 
For example, "commits" would be ```/api/commits```

## The JSON body

Include a JSON body that consists of an array of objects you wish to insert.
Please Note: You must know the schema you are inserting into (column names, types, etc.)
```
	[
		{
			"origin_key": "gitlab...etc",
			"sha": "osidjfoawehfwh08",
      "additions": 89,
      ...
		}
	]
```

## Docker Implementation - Coming Soon

Steps To Do:

- Fix Dockerfile
- Test image
- Add to docker-compose.yml


