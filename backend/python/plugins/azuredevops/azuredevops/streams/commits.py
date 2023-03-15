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

import iso8601 as iso8601

from azuredevops.api import AzureDevOpsAPI
from azuredevops.models import GitRepository, GitCommit
from pydevlake import Stream, DomainType, Context
from pydevlake.domain_layer.code import Commit as DomainCommit
from pydevlake.domain_layer.code import RepoCommit as DomainRepoCommit


class GitCommits(Stream):
    tool_model = GitCommit
    domain_types = [DomainType.CODE]

    def collect(self, state, context) -> Iterable[tuple[object, dict]]:
        connection = context.connection
        repo: GitRepository = context.scope
        azuredevops_api = AzureDevOpsAPI(connection.base_url, connection.pat)
        response = azuredevops_api.commits(repo.org_id, repo.project_id, repo.id)
        for raw_commit in response:
            raw_commit["repo_id"] = repo.id
            yield raw_commit, state

    def extract(self, raw_data: dict) -> GitCommit:
        return extract_raw_commit(raw_data)

    def convert(self, commit: GitCommit, ctx: Context) -> Iterable[DomainCommit]:
        yield DomainCommit(
            sha=commit.commit_sha,
            additions=commit.additions,
            deletions=commit.deletions,
            message=commit.comment,
            author_name=commit.author_name,
            author_email=commit.author_email,
            authored_date=commit.authored_date,
            author_id=commit.author_name,
            committer_name=commit.committer_name,
            committer_email=commit.committer_email,
            committed_date=commit.commit_date,
            committer_id=commit.committer_name,
        )

        yield DomainRepoCommit(
                repo_id=commit.repo_id,
                commit_sha=commit.commit_sha,
        )


def extract_raw_commit(stream: Stream, raw_data: dict, ctx: Context) -> GitCommit:
    commit: GitCommit = stream.tool_model(**raw_data)
    repo: GitRepository = ctx.scope
    commit.project_id = repo.project_id
    commit.repo_id = raw_data["repo_id"]
    commit.commit_sha = raw_data["commitId"]
    commit.author_name = raw_data["author"]["name"]
    commit.author_email = raw_data["author"]["email"]
    commit.authored_date = iso8601.parse_date(raw_data["author"]["date"])
    commit.committer_name = raw_data["committer"]["name"]
    commit.committer_email = raw_data["committer"]["email"]
    commit.commit_date = iso8601.parse_date(raw_data["committer"]["date"])
    if "changeCounts" in raw_data:
        commit.additions = raw_data["changeCounts"]["Add"]
        commit.deletions = raw_data["changeCounts"]["Delete"]
    return commit
