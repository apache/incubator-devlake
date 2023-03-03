from typing import Iterable, Tuple

from azure.api import AzureDevOpsAPI
from azure.models import GitRepository
from pydevlake import Stream, DomainType
from pydevlake.domain_layer.code import Repo as DomainRepo


class GitRepositories(Stream):
    tool_model = GitRepository
    domain_types = [DomainType.CODE]

    def collect(self, state, context) -> Iterable[Tuple[object, dict]]:
        connection = context.connection
        options = context.options
        api = AzureDevOpsAPI(connection.base_url, connection.pat)
        response = api.git_repos(options['org'], options['project'])
        for raw_repo in api.parse_response(response):
            yield raw_repo, state

    def extract(self, raw_data: dict, context) -> GitRepository:
        repo: GitRepository = self.tool_model(**raw_data)
        if not repo.defaultBranch:
            return None
        project: dict = raw_data['project']
        repo.project_id = project["id"]
        if "parentRepository" in raw_data:
            repo.parentRepositoryUrl = raw_data["parentRepository"]["url"]
        return repo

    def convert(self, repo: GitRepository, context) -> Iterable[DomainRepo]:
        # dummy return for now
        yield DomainRepo(
            name=repo.name,
            url=repo.url,
            forked_from=repo.parentRepositoryUrl,
            deleted=repo.isDisabled,
        )