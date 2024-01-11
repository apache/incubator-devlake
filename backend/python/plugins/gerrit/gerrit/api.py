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

from base64 import b64encode
from typing import Optional
from pydevlake.api import API, Request, Response, request_hook, response_hook, Paginator


# TODO: implement pagination
class GerritPaginator(Paginator):
    def get_items(self, response) -> Optional[list[object]]:
        return response.json

    def get_next_page_id(self, response) -> Optional[str]:
        return []

    def set_next_page_param(self, request, next_page_id):
        pass


class GerritApi(API):
    # paginator = GerritPaginator()

    @property
    def base_url(self):
        return self.connection.url

    @request_hook
    def authenticate(self, request: Request):
        conn = self.connection
        if conn.username and conn.password:
            user_and_pass = f"{conn.username}:{conn.password.get_secret_value()}".encode()
            basic_auth = b64encode(user_and_pass).decode()
            request.headers["Authorization"] = f"Basic {basic_auth}"

    @response_hook
    def remove_extra_content_in_json(self, response: Response):
        # remove ")]}'"
        if response.body.startswith(b")]}'"):
            response.body = response.body[4:]

    def my_profile(self):
        return self.get("accounts/self")

    def projects(self):
        # TODO: use pagination
        projects_uri = "projects/?type=CODE&n=10000"
        if self.connection.pattern:
            projects_uri += f"&r={self.connection.pattern}"
        return self.get(projects_uri)

    def changes(self, project_name: str):
        # TODO: use pagination
        return self.get(f"changes/?q=p:{project_name}&o=CURRENT_REVISION&o=ALL_COMMITS&o=DETAILED_ACCOUNTS&no-limit")

    def change_detail(self, change_id: str):
        return self.get(f"changes/{change_id}/detail")

    def account(self, account_id: int):
        return self.get(f"accounts/{account_id}")
