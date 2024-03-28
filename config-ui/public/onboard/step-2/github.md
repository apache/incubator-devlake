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

##### Q1. How to generate a GitHub personal access token (classic)?

1. Sign in at [github.com](https://github.com).
2. Click your profile photo in the upper-right corner, select **Settings**.
3. On the left sidebar, click **Developer settings**.
4. Click **Tokens (classic)** under **Personal access tokens**.
5. Select the scopes: `repo:status`, `repo:deployment`, `read:user` and `read:org`.

For detailed instructions, refer to [this doc](https://devlake.apache.org/docs/Configuration/GitHub/#personal-access-tokens).

##### Q2. Which scopes should be included in a token?

Typically, the necessary scopes are:
`repo:status` `repo:deployment` `read:user` `read:org`

For private repositories, extend permissions with:
`repo` `read:user` `read:org`

##### Q3. Is connecting to the GitHub Server version possible?

Yes, you can.

1. Navigate to the **Connections** page.
2. Click on GitHub, and click on **Create a New Connection**.
3. Select 'GitHub Server' and finish the configuration.

Please note that GitHub Server integration is not included to ease the onboarding process.
