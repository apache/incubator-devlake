## pipeline-plan

### Summary

Users can set pipepline plan by config-ui to create schedule jobs.
And config-ui will send pipeline-plan request with cronConfig in crontab format.

### Cron Job

cronConfig should look like this: "M H D M WD"

M: minute

H: hour

D: day(month)

M: month

WD: day(week)

Please check cron time format in https://crontab.guru/

### API

POST /pipeline-plans
```json
Request
{
	"name": "COLLECT 1648121282469",
	"tasks": [
		[
			{
				"plugin": "github",
				"options": {
					"repo": "lake",
					"owner": "merico-dev"
				}
			}
		]
	],
	"enable": true,
	"cronConfig": "103 13 /13 * *"
}
Response
{
	"id": 7,
	"createdAt": "2022-03-27T10:16:20.046+08:00",
	"updatedAt": "2022-03-27T10:16:20.046+08:00",
	"name": "COLLECT 1648121282469",
	"tasks": [
		[
			{
				"plugin": "github",
				"options": {
					"owner": "merico-dev",
					"repo": "lake"
				}
			}
		]
	],
	"enable": true,
	"cronConfig": "103 13 /13 * *"
}
```
GET /pipeline-plans
```json
Response
{
	"id": 7,
	"createdAt": "2022-03-27T10:16:20.046+08:00",
	"updatedAt": "2022-03-27T10:16:20.046+08:00",
	"name": "COLLECT 1648121282469",
	"tasks": [
		[
			{
				"plugin": "github",
				"options": {
					"owner": "merico-dev",
					"repo": "lake"
				}
			}
		]
	],
	"enable": true,
	"cronConfig": "103 13 /13 * *"

}
```

GET /pipeline-plans/:pipelinePlanId
```json
Response
{
	"id": 7,
	"createdAt": "2022-03-27T10:16:20.046+08:00",
	"updatedAt": "2022-03-27T10:16:20.046+08:00",
	"name": "COLLECT 1648121282469",
	"tasks": [
		[
			{
				"plugin": "github",
				"options": {
					"owner": "merico-dev",
					"repo": "lake"
				}
			}
		]
	],
	"enable": true,
	"cronConfig": "103 13 /13 * *"
}
```


PATCH /pipeline-plans/:pipelinePlanId
```json
Request
{
	"name": "COLLECT 1648121282469",
	"tasks": [
		[
			{
				"plugin": "github",
				"options": {
					"repo": "lake",
					"owner": "merico-dev"
				}
			}
		]
	],
	"enable": true,
	"cronConfig": "103 13 /13 * *"
}
Response
{
	"id": 7,
	"createdAt": "2022-03-27T10:16:20.046+08:00",
	"updatedAt": "2022-03-27T10:16:20.046+08:00",
	"name": "COLLECT 1648121282469",
	"tasks": [
		[
			{
				"plugin": "github",
				"options": {
					"owner": "merico-dev",
					"repo": "lake"
				}
			}
		]
	],
	"enable": true,
	"cronConfig": "103 13 /13 * *"
}
```
