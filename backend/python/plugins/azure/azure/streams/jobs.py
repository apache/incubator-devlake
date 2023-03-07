from typing import Iterable

from azure.api import AzureDevOpsAPI
from azure.models import AzureDevOpsConnection, Job, Build
from azure.streams.builds import Builds
from pydevlake import Context, Substream, DomainType
from pydevlake.domain_layer.devops import *


class Jobs(Substream):
    tool_model = Job
    domain_types = [DomainType.CICD]
    parent_stream = Builds

    def collect(self, state, context, parent: Build) -> Iterable[tuple[object, dict]]:
        connection: AzureDevOpsConnection = context.connection
        options = context.options
        azure_api = AzureDevOpsAPI(connection.base_url, connection.pat)
        # grab this info off the parent results
        response = azure_api.jobs(options["org"], options["project"], parent.id)
        if response.status != 200:
            yield None, state
        else:
            for raw_job in response.json["records"]:
                raw_job["build_id"] = parent.id
                raw_job["repo_id"] = parent.repo_id
                yield raw_job, state

    def extract(self, raw_data: dict) -> Job:
        job: Job = self.tool_model(**raw_data)
        if job.type != job.type.Job:
            return None
        return job

    def convert(self, j: Job, ctx: Context) -> Iterable[CICDPipeline]:
        yield CICDTask(
            name=j.id,
            pipeline_id=j.build_id,
            status=j.state.value,
            created_date=j.startTime,
            finished_date=j.finishTime,
            result=j.result.value,
            type=CICDType.DEPLOYMENT.value,
            duration_sec=abs(j.finishTime.second-j.startTime.second),
            environment=CICDEnvironment.PRODUCTION.value,
            cicd_scope_id=j.repo_id
        )
