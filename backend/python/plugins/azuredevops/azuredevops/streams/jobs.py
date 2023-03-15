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

from azuredevops.api import AzureDevOpsAPI
from azuredevops.models import AzureDevOpsConnection, Job, Build, GitRepository
from azuredevops.streams.builds import Builds
from pydevlake import Context, Substream, DomainType
import pydevlake.domain_layer.devops as devops


class Jobs(Substream):
    tool_model = Job
    domain_types = [DomainType.CICD]
    parent_stream = Builds

    def collect(self, state, context, parent: Build) -> Iterable[tuple[object, dict]]:
        connection: AzureDevOpsConnection = context.connection
        repo: GitRepository = context.scope
        azuredevops_api = AzureDevOpsAPI(connection.pat)
        response = azuredevops_api.jobs(repo.org_id, repo.project_id, parent.id)
        if response.status != 200:
            yield None, state
        else:
            for raw_job in response.json["records"]:
                raw_job["build_id"] = parent.id
                raw_job["repo_id"] = parent.repo_id
                yield raw_job, state


    def convert(self, j: Job, ctx: Context) -> Iterable[devops.CICDPipeline]:
        result = None
        match j.result:
            case Job.Result.Abandoned:
                result = devops.CICDResult.ABORT
            case Job.Result.Canceled:
                result = devops.CICDResult.ABORT
            case Job.Result.Failed:
                result = devops.CICDResult.FAILURE
            case Job.Result.Skipped:
                result = devops.CICDResult.ABORT
            case Job.Result.Succeeded:
                result = devops.CICDResult.SUCCESS
            case Job.Result.SucceededWithIssues:
                result = devops.CICDResult.FAILURE

        status = None
        match j.state:
            case Job.State.Completed:
                status = devops.CICDStatus.DONE
            case Job.State.InProgress:
                status = devops.CICDStatus.IN_PROGRESS
            case Job.State.Pending:
                status = devops.CICDStatus.IN_PROGRESS

        yield devops.CICDTask(
            id=j.id,
            name=j.name,
            pipeline_id=j.build_id,
            status=status,
            created_date=j.startTime,
            finished_date=j.finishTime,
            result=result,
            type=devops.CICDType.BUILD,
            duration_sec=abs(j.finishTime.second-j.startTime.second),
            environment=devops.CICDEnvironment.PRODUCTION,
            cicd_scope_id=j.repo_id
        )
