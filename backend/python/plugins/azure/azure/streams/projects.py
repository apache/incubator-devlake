from typing import Iterable, Tuple

from azure.api import AzureDevOpsAPI
from azure.models import Project
from pydevlake import Stream, DomainType


# TODO probably not needed except for remote-scopes
class Projects(Stream):
    tool_model = Project
    domain_types = [DomainType.CROSS] # ??

    def collect(self, state, context) -> Iterable[Tuple[object, dict]]:
        connection = context.connection
        options = context.options
        api = AzureDevOpsAPI(connection.base_url, connection.pat)
        response = api.projects(options['org'])
        for raw_project in api.parse_response(response):
            yield raw_project, state

    def convert(self, project: Project, context) -> Iterable[object]:
        # dummy return for now
        return None
