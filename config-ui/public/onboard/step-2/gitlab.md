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

##### Q1. How to generate a GitLab personal access token?

1. Sign in at [gitlab.com](https://gitlab.com).
2. On the left sidebar, select your avatar, then select **Edit profile**.
3. On the left sidebar, select **Access Tokens**.
4. Select **Add new token**.
5. Enter a name and expiry date for the token. Select the desired scopes, `api` or `read_api`.
6. Select **Create personal access token**.

For detailed instructions, refer to [this doc](https://devlake.apache.org/docs/Configuration/GitLab/#personal-access-token).

##### Q2. Which scopes should be included in a token?

At least one of the following scopes must be included:
`api` `read_api`

Also, ensure proper user permissions are set for the GitLab project you intend to collect data from:

1. Navigate to the **Project information > Members** page of the GitLab project.
2. Confirm your role under the **Max role** column. Avoid the Guest role to ensure data collection capabilities.

##### Q3. Is connecting to the GitLab Server possible?

Yes, you can.

1. Navigate to the **Connections** page.
2. Click on GitLab, and click on **Create a New Connection**.
3. Select 'GitLab Server' and finish the configuration.

Please note that GitLab Server integration is not included to ease the onboarding process.
