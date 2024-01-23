<!--
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->
# Trello

## Summary

This plugin collects `Trello` data through [Trello's rest api](https://developer.atlassian.com/cloud/trello/guides/rest-api/api-introduction/).

## Configuration

In order to fully use this plugin, you will need to get `apikey` and `token` on the [Trello website](https://developer.atlassian.com/cloud/trello/guides/rest-api/api-introduction/).

A connection should be created before you can collect any data. Currently, this plugin supports creating connection by requesting `connections` API:

```
curl 'http://localhost:8080/plugins/trello/connections' \
--header 'Content-Type: application/json' \
--data-raw '
{
    "name": "trello",
    "endpoint": "https://api.trello.com/",
    "rateLimitPerHour": 20000,
    "appId": "<YOUR_APIKEY>",
    "secretKey": "<YOUR_TOKEN>"
}
'
```

## Collect data from Trello

In order to collect data, you have to make a POST request to `/pipelines`.

```
curl 'http://localhost:8080/pipelines' \
--header 'Content-Type: application/json' \
--data-raw '
{
    "name":"MY PIPELINE",
    "plan":[
        [
            {
                "plugin":"trello",
                "options":{
                    "connectionId":<CONNECTION_ID>,
                    "boardId":"<BOARD_ID>"
                }
            }
        ]
    ]
}
'
```

You can make the following request to get all the boards.

```
curl 'http://localhost:8080/plugins/trello/connections/<CONNECTION_ID>/proxy/rest/1/members/me/boards?fields=name,id'
```
