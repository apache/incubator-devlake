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
from azuredevops.models import GitPullRequest, GitCommit, GitRepository
from azuredevops.streams.commits import extract_raw_commit
from azuredevops.streams.pull_requests import GitPullRequests
from pydevlake import Substream, DomainType
from pydevlake.domain_layer.code import PullRequestCommit as DomainPullRequestCommit


class GitPullRequestCommits(Substream):
    tool_model = GitCommit
    domain_types = [DomainType.CODE]
    parent_stream = GitPullRequests

    def collect(self, state, context, parent: GitPullRequest) -> Iterable[tuple[object, dict]]:
        repo: GitRepository = context.scope
        azuredevops_api = AzureDevOpsAPI(context.connection)
        response = azuredevops_api.git_repo_pull_request_commits(repo.org_id, repo.project_id, parent.repo_id, parent.id)
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
