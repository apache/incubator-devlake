from typing import Iterable

from azure.api import AzureDevOpsAPI
from azure.models import GitPullRequest, GitCommit
from azure.streams.commits import extract_raw_commit
from azure.streams.pull_requests import GitPullRequests
from pydevlake import Substream, DomainType
from pydevlake.domain_layer.code import PullRequestCommit as DomainPullRequestCommit



class GitPullRequestCommits(Substream):
    tool_model = GitCommit
    domain_types = [DomainType.CODE]
    parent_stream = GitPullRequests

    def collect(self, state, context, parent: GitPullRequest) -> Iterable[tuple[object, dict]]:
        connection = context.connection
        options = context.options
        azure_api = AzureDevOpsAPI(connection.base_url, connection.pat)
        # grab this info off the parent results
        response = azure_api.git_repo_pull_request_commits(options["org"], options["project"], parent.repo_id, parent.id)
        for raw_commit in response:
            raw_commit["repo_id"] = parent.repo_id
            yield raw_commit, state

    def extract(self, raw_data: dict) -> GitCommit:
        return extract_raw_commit(self, raw_data)

    def convert(self, commit: GitCommit, context) -> Iterable[DomainPullRequestCommit]:
        yield DomainPullRequestCommit(
            commit_sha=commit.commit_sha,
            pull_request_id=commit.repo_id,
        )
