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

##### Q1. How to create an Azure DevOps token?

1. Sign in to your organization (https://dev.azure.com/{yourorganization}) and go to the homepage.
2. Open **User Settings** in the upper-right corner of the top navigation bar and select **Personal access tokens**.
3. Select **New Token**.
4. Name your token, select 'All accessible organizations' in the Organization field.
5. Select the scopes. See Q2.
6. Select **Create**.

Check [this doc](https://devlake.apache.org/docs/Configuration/AzureDevOps/#token) for more details.

##### Q2. What token scopes should I choose?

Please select 'Full access'.

##### Q3. Can I connect to the Azure DevOps server?

No. Azure DevOps server is not supported yet.
