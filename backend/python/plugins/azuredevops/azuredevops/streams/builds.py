# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from typing import Iterable

import iso8601 as iso8601

from azuredevops.api import AzureDevOpsAPI
from azuredevops.models import GitRepository
from azuredevops.models import Build
from pydevlake import Context, DomainType, Stream
import pydevlake.domain_layer.devops as devops


class Builds(Stream):
    tool_model = Build
    domain_types = [DomainType.CICD]

    def collect(self, state, context) -> Iterable[tuple[object, dict]]:
        repo: GitRepository = context.scope
        api = AzureDevOpsAPI(context.connection)
        response = api.builds(repo.org_id, repo.project_id, repo.id, 'tfsgit')
        for raw_build in response:
            yield raw_build, state

    def convert(self, b: Build, ctx: Context):
        result = None
        if b.result == Build.BuildResult.Canceled:
            result = devops.CICDResult.ABORT
        elif b.result == Build.BuildResult.Failed:
            result = devops.CICDResult.FAILURE
        elif b.result == Build.BuildResult.PartiallySucceeded:
            result = devops.CICDResult.SUCCESS
        elif b.result ==  Build.BuildResult.Succeeded:
            result = devops.CICDResult.SUCCESS

        status = None
        if b.status == Build.BuildStatus.Cancelling:
            status = devops.CICDStatus.DONE
        elif b.status == Build.BuildStatus.Completed:
            status = devops.CICDStatus.DONE
        elif b.status ==  Build.BuildStatus.InProgress:
            status = devops.CICDStatus.IN_PROGRESS
        elif b.status == Build.BuildStatus.NotStarted:
            status = devops.CICDStatus.IN_PROGRESS
        elif b.status ==  Build.BuildStatus.Postponed:
            status = devops.CICDStatus.IN_PROGRESS

        type = devops.CICDType.BUILD
        if ctx.transformation_rule and ctx.transformation_rule.deployment_pattern.search(b.name):
            type = devops.CICDType.DEPLOYMENT
        environment = devops.CICDEnvironment.TESTING
        if ctx.transformation_rule and ctx.transformation_rule.production_pattern.search(b.name):
            environment = devops.CICDEnvironment.PRODUCTION

        yield devops.CICDPipeline(
            name=b.name,
            status=status,
            created_date=b.start_time,
            finished_date=b.finish_time,
            result=result,
            duration_sec=abs(b.finish_time.second-b.start_time.second),
            environment=environment,
            type=type,
            cicd_scope_id=ctx.scope.domain_id(),
        )

        yield devops.CiCDPipelineCommit(
            pipeline_id=b.domain_id(),
            commit_sha=b.source_version,
            branch=b.source_branch,
            repo_id=ctx.scope.domain_id(),
            repo_url=ctx.scope.url,
        )
