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

from http import HTTPStatus
from typing import Iterable, Optional

import pydevlake.domain_layer.devops as devops
from azuredevops.api import AzureDevOpsAPI
from azuredevops.models import Job, Build, GitRepository
from azuredevops.streams.builds import Builds
from pydevlake import Context, Substream, DomainType
from pydevlake.api import APIException


def extract_environment_name(name: str, identifier: Optional[str], context: Context) -> Optional[str]:
    """
    Extract environment name from job/stage name or identifier using environment_pattern.

    The environment_pattern should contain a capture group to extract the environment name.
    For example: r'(?:deploy|predeploy)[_-](.+?)(?:[_-](?:helm|terraform))?$'
    This would extract 'xxxx-prod' from 'deploy_xxxx-prod_helm'
    """
    if not context.scope_config.environment_pattern:
        return None

    # Try to match against the name first
    match = context.scope_config.environment_pattern.search(name)
    if match and match.groups():
        return match.group(1)

    # If no match on name and identifier is available, try identifier
    if identifier:
        match = context.scope_config.environment_pattern.search(identifier)
        if match and match.groups():
            return match.group(1)

    return None


class Jobs(Substream):
    tool_model = Job
    domain_types = [DomainType.CICD]
    parent_stream = Builds
    domain_models = [devops.CICDTask]

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
            # Collect both Job and Stage records to support environment detection from stages
            if raw_job["type"] in ("Job", "Stage"):
                raw_job["build_id"] = parent.domain_id()
                raw_job["x_request_url"] = response.get_url_with_query_string()
                raw_job["x_request_input"] = {
                    "OrgId": repo.org_id,
                    "ProjectId": repo.project_id,
                    "BuildId": parent.id,
                }
                yield raw_job, state

    def convert(self, j: Job, ctx: Context) -> Iterable[devops.CICDPipeline]:
        if not j.start_time:
            return

        result = devops.CICDResult.RESULT_DEFAULT
        if j.result == Job.JobResult.Abandoned:
            result = devops.CICDResult.RESULT_DEFAULT
        elif j.result == Job.JobResult.Canceled:
            result = devops.CICDResult.FAILURE
        elif j.result == Job.JobResult.Failed:
            result = devops.CICDResult.FAILURE
        elif j.result == Job.JobResult.Skipped:
            result = devops.CICDResult.RESULT_DEFAULT
        elif j.result == Job.JobResult.Succeeded:
            result = devops.CICDResult.SUCCESS
        elif j.result == Job.JobResult.SucceededWithIssues:
            result = devops.CICDResult.FAILURE

        status = devops.CICDStatus.STATUS_OTHER
        if j.state == Job.JobState.Completed:
            status = devops.CICDStatus.DONE
        elif j.state == Job.JobState.InProgress:
            status = devops.CICDStatus.IN_PROGRESS
        if j.state == Job.JobState.Pending:
            status = devops.CICDStatus.IN_PROGRESS

        type = devops.CICDType.BUILD
        if ctx.scope_config.deployment_pattern and ctx.scope_config.deployment_pattern.search(j.name):
            type = devops.CICDType.DEPLOYMENT

        # Extract environment name using the new environment_pattern if configured
        extracted_env_name = extract_environment_name(j.name, j.identifier, ctx)

        # Determine if this is a production environment
        # Priority: 1) Use extracted environment name with production_pattern
        #           2) Fall back to matching production_pattern against job name
        environment = None
        if ctx.scope_config.production_pattern is not None:
            # If we extracted an environment name, use it for production matching
            if extracted_env_name:
                if ctx.scope_config.production_pattern.search(extracted_env_name):
                    environment = devops.CICDEnvironment.PRODUCTION
            # Fall back to matching against job name
            elif ctx.scope_config.production_pattern.search(j.name):
                environment = devops.CICDEnvironment.PRODUCTION
        else:
            # No production_pattern configured - default to PRODUCTION for deployments
            if type == devops.CICDType.DEPLOYMENT:
                environment = devops.CICDEnvironment.PRODUCTION

        if j.finish_time:
            duration_sec = abs(j.finish_time.timestamp() - j.start_time.timestamp())
        else:
            duration_sec = float(0.0)

        yield devops.CICDTask(
            id=j.id,
            name=j.name,
            pipeline_id=j.build_id,
            status=status,
            original_status=str(j.state),
            original_result=str(j.result),
            created_date=j.start_time,
            started_date=j.start_time,
            finished_date=j.finish_time,
            result=result,
            type=type,
            duration_sec=duration_sec,
            environment=environment,
            cicd_scope_id=ctx.scope.domain_id()
        )
