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

##### Q1. What Azure DevOps data will be collected?

`repos` `commits` `branches` `pull requests` `pr comments` `pipeline runs` `job runs` `users`

Check [this doc](https://devlake.apache.org/docs/Overview/SupportedDataSources/#data-collection-scope-by-each-plugin) for more details.

##### Q2. The data from which time range will be collected?

Only the data from last 14 days will be collected to speed up the sync up time. You can always change the time range in the project details page later.

##### Q3. Can I do transformations on the collected data?

Yes. You can do transformations by adding a Scope Config to the repositories you choose later.

##### Q4. What is the frequency to sync up data?

The data will be synced daily, you can change it in the project details page later.
