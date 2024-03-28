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

##### Q1. How to create a GitLab personal access token?

1. Log in to [gitlab.com](https://gitlab.com).
2. On the left sidebar, select your avatar, then select **Edit profile**.
3. On the left sidebar, select **Access Tokens**.
4. Select **Add new token**.
5. Enter a name and expiry date for the token. Select the desired scopes, `api` or `read_api`.
6. Select **Create personal access token**.

Check [this doc](https://devlake.apache.org/docs/Configuration/GitLab/#personal-access-token) for more details.

##### Q2.What token scopes should I choose?

Only one of the following scopes is required:
`api` `read_api`

You also have to double-check your GitLab user
permission settings:

1. Go to the Project information > Members page of the GitLab projects you wish to collect.
2. Check your role in this project from the Max role column. Make sure you are not the Guest role, otherwise, you will not be able to collect data from this project.

##### Q3. Can I connect to the GitLab server?

Sure, you could

1. Go to the Connections page.
2. Click Create a New Connection.
3. Choose 'GitLab server' and finish the connection.

GitLab server is not supported to simplify this onboard process.
