import typing
from typing import Iterable, Tuple

from azure.api import AzureDevOpsAPI
from azure.models import Project
from pydevlake import Stream
from pydevlake.model import DomainModel, ToolModel


# TODO probably not needed except for remote-scopes
class Projects(Stream):

    @property
    def tool_model(self) -> typing.Type[ToolModel]:
        return Project

    @property
    def domain_models(self) -> Iterable[typing.Type[DomainModel]]:
        return [DomainModel]

    def collect(self, state, context) -> Iterable[Tuple[object, dict]]:
        connection = context.connection
        options = context.options
        api = AzureDevOpsAPI(connection.base_url, connection.pat)
        response = api.projects(options['org'])
        for raw_project in api.parse_response(response):
            yield raw_project, state

    def convert(self, project: Project, context) -> Iterable[DomainModel]:
        # dummy return for now
        return None
