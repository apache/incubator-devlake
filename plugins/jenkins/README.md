# Jenkins

## Summary

This plugin collects Jenkins data through [Remote Access API](https://www.jenkins.io/doc/book/using/remote-access-api/). It then computes and visualizes various devops metrics from the Jenkins data.

![image](https://user-images.githubusercontent.com/61080/141943122-dcb08c35-cb68-4967-9a7c-87b63c2d6988.png)

## Metrics

Metric Name | Description
:------------ | :-------------
Build Count | The number of builds created
Build Success Rate | The percentage of successful builds

## Configuration

In order to fully use this plugin, you will need to set various configurations.
Either by `.env` or via Dev Lake's configuration UI.

### By `.env` file directly

In your `.env` file, you will need to set up
```
# Jenkins configuration
JENKINS_ENDPOINT=https://jenkins.merico.cn/
JENKINS_USERNAME=your user name here
JENKINS_PASSWORD=your password or jenkins token here
```

### By `config-ui`

Open config-ui page (default: http://localhost:4000), go to **Data Integrations** click on **Jenkins**, then **Settings**, enter your configuration and click **Save** button

### How to generate token?
You can generate your Jenkins access token at `User` -> `Configure` -> `API Token` section on Jenkins.

## Collect Data From Jenkins

In order to collect data from JIRA, you have to compose a JSON looks like following one, and send it via `Triggers` page on `config-ui`:


```json
[
  [
    {
      "plugin": "jenkins",
      "options": {}
    }
  ]
]'
```
