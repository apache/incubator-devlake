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

from azuredevops.api import AzureDevOpsAPI
from azuredevops.models import AzureDevOpsConnection, GitRepository
from azuredevops.streams.builds import Builds
from azuredevops.streams.commits import GitCommits
from azuredevops.streams.jobs import Jobs
from azuredevops.streams.pull_request_commits import GitPullRequestCommits
from azuredevops.streams.pull_requests import GitPullRequests

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
        api = AzureDevOpsAPI(ctx.connection)
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
        api = AzureDevOpsAPI(ctx.connection)
        for raw_repo in api.git_repos(org, proj):
            repo = GitRepository(**raw_repo, project_id=proj, org_id=org)
            if not repo.defaultBranch:
                return None
            if "parentRepository" in raw_repo:
                repo.parentRepositoryUrl = raw_repo["parentRepository"]["url"]
            yield repo

    def test_connection(self, connection: AzureDevOpsConnection):
        resp = AzureDevOpsAPI(connection).my_profile()
        if resp.status != 200:
            raise Exception(f"Invalid token: {connection.token}")

    @property
    def streams(self):
        return [
            GitPullRequests,
            GitPullRequestCommits,
            GitCommits,
            Builds,
            Jobs,
        ]


if __name__ == '__main__':
    AzureDevOpsPlugin.start()
