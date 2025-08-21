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

##### Q1. How to generate a Bitbucket app password?

1. Sign in at [bitbucket.org](https://bitbucket.org).
2. Select the **Settings** cog in the upper-right corner of the top navigation bar.
3. Under **Personal settings**, select **Personal Bitbucket settings**.
4. On the left sidebar, select **App passwords**.
5. Select **Create app password**.
6. Give the 'App password' a name.
7. Select the permissions the 'App password needs'. See **Q2**.
8. Select the **Create** button.

For detailed instructions, refer to [this doc](https://devlake.apache.org/docs/Configuration/BitBucket/#username-and-app-password).

##### Q2. Which app password permissions should be included in a token?

The following permissions are required to collect data from Bitbucket repositories:
`Account:Read` `Workspace` `membership:Read` `Repositories:Read` `Projects:Read` `Pull requests:Read` `Issues:Read` `Pipelines:Read` `Runners:Read`

##### Q3. Is connecting to the Bitbucket Server/Data Center possible?

Yes, you can.

1. Navigate to the **Connections** page.
2. Click on Bitbucket Server, and Click on **Create a New Connection**.
3. Finish the configuration.

Please note that GitLab Server integration is not included to ease the onboarding process.
