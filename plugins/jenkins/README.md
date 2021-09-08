# Jenkins pond

## Metrics

Metric Name | Description
:------------ | :-------------
Build Count | The number of builds created
Build Success Rate | The percentage of successful builds

## Configuration

In your `.env` file, you will need to set up

```
# Jenkins configuration
JENKINS_ENDPOINT=https://jenkins.merico.cn/
JENKINS_USERNAME=your user name here
JENKINS_PASSWORD=your password or jenkins token here
```

You can generate access token at `User` -> `Configure` -> `API Token` section.

## Gathering data with jenkins

To collect data, you can make a POST request to `/task`

```
curl --location --request POST 'localhost:8080/task' \
  --header 'Content-Type: application/json' \
  --data-raw '[{
      "Plugin": "jenkins",
      "Options": {}
  }]'
```
