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
Please note, that the current implementation lacks a database migration script to retain the data collected by the
Python plugin. As a result, the plugin has been temporarily excluded from the build process to prevent conflicts with
the existing Python version.

Please be careful and only activate the plugin if you know what you do. Once activated following tool tables
will be dropped. The data aggregated in the domain tables will remain untouched which can leave your database in an
undesired state.

```
"_raw_azuredevops_builds",
"_raw_azuredevops_gitpullrequestcommits",
"_raw_azuredevops_gitpullrequests",
"_raw_azuredevops_jobs"

"_tool_azuredevops_builds",
"_tool_azuredevops_gitpullrequestcommits",
"_tool_azuredevops_gitpullrequests",
"_tool_azuredevops_gitrepositories",
"_tool_azuredevops_gitrepositoryconfigs",
"_tool_azuredevops_jobs"
```

To enable the plugin navigate to the build-plugin script and remove the exclude flag `-not -name azuredevops`.

```
PLUGINS=$(find $PLUGIN_SRC_DIR/* -maxdepth 0 -type d -not -name core -not -name helper -not -name logs -not -empty -not -name azuredevops)
```

The plugin requires **read access** to the following Azure DevOps Scopes:

- Build
- Code
- Graph (collectAccounts task)
- Release

Access to Service Connections has been removed as they usually contain sensitive security information.