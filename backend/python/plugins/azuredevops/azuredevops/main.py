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

from urllib.parse import urlparse

from azuredevops.api import AzureDevOpsAPI
from azuredevops.models import AzureDevOpsConnection, GitRepository
from azuredevops.streams.builds import Builds
from azuredevops.streams.commits import GitCommits
from azuredevops.streams.jobs import Jobs
from azuredevops.streams.pull_request_commits import GitPullRequestCommits
from azuredevops.streams.pull_requests import GitPullRequests

from pydevlake import Plugin, RemoteScopeGroup, DomainType
from pydevlake.domain_layer.code import Repo
from pydevlake.domain_layer.devops import CicdScope
from pydevlake.pipeline_tasks import gitextractor, refdiff


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

    def remote_scope_groups(self, connection) -> list[RemoteScopeGroup]:
        api = AzureDevOpsAPI(connection)
        member_id = api.my_profile().json['id']
        accounts = api.accounts(member_id).json
        for account in accounts['value']:
            org = account['accountName']
            for proj in api.projects(org):
                proj_name = proj['name']

                yield RemoteScopeGroup(
                    id=f'{org}/{proj_name}',
                    name=proj_name
                )

    def remote_scopes(self, connection, group_id: str) -> list[GitRepository]:
        org, proj = group_id.split('/')
        api = AzureDevOpsAPI(connection)
        for raw_repo in api.git_repos(org, proj):
            url = urlparse(raw_repo['remoteUrl'])
            url = url._replace(netloc=f'{url.username}:{connection.pat}@{url.hostname}')
            repo = GitRepository(**raw_repo, project_id=proj, org_id=org, url=url.geturl())
            if not repo.defaultBranch:
                return None
            if "parentRepository" in raw_repo:
                repo.parentRepositoryUrl = raw_repo["parentRepository"]["url"]
            yield repo

    def test_connection(self, connection: AzureDevOpsConnection):
        resp = AzureDevOpsAPI(connection).my_profile()
        if resp.status != 200:
            raise Exception(f"Invalid token: {connection.token}")

    def extra_tasks(self, scope: GitRepository, entity_types: list[str], connection_id: int):
        if DomainType.CODE in entity_types:
            # TODO: pass proxy
            return [gitextractor(scope.url, scope.id, None)]
        else:
            return []

    def extra_stages(self, scopes: list[GitRepository], entity_types: list[str], connection_id: int):
        if DomainType.CODE in entity_types:
            # TODO: pass refdiff options
            return [refdiff(scope.id, {}) for scope in scopes]
        else:
            return []

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
