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

from typing import Iterable

from azuredevops.api import AzureDevOpsAPI
from azuredevops.models import GitRepository, GitPullRequest
from pydevlake import Stream, DomainType, domain_id
import pydevlake.domain_layer.code as code


class GitPullRequests(Stream):
    tool_model = GitPullRequest
    domain_types = [DomainType.CODE]

    def collect(self, state, context) -> Iterable[tuple[object, dict]]:
        api = AzureDevOpsAPI(context.connection)
        repo: GitRepository = context.scope
        response = api.git_repo_pull_requests(repo.org_id, repo.project_id, repo.id)
        for raw_pr in response:
            yield raw_pr, state

    def convert(self, pr: GitPullRequest, ctx):
        repo_id = ctx.scope.domain_id()
        # If the PR is from a fork, we forge a new repo ID for the base repo but it doesn't correspond to a real repo
        base_repo_id = domain_id(GitRepository, ctx.connection.id, pr.fork_repo_id) if pr.fork_repo_id is not None else repo_id

        # Use the same status values as GitHub plugin
        status = None
        if pr.status == GitPullRequest.PRStatus.Abandoned:
            status = 'CLOSED'
        elif pr.status == GitPullRequest.PRStatus.Active:
            status = 'OPEN'
        elif pr.status == GitPullRequest.PRStatus.Completed:
            status = 'MERGED'

        yield code.PullRequest(
            base_repo_id=base_repo_id,
            head_repo_id=repo_id,
            status=status,
            title=pr.title,
            description=pr.description,
            url=pr.url,
            author_name=pr.created_by_name,
            author_id=pr.created_by_id,
            pull_request_key=pr.pull_request_id,
            created_date=pr.creation_date,
            merged_date=pr.closed_date,
            closed_date=pr.closed_date,
            type=pr.type,
            component="", # not supported
            merge_commit_sha=pr.merge_commit_sha,
            head_ref=pr.source_ref_name,
            base_ref=pr.target_ref_name,
            head_commit_sha=pr.source_commit_sha,
            base_commit_sha=pr.target_commit_sha
        )
