/*
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
*/

package plugin

// Scope represents the top level entity for a data source, i.e. github repo,
// gitlab project, jira board. They turn into repo, board in Domain Layer. In
// Apache Devlake, a Project is essentially a set of these top level entities,
// for the framework to maintain these relationships dynamically and
// automatically, all Domain Layer Top Level Entities should implement this
// interface
type Scope interface {
	ScopeId() string
	ScopeName() string
	TableName() string
}

type ToolLayerScope interface {
	ScopeId() string
	ScopeName() string
	TableName() string
}

type ApiScope interface {
	ConvertApiScope() ToolLayerScope
}

type ApiGroup interface {
	GroupId() string
	GroupName() string
}
