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

##### Q1. How to generate a Bitbucket API token?

**⚠️ Important: App passwords are being deprecated!**
- Creation of App passwords will be discontinued on **September 9, 2025**
- All existing App passwords will be deactivated on **June 9, 2026**
- Please use API tokens for all new connections

**To create an API token:**

1. Sign in at [https://id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens).
2. Select **Create API token with scopes**.
3. Give the API token a name and an expiry date (ex: 365 days), then select **Next**.
4. Select **Bitbucket** as the app and select **Next**.
5. Select the required scopes (see **Q2**) and select **Next**.
6. Review your token and select **Create token**.
7. **Copy the generated API token immediately** - it's only displayed once and can't be retrieved later.

For detailed instructions, refer to [Atlassian's API token documentation](https://support.atlassian.com/bitbucket-cloud/docs/create-an-api-token/).

##### Q2. Which permissions (scopes) should be included in an API token?

The following scopes are **required** to collect data from Bitbucket repositories:

- `read:account` - Required to view users profiles
- `read:issue:bitbucket` - View your issues
- `read:pipeline:bitbucket` - View your pipelines
- `read:project:bitbucket` - View your projects
- `read:pullrequest:bitbucket` - View your pull requests
- `read:repository:bitbucket` - View your repositories
- `read:runner:bitbucket` - View your workspaces/repositories' runners
- `read:user:bitbucket` - View user info (required for connection test)
- `read:workspace:bitbucket` - View your workspaces

##### Q3. Is connecting to the Bitbucket Server/Data Center possible?

Yes, you can.

1. Navigate to the **Connections** page.
2. Click on Bitbucket Server, and Click on **Create a New Connection**.
3. Finish the configuration.

Please note that GitLab Server integration is not included to ease the onboarding process.
