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

##### Q1. How to create a Bitbucket app password?

1. Log in to [bitbucket.org](https://bitbucket.org).
2. Select the **Settings** cog in the upper-right corner of the top navigation bar.
3. Under **Personal settings**, select **Personal Bitbucket settings**.
4. On the left sidebar, select **App passwords**.
5. Select **Create app password**.
6. Give the App password a name.
7. Select the permissions the App password needs. See Q2.
8. Select the **Create** button.

Check [this doc](https://devlake.apache.org/docs/Configuration/BitBucket/#username-and-app-password) for more details.

##### Q2. What app password permission should I choose?

The following permissions are required to collect data from Bitbucket repositories:
`Account:Read` `Workspace` `membership:Read` `Repositories:Read` `Projects:Read` `Pull requests:Read` `Issues:Read` `Pipelines:Read` `Runners:Read`

##### Q3. Can I connect to the Bitbucket server?

Sure, you could

1. Go to the Connections page.
2. Click Create a New Connection.
3. Choose ‘Bitbucket server’ and finish the connection.

Bitbucket server is not supported to simplify this onboard process.
