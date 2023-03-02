import typing
from datetime import datetime
from typing import Iterable

import iso8601 as iso8601

from azure.api import AzureDevOpsAPI
from azure.helper import db
from azure.models import GitRepository, GitPullRequest, GitCommit
from azure.streams.repositories import GitRepositories
from pydevlake import Substream, Stream, Context
from pydevlake.domain_layer.code import PullRequest as DomainPullRequest
from pydevlake.model import DomainModel, ToolModel


class GitPullRequests(Substream):

    @property
    def tool_model(self) -> typing.Type[ToolModel]:
        # TODO define pr model
        return GitPullRequest

    @property
    def domain_models(self) -> Iterable[typing.Type[DomainModel]]:
        return [DomainPullRequest]

    @property
    def parent_stream(self) -> Stream:
        return GitRepositories(self.plugin_name)

    def collect(self, state, context, parent: GitRepository) -> Iterable[tuple[object, dict]]:
        connection = context.connection
        options = context.options
        azure_api = AzureDevOpsAPI(connection.base_url, connection.pat)
        # grab this info off the parent results
        response = azure_api.git_repo_pull_requests(options["org"], options["project"], parent.name)
        for raw_pr in azure_api.parse_response(response):
            yield raw_pr, state

    def extract(self, raw_data: dict, context) -> ToolModel:
        pr: GitPullRequest = self.tool_model(**raw_data)
        pr.id = raw_data["pullRequestId"]
        pr.project_id = context.options["project"]
        pr.created_by_id = raw_data["createdBy"]["id"]
        pr.created_by_name = raw_data["createdBy"]["displayName"]
        if "closedDate" in raw_data:
            pr.closed_date = iso8601.parse_date(raw_data["closedDate"])
        pr.creation_date = iso8601.parse_date(raw_data["creationDate"])
        pr.code_review_id = raw_data["codeReviewId"]
        pr.repo_id = raw_data["repository"]["id"]
        pr.title = raw_data["title"]
        pr.description = raw_data["description"]
        pr.source_commit_sha = raw_data["lastMergeSourceCommit"]["commitId"]
        pr.target_commit_sha = raw_data["lastMergeTargetCommit"]["commitId"]
        pr.merge_commit_sha = raw_data["lastMergeCommit"]["commitId"]
        pr.source_ref_name = raw_data["sourceRefName"]
        pr.target_ref_name = raw_data["targetRefName"]
        pr.status = raw_data["status"]
        pr.url = raw_data["url"]
        if "labels" in raw_data:
            # TODO get this off transformation rules regex
            pr.type = raw_data["labels"][0]["name"]
        if "forkSource" in raw_data:
            pr.fork_repo_id = raw_data["forkSource"]["repository"]["id"]
        return pr

    def convert(self, pr: GitPullRequest, context: Context) -> Iterable[DomainPullRequest]:
        merged_date: datetime = None
        if pr.status == GitPullRequest.Status.Completed:
            # query from commits
            merge_commit: GitCommit = db.get(context, GitCommit, GitCommit.commit_sha == pr.merge_commit_sha)
            merged_date = merge_commit.commit_date
        yield DomainPullRequest(
            base_repo_id=(pr.fork_repo_id if pr.fork_repo_id is not None else pr.repo_id),
            head_repo_id=pr.repo_id,
            status=pr.status.value,
            title=pr.title,
            description=pr.description,
            url=pr.url,
            author_name=pr.created_by_name,
            author_id=pr.created_by_id,
            pull_request_key=pr.id,
            created_date=pr.creation_date,
            merged_date=merged_date,
            closed_date=pr.closed_date,
            type=pr.type,
            component="", # not supported
            merge_commit_sha=pr.merge_commit_sha,
            head_ref=pr.target_ref_name,
            base_ref=pr.source_ref_name,
            base_commit_sha=pr.source_commit_sha,
            head_commit_sha=pr.target_commit_sha,
        )
