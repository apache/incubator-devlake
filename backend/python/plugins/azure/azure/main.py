from azure.api import AzureDevOpsAPI
from azure.models import AzureDevOpsConnection, GitRepository
from azure.streams.builds import Builds
from azure.streams.commits import GitCommits
from azure.streams.jobs import Jobs
from azure.streams.pull_request_commits import GitPullRequestCommits
from azure.streams.pull_requests import GitPullRequests
from azure.streams.repositories import GitRepositories

from pydevlake import Plugin, RemoteScopeGroup
from pydevlake.domain_layer.code import Repo
from pydevlake.domain_layer.devops import CicdScope


class AzureDevOpsPlugin(Plugin):

    @property
    def connection_type(self):
        return AzureDevOpsConnection

    @property
    def tool_scope_type(self):
        return GitRepository

    def domain_scopes(self, git_repo: GitRepository):
        yield Repo(
            name=git_repo.name,
            url=git_repo.url,
            forked_from=git_repo.parentRepositoryUrl,
            deleted=git_repo.isDisabled,
        )

        yield CicdScope(
            name=git_repo.name,
            description=git_repo.name,
            url=git_repo.url
        )

    def remote_scope_groups(self, ctx) -> list[RemoteScopeGroup]:
        api = AzureDevOpsAPI(ctx.connection.base_url, ctx.connection.pat)
        member_id = api.my_profile.json['id']
        accounts = api.accounts(member_id).json
        orgs = [acc['accountId'] for acc in accounts]
        for org in orgs:
            for proj in api.projects(org):
                yield RemoteScopeGroup(
                    id=f'{org}/{proj["name"]}',
                    name=proj['name']
                )

    def remote_scopes(self, ctx, group_id: str) -> list[GitRepository]:
        org, proj = group_id.split('/')
        api = AzureDevOpsAPI(ctx.connection.base_url, ctx.connection.pat)
        for raw_repo in api.git_repos(org, proj):
            repo = GitRepository(**raw_repo)
            if not repo.defaultBranch:
                return None
            repo.project_id = raw_repo['project']["id"]
            if "parentRepository" in raw_repo:
                repo.parentRepositoryUrl = raw_repo["parentRepository"]["url"]
            yield repo

    @property
    def name(self) -> str:
        return "azure"

    def test_connection(self, connection: AzureDevOpsConnection):
        resp = AzureDevOpsAPI(connection.base_url, connection.pat).projects(connection.org)
        if resp.status != 200:
            raise Exception(f"Invalid connection: {resp.json}")

    @property
    def streams(self):
        return [
            # Projects,
            GitRepositories,
            GitPullRequests,
            GitPullRequestCommits,
            GitCommits,
            Builds,
            Jobs,
        ]


if __name__ == '__main__':
    AzureDevOpsPlugin.start()
