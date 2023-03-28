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
        response = azuredevops_api.git_repo_pull_request_commits(repo.org_id, repo.project_id, parent.repo_id, parent.id)
        for raw_commit in response:
            raw_commit["pull_request_id"] = parent.domain_id()
            yield raw_commit, state

    def extract(self, raw_data: dict) -> GitPullRequestCommit:
        return GitPullRequestCommit(
            **raw_data,
            commit_sha = raw_data["commitId"],
            author_name = raw_data["author"]["name"],
            author_email = raw_data["author"]["email"],
            authored_date = raw_data["author"]["date"],
            committer_name = raw_data["committer"]["name"],
            committer_email = raw_data["committer"]["email"],
            commit_date = raw_data["committer"]["date"],
            additions = raw_data["changeCounts"]["Add"] if "changeCounts" in raw_data else 0,
            deletions = raw_data["changeCounts"]["Delete"] if "changeCounts" in raw_data else 0
        )

    def convert(self, commit: GitPullRequestCommit, context) -> Iterable[code.PullRequestCommit]:
        yield code.PullRequestCommit(
            commit_sha=commit.commit_sha,
            pull_request_id=commit.pull_request_id,
        )
