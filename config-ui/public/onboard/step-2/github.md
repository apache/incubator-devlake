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

##### Q1. How to create a classic GitHub personal access token?

1. Log in to [github.com](https://github.com).
2. In the upper-right corner of any page, click your profile photo, then click Settings.
3. In the left sidebar, click <> Developer settings.
4. In the left sidebar, under 'Personal access tokens', click Tokens (classic).
5. Choose the following scopes: repo:status, repo:deployment, read:user and read:org.

Check [this doc](https://devlake.apache.org/docs/Configuration/GitHub/#personal-access-tokens) for more details.

##### Q2. What token scopes should I choose?

Normally, only the following scopes are required:
`repo:status` `repo:deployment` `read:user` `read:org`

However, if you want to collect data from private repositories, you need to give full permission to repo:
`repo` `read:user` `read:org`

##### Q3. Can I connect to the GitHub server?

Sure, you could

1. Go to the Connections page.
2. Click Create a New Connection.
3. Choose 'GitHub server' and finish the connection.

GitHub server is not supported to simplify this onboard process.
