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
from azuredevops.models import GitPullRequest, GitPullRequestCommit, GitRepository
from azuredevops.streams.pull_requests import GitPullRequests
from pydevlake import Substream, DomainType
import pydevlake.domain_layer.code as code


class GitPullRequestCommits(Substream):
    tool_model = GitPullRequestCommit
    domain_types = [DomainType.CODE]
    parent_stream = GitPullRequests

    def collect(self, state, context, parent: GitPullRequest) -> Iterable[tuple[object, dict]]:
        repo: GitRepository = context.scope
        azuredevops_api = AzureDevOpsAPI(context.connection)
        response = azuredevops_api.git_repo_pull_request_commits(repo.org_id, repo.project_id, repo.id, parent.pull_request_id)
        for raw_commit in response:
            raw_commit["pull_request_id"] = parent.domain_id()
            yield raw_commit, state

    def convert(self, commit: GitPullRequestCommit, context) -> Iterable[code.PullRequestCommit]:
        yield code.PullRequestCommit(
            commit_sha=commit.commit_id,
            pull_request_id=commit.pull_request_id,
            commit_author_name=commit.author_name,
            commit_author_email=commit.author_email,
            commit_authored_date=commit.author_date
        )
