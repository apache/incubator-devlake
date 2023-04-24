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
from azuredevops.models import AzureDevOpsConnection, GitRepository, AzureDevOpsTransformationRule
from azuredevops.streams.builds import Builds
from azuredevops.streams.jobs import Jobs
from azuredevops.streams.pull_request_commits import GitPullRequestCommits
from azuredevops.streams.pull_requests import GitPullRequests

from pydevlake import Plugin, RemoteScopeGroup, DomainType, ScopeTxRulePair
from pydevlake.domain_layer.code import Repo
from pydevlake.domain_layer.devops import CicdScope
from pydevlake.pipeline_tasks import gitextractor, refdiff
from pydevlake.api import APIException


class AzureDevOpsPlugin(Plugin):

    @property
    def connection_type(self):
        return AzureDevOpsConnection

    @property
    def tool_scope_type(self):
        return GitRepository

    @property
    def transformation_rule_type(self):
        return AzureDevOpsTransformationRule

    def domain_scopes(self, git_repo: GitRepository):
        yield Repo(
            name=git_repo.name,
            url=git_repo.url,
            forked_from=git_repo.parent_repository_url
        )

        yield CicdScope(
            name=git_repo.name,
            description=git_repo.name,
            url=git_repo.url
        )

    def remote_scope_groups(self, connection) -> list[RemoteScopeGroup]:
        api = AzureDevOpsAPI(connection)
        if connection.organization:
            orgs = [connection.organization]
        else:
            member_id = api.my_profile().json['id']
            accounts = api.accounts(member_id).json
            orgs = [account['accountName'] for account in accounts['value']]

        for org in orgs:
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
            raw_repo['project_id'] = proj
            raw_repo['org_id'] = org
            # remove username from url
            url = urlparse(raw_repo['remoteUrl'])
            url = url._replace(netloc=url.hostname)
            raw_repo['url'] = url.geturl()
            repo = GitRepository(**raw_repo)
            if not repo.default_branch:
                continue
            if "parentRepository" in raw_repo:
                repo.parent_repository_url = raw_repo["parentRepository"]["url"]
            yield repo

    def test_connection(self, connection: AzureDevOpsConnection):
        api = AzureDevOpsAPI(connection)
        if connection.organization is None:
            try:
                api.my_profile()
            except APIException as e:
                if e.response.status == 401:
                    raise Exception(f"Invalid token {e}. You may need to set organization name in connection or edit your token to set organization to 'All accessible organizations'")
                raise
        else:
            try:
                api.projects(connection.organization)
            except APIException as e:
                raise Exception(f"Invalid token: {e}")

    def extra_tasks(self, scope: GitRepository, tx_rule: AzureDevOpsTransformationRule, entity_types: list[DomainType], connection: AzureDevOpsConnection):
        if DomainType.CODE in entity_types:
            url = urlparse(scope.remoteUrl)
            url = url._replace(netloc=f'{url.username}:{connection.token}@{url.hostname}')
            yield gitextractor(url.geturl(), scope.domain_id(), connection.proxy)

    def extra_stages(self, scope_tx_rule_pairs: list[ScopeTxRulePair], entity_types: list[DomainType], _):
        if DomainType.CODE in entity_types:
            for scope, tx_rule in scope_tx_rule_pairs:
                options = tx_rule.refdiff if tx_rule else None
                yield [refdiff(scope.id, options)]

    @property
    def streams(self):
        return [
            GitPullRequests,
            GitPullRequestCommits,
            Builds,
            Jobs,
        ]


if __name__ == '__main__':
    AzureDevOpsPlugin.start()
