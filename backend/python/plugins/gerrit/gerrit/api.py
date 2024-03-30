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
from os import environ
from typing import Optional
from urllib.parse import urlparse
from datetime import datetime, timedelta
from MySQLdb import connect as mysql_connect, Error as MySQLError
from pydevlake.api import API, Request, Response, request_hook, response_hook, Paginator
from gerrit.models import GerritChange


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

    def __init__(self, connection=None):
        super().__init__(connection)
        self.db_conn = None

    def auto_connect(self):
        if self.db_conn:
            try:
                self.db_conn.ping()
                return
            except MySQLError as e:
                self.db_conn.close()
        self.db_conn = None
        if 'DB_URL' in environ:
            parsed_url = urlparse((environ['DB_URL']))
            connection_args = {
                'user': parsed_url.username,
                'password': parsed_url.password,
                'host': parsed_url.hostname,
                'port': parsed_url.port or 3306,  # Default MySQL port
                # Remove leading slash from path
                'database': parsed_url.path[1:]
            }
            try:
                self.db_conn = mysql_connect(**connection_args)
            except MySQLError as e:
                print(f"Error connecting to MySQL: {e}")

    @property
    def base_url(self):
        return self.connection.url

    @request_hook
    def authenticate(self, request: Request):
        conn = self.connection
        if conn.username and conn.password:
            user_pass = f"{conn.username}:{conn.password.get_secret_value()}".encode()
            basic_auth = b64encode(user_pass).decode()
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
        self.auto_connect()
        start_date = None
        if self.db_conn:
            cursor = self.db_conn.cursor()
            try:
                cursor.execute(
                    f"SELECT updated_at FROM _tool_gerrit_gerritchanges WHERE id like '{project_name}~%' ORDER BY updated_at desc limit 1")
                last_updated = cursor.fetchone()
                if last_updated and len(last_updated) > 0:
                    last_updated = last_updated[0] - timedelta(days=1)
                    start_date = datetime.strftime(last_updated, "%Y-%m-%d")
            except MySQLError as e:
                print(f"Error fetching last updated date: {e}")
            cursor.close()
        if start_date:
            return self.get(f"changes/?q=p:{project_name}+after:{start_date}&o=CURRENT_REVISION&o=ALL_COMMITS&o=DETAILED_ACCOUNTS&no-limit")
        return self.get(f"changes/?q=p:{project_name}&o=CURRENT_REVISION&o=ALL_COMMITS&o=DETAILED_ACCOUNTS&no-limit")

    def change_detail(self, change_id: str):
        return self.get(f"changes/{change_id}/detail")

    def account(self, account_id: int):
        return self.get(f"accounts/{account_id}")
