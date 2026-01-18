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

import pydevlake.domain_layer.devops as devops
from azuredevops.api import AzureDevOpsAPI
from azuredevops.models import Build
from azuredevops.models import GitRepository
from pydevlake import Context, DomainType, Stream


class Builds(Stream):
    tool_model = Build
    domain_types = [DomainType.CICD]
    domain_models = [devops.CiCDPipelineCommit, devops.CICDPipeline]

    def collect(self, state, context) -> Iterable[tuple[object, dict]]:
        repo: GitRepository = context.scope
        api = AzureDevOpsAPI(context.connection)
        provider = repo.provider or 'tfsgit'
        response = api.builds(repo.org_id, repo.project_id, repo.id, repo.provider or 'tfsgit')
        for raw_build in response:
            raw_build["x_request_url"] = response.get_url_with_query_string()
            raw_build["x_request_input"] = {
                "OrgId": repo.org_id,
                "ProjectId": repo.project_id,
                "RepoId": repo.id,
                "Provider": provider,
            }
            yield raw_build, state

    def convert(self, b: Build, ctx: Context):
        if not b.start_time:
            return

        result = devops.CICDResult.RESULT_DEFAULT
        if b.result == Build.BuildResult.Canceled:
            result = devops.CICDResult.FAILURE
        elif b.result == Build.BuildResult.Failed:
            result = devops.CICDResult.FAILURE
        elif b.result == Build.BuildResult.PartiallySucceeded:
            result = devops.CICDResult.SUCCESS
        elif b.result == Build.BuildResult.Succeeded:
            result = devops.CICDResult.SUCCESS

        status = devops.CICDStatus.STATUS_OTHER
        if b.status == Build.BuildStatus.Cancelling:
            status = devops.CICDStatus.DONE
        elif b.status == Build.BuildStatus.Completed:
            status = devops.CICDStatus.DONE
        elif b.status == Build.BuildStatus.InProgress:
            status = devops.CICDStatus.IN_PROGRESS
        elif b.status == Build.BuildStatus.NotStarted:
            status = devops.CICDStatus.IN_PROGRESS
        elif b.status == Build.BuildStatus.Postponed:
            status = devops.CICDStatus.IN_PROGRESS

        type = devops.CICDType.BUILD
        if ctx.scope_config.deployment_pattern and ctx.scope_config.deployment_pattern.search(b.name):
            type = devops.CICDType.DEPLOYMENT

        # Determine if this is a production environment
        # Match production_pattern against pipeline name
        environment = None
        if ctx.scope_config.production_pattern is not None:
            if ctx.scope_config.production_pattern.search(b.name):
                environment = devops.CICDEnvironment.PRODUCTION
        else:
            # No production_pattern configured - default to PRODUCTION for deployments
            if type == devops.CICDType.DEPLOYMENT:
                environment = devops.CICDEnvironment.PRODUCTION

        if b.finish_time:
            duration_sec = abs(b.finish_time.timestamp() - b.start_time.timestamp())
        else:
            duration_sec = float(0.0)

        yield devops.CICDPipeline(
            name=b.name,
            status=status,
            result=result,
            original_status=str(b.status),
            original_result=str(b.result),
            created_date=b.queue_time,
            queued_date=b.queue_time,
            started_date=b.start_time,
            finished_date=b.finish_time,
            duration_sec=duration_sec,
            environment=environment,
            type=type,
            cicd_scope_id=ctx.scope.domain_id(),
            display_title=b.display_title,
            url=b.url,
        )

        if b.source_version is not None:
            yield devops.CiCDPipelineCommit(
                pipeline_id=b.domain_id(),
                commit_sha=b.source_version,
                branch=b.source_branch,
                repo_id=ctx.scope.domain_id(),
                repo_url=ctx.scope.url,
                display_title=b.display_title,
                url=b.url,
            )
