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

from typing import Optional
import base64

from pydevlake.api import API, request_hook, Paginator, Request


class AzurePaginator(Paginator):
    def get_items(self, response) -> Optional[list[object]]:
        return response.json['value']

    def get_next_page_id(self, response) -> Optional[str]:
        return response.headers.get('x-ms-continuation')

    def set_next_page_param(self, request, next_page_id):
        request.query_args['continuationToken'] = next_page_id


class AzureDevOpsAPI(API):
    paginator = AzurePaginator()
    base_url = "https://dev.azure.com/"

    @request_hook
    def authenticate(self, request: Request):
        token_b64 = base64.b64encode((':' + self.connection.token).encode()).decode()
        request.headers['Authorization'] = 'Basic ' + token_b64

    @request_hook
    def set_api_version(self, request: Request):
        request.query_args['api-version'] = "7.0"

    def my_profile(self):
        req = Request('https://app.vssps.visualstudio.com/_apis/profile/profiles/me')
        return self.send(req)

    def accounts(self, member_id: str):
        req = Request('https://app.vssps.visualstudio.com/_apis/accounts', query_args={"memberId": member_id})
        return self.send(req)

    def projects(self, org: str):
        return self.get(org, '_apis/projects')

    def git_repos(self, org: str, project: str):
        return self.get(org, project, '_apis/git/repositories')

    def git_repo_pull_requests(self, org: str, project: str, repo_id: str):
        return self.get(org, project, '_apis/git/repositories', repo_id, 'pullrequests?searchCriteria.status=all')

    def git_repo_pull_request_commits(self, org: str, project: str, repo_id: str, pull_request_id: int):
        return self.get(org, project, '_apis/git/repositories', repo_id, 'pullRequests', pull_request_id, 'commits')

    def git_repo_pull_request_comments(self, org: str, project: str, repo_id: str, pull_request_id: int):
        return self.get(org, project, '_apis/git/repositories', repo_id, 'pullRequests', pull_request_id, 'threads')

    def commits(self, org: str, project: str, repo_id: str):
        return self.get(org, project, '_apis/git/repositories', repo_id, 'commits')

    def builds(self, org: str, project: str, repository_id: str, provider: str):
        return self.get(org, project, '_apis/build/builds', repositoryId=repository_id, repositoryType=provider)

    def jobs(self, org: str, project: str, build_id: int):
        return self.get(org, project, '_apis/build/builds', build_id, 'timeline')
