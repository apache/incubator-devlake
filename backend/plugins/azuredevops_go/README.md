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

# Azure Devops Python Plugin

This is a revamped version of the Python Azure DevOps Plugin, originally located at `../../python/plugins/azuredevops`.
The plugin is able to coexist with the Python version as both implementations come with their own `_raw` and `_tool` tables.

**Read access** to the following Azure DevOps Scopes is required:

- Build
- Code
- Graph (collectAccounts task)
- Release

Access to Service Connections has been removed as they usually contain sensitive security information.