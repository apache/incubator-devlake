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

from http import HTTPStatus

from azuredevops.api import AzureDevOpsAPI
from azuredevops.models import Job, Build, GitRepository
from azuredevops.streams.builds import Builds
from pydevlake import Context, Substream, DomainType
from pydevlake.api import APIException
import pydevlake.domain_layer.devops as devops


class Jobs(Substream):
    tool_model = Job
    domain_types = [DomainType.CICD]
    parent_stream = Builds

    def collect(self, state, context, parent: Build) -> Iterable[tuple[object, dict]]:
        repo: GitRepository = context.scope
        api = AzureDevOpsAPI(context.connection)
        try:
            response = api.jobs(repo.org_id, repo.project_id, parent.id)
        except APIException as e:
            # Asking for the timeline of a deleted build returns a 204.
            # But a "deleted" build may be "cleaned" (i.e. deleted for real)
            # after some time. In this case, the timeline endpoint returns a
            # 404 instead.
            if e.response.status == HTTPStatus.NOT_FOUND:
                return
            raise
        # If a build has failed before any jobs have started, e.g. due to a
        # bad YAML file, then the jobs endpoint will return a 204 NO CONTENT.
        if response.status == HTTPStatus.NO_CONTENT:
            return
        for raw_job in response.json["records"]:
            if raw_job["type"] == "Job":
                raw_job["build_id"] = parent.domain_id()
                yield raw_job, state

    def convert(self, j: Job, ctx: Context) -> Iterable[devops.CICDPipeline]:
        if not j.startTime:
            return

        result = None
        if j.result == Job.JobResult.Abandoned:
            result = devops.CICDResult.ABORT
        elif j.result == Job.JobResult.Canceled:
            result = devops.CICDResult.ABORT
        elif j.result == Job.JobResult.Failed:
            result = devops.CICDResult.FAILURE
        elif j.result == Job.JobResult.Skipped:
            result = devops.CICDResult.ABORT
        elif j.result == Job.JobResult.Succeeded:
            result = devops.CICDResult.SUCCESS
        elif j.result == Job.JobResult.SucceededWithIssues:
            result = devops.CICDResult.FAILURE

        status = None
        if j.state == Job.JobState.Completed:
            status = devops.CICDStatus.DONE
        elif j.state == Job.JobState.InProgress:
            status = devops.CICDStatus.IN_PROGRESS
        if j.state == Job.JobState.Pending:
            status = devops.CICDStatus.IN_PROGRESS

        type = devops.CICDType.BUILD
        if ctx.transformation_rule and ctx.transformation_rule.deployment_pattern.search(j.name):
            type = devops.CICDType.DEPLOYMENT
        environment = devops.CICDEnvironment.TESTING
        if ctx.transformation_rule and ctx.transformation_rule.production_pattern.search(j.name):
            environment = devops.CICDEnvironment.PRODUCTION

        if j.finishTime:
            duration_sec = abs(j.finishTime.second-j.startTime.second)
        else:
            duration_sec = 0

        yield devops.CICDTask(
            id=j.id,
            name=j.name,
            pipeline_id=j.build_id,
            status=status,
            created_date=j.startTime,
            finished_date=j.finishTime,
            result=result,
            type=type,
            duration_sec=duration_sec,
            environment=environment,
            cicd_scope_id=ctx.scope.domain_id()
        )
