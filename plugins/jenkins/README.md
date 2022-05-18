# Jenkins

## Summary

This plugin collects Jenkins data through [Remote Access API](https://www.jenkins.io/doc/book/using/remote-access-api/). It then computes and visualizes various devops metrics from the Jenkins data.

![image](https://user-images.githubusercontent.com/61080/141943122-dcb08c35-cb68-4967-9a7c-87b63c2d6988.png)

## Metrics

| Metric Name        | Description                         |
|:-------------------|:------------------------------------|
| Build Count        | The number of builds created        |
| Build Success Rate | The percentage of successful builds |

## Configuration

In order to fully use this plugin, you will need to set various configurations via Dev Lake's `config-ui`.

### By `config-ui`

The connection aspect of the configuration screen requires the following key fields to connect to the Jenkins API. As Jenkins is a single-source data provider at the moment, the connection name is read-only as there is only one instance to manage. As we continue our development roadmap we may enable multi-source connections for Jenkins in the future.

- Connection Name [READONLY]
  - ⚠️ Defaults to "Jenkins" and may not be changed.
- Endpoint URL (REST URL, starts with `https://` or `http://`i, ends with `/`)
  - This should be a valid REST API Endpoint eg. `https://ci.jenkins.io/`
- Username (E-mail)
  - Your User ID for the Jenkins Instance.
- Password (Secret Phrase or API Access Token)
  - Secret password for common credentials.
  - For help on Username and Password, please see official Jenkins Docs on Using Credentials
  - Or you can use **API Access Token** for this field, which can be generated at `User` -> `Configure` -> `API Token` section on Jenkins.

Click Save Connection to update connection settings.

## Collect Data From Jenkins

In order to collect data from Jenkins, you have to compose a JSON looks like following one, and send it via `Triggers` page on `config-ui`:


```json
[
  [
    {
      "plugin": "jenkins",
      "options": {}
    }
  ]
]
```

## Relationship between job and build

Build is kind of a snapshot of job. Running job each time creates a build.