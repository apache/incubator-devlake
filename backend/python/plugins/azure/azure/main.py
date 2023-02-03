from abc import ABC

from azure.api import AzureDevOpsAPI
from azure.models import AzureDevOpsConnection
from azure.streams.builds import Builds
from azure.streams.commits import GitCommits
from azure.streams.jobs import Jobs
from azure.streams.pull_request_commits import GitPullRequestCommits
from azure.streams.pull_requests import GitPullRequests
from azure.streams.repositories import GitRepositories
from pydevlake.plugin import Plugin


class AzureDevOpsPlugin(Plugin, ABC):

    @property
    def connection_type(self):
        return AzureDevOpsConnection

    @property
    def name(self) -> str:
        return "azure"

    def test_connection(self, connection: AzureDevOpsConnection):
        resp = AzureDevOpsAPI(connection.base_url, connection.pat).projects(connection.org)
        print(resp)
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
