from datetime import datetime
from typing import Iterable

import iso8601 as iso8601

from azure.api import AzureDevOpsAPI
from azure.helper import db
from azure.models import GitRepository, GitPullRequest, GitCommit
from pydevlake import Stream, Context, DomainType
from pydevlake.domain_layer.code import PullRequest as DomainPullRequest


class GitPullRequests(Stream):
    tool_model = GitPullRequest
    domain_types = [DomainType.CODE]

    def collect(self, state, context) -> Iterable[tuple[object, dict]]:
        connection = context.connection
        azure_api = AzureDevOpsAPI(connection.base_url, connection.pat)
        repo: GitRepository = context.scope
        response = azure_api.git_repo_pull_requests(repo.org_id, repo.project_id, repo.id)
        for raw_pr in response:
            yield raw_pr, state

    def extract(self, raw_data: dict) -> GitPullRequest:
        pr = GitPullRequest(**raw_data)
        pr.id = raw_data["pullRequestId"]
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
