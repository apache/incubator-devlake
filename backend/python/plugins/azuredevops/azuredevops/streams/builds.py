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
from pydevlake import Context, DomainType, Stream, logger
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

    def extract(self, raw_data: dict) -> Build:
        build: Build = self.tool_model(**raw_data)
        build.project_id = raw_data["project"]["id"]
        build.repo_id = raw_data["repository"]["id"]
        build.repo_type = raw_data["repository"]["type"]
        build.build_number = raw_data["buildNumber"]
        build.tags = ",".join(raw_data["tags"])
        build.build_result = Build.Result(raw_data["result"])
        trigger_info: dict = raw_data["triggerInfo"]
        if "ci.sourceSha" in trigger_info: # this key is not guaranteed to be in here per docs
            assert build.source_version == trigger_info["ci.sourceSha"]
        return build

    def convert(self, b: Build, ctx: Context):
        result = None
        match b.build_result:
            case Build.Result.Canceled:
                result = devops.CICDResult.ABORT
            case Build.Result.Failed:
                result = devops.CICDResult.FAILURE
            case Build.Result.PartiallySucceeded:
                result = devops.CICDResult.SUCCESS
            case Build.Result.Succeeded:
                result = devops.CICDResult.SUCCESS

        status = None
        match b.status:
            case Build.Status.All:
                status = devops.CICDStatus.IN_PROGRESS
            case Build.Status.Cancelling:
                status = devops.CICDStatus.DONE
            case Build.Status.Completed:
                status = devops.CICDStatus.DONE
            case Build.Status.InProgress:
                status = devops.CICDStatus.IN_PROGRESS
            case Build.Status.NotStarted:
                status = devops.CICDStatus.IN_PROGRESS
            case Build.Status.Postponed:
                status = devops.CICDStatus.IN_PROGRESS

        yield devops.CICDPipeline(
            name=b.id,
            status=status,
            created_date=b.start_time,
            finished_date=b.finish_time,
            result=result,
            duration_sec=abs(b.finish_time.second-b.start_time.second),
            environment=devops.CICDEnvironment.PRODUCTION,
            type=devops.CICDType.DEPLOYMENT,
            cicd_scope_id=b.repo_id,
        )

        repo_url = None
        if b.repo_type == 'GitHub':
            repo_url = f'https://github.com/{b.repo_id}'

        yield devops.CiCDPipelineCommit(
            pipeline_id=b.id,
            commit_sha=b.source_version,
            branch=b.source_branch,
            repo_id=b.repo_id,
            repo=repo_url,
        )
