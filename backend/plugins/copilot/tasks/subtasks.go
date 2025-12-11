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

package tasks

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var CollectCopilotOrgMetricsMeta = plugin.SubTaskMeta{
	Name:             "collectCopilotOrgMetrics",
	EntryPoint:       CollectCopilotOrgMetrics,
	EnabledByDefault: true,
	Description:      "Collect GitHub Copilot organization metrics",
}

var CollectCopilotSeatAssignmentsMeta = plugin.SubTaskMeta{
	Name:             "collectCopilotSeatAssignments",
	EntryPoint:       CollectCopilotSeatAssignments,
	EnabledByDefault: true,
	Description:      "Collect GitHub Copilot seat assignments",
}

var ExtractCopilotOrgMetricsMeta = plugin.SubTaskMeta{
	Name:             "extractCopilotOrgMetrics",
	EntryPoint:       ExtractCopilotOrgMetrics,
	EnabledByDefault: true,
	Description:      "Extract Copilot metrics into tool-layer tables",
}

// CollectCopilotOrgMetrics is a placeholder implementation that will be filled in future phases.
func CollectCopilotOrgMetrics(taskCtx plugin.SubTaskContext) errors.Error {
	return nil
}

// CollectCopilotSeatAssignments is a placeholder implementation that will be filled in future phases.
func CollectCopilotSeatAssignments(taskCtx plugin.SubTaskContext) errors.Error {
	return nil
}

// ExtractCopilotOrgMetrics is a placeholder implementation that will be filled in future phases.
func ExtractCopilotOrgMetrics(taskCtx plugin.SubTaskContext) errors.Error {
	return nil
}
