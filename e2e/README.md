## E2E Tests

## Why?

To ensure data integrity of the application, we need to make sure that the result
data matches what we expect to get from accessing a real API as if a real user is
using it. Previously, we relied on manual testing.

Automated tests allow us to do this in a very convenient, low cost way 
that is easily repeatable.

## What it does

1. Automatically or Manually trigger all collection / enrichment / conversion tasks
2. Tests access all key data models from our DB to determine if the expected number 
of rows were collected and processed or not.

## How to run

1. Collect all data normally for all repos. 
2. Wait until all jobs are done
3. Then you can run tests with this command: `make real-e2e-test`

## JSON samples to send to DevLake (POST /pipelines)

### Jira

{
  "name": "jira",
  "tasks": [[{
    "Plugin": "jira",
    "Options": {
        "sourceId": `<your_source_id>`,
        "boardId": 8
    }
  }]]
}

### GitLab

{
  "name": "gitlab",
  "tasks": [[{
    "Plugin": "gitlab",
    "Options": {
        "projectId": 8967944
    }
  }]]
}

### GitHub

{
  "name": "github",
  "tasks": [[{
    "Plugin": "github",
    "Options": {
        "repo": "lake",
        "owner": "merico-dev"
      }
  }]]
}

### Jenkins

{
  "name": "jenkins",
  "tasks": [[
    {
        "plugin": "jenkins",
        "options": {}
    }
  ]]
}