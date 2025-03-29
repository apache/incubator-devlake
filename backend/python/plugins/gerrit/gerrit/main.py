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

import logging
from urllib.parse import urlparse
from gerrit.streams.changes import GerritChanges
from gerrit.streams.change_commits import GerritChangeCommits
from gerrit.api import GerritApi
from gerrit.models import GerritConnection, GerritProject, GerritProjectConfig

from pydevlake.api import APIException
from pydevlake.domain_layer.code import Repo
from pydevlake.message import (
    PipelineTask,
    RemoteScopeGroup,
    TestConnectionResult,
)
from pydevlake.model import (
    Connection,
    DomainType,
    ScopeConfig,
)
from pydevlake.pipeline_tasks import gitextractor, refdiff
from pydevlake.plugin import Plugin
from pydevlake.stream import Stream


logger = logging.getLogger()


class GerritPlugin(Plugin):
    @property
    def connection_type(self):
        return GerritConnection

    @property
    def tool_scope_type(self):
        return GerritProject

    @property
    def scope_config_type(self):
        return GerritProjectConfig

    def domain_scopes(self, gerrit_project: GerritProject):
        yield Repo(
            name=gerrit_project.name,
            url=gerrit_project.url,
        )

    def remote_scope_groups(self, connection: Connection) -> list[RemoteScopeGroup]:
        yield RemoteScopeGroup(
            id=f"{connection.id}:default",
            name="Code Repositories",
        )

    def remote_scopes(self, connection: Connection, group_id: str) -> list[GerritProject]:
        api = GerritApi(connection)
        json_data = api.projects().json
        for project_name in json_data:
            yield GerritProject(
                id=project_name,
                name=project_name,
                url=connection.url + project_name,
            )

    def test_connection(self, connection: Connection):
        api = GerritApi(connection)
        message = None
        try:
            res = api.projects()
        except APIException as e:
            res = e.response
            message = "HTTP Error: " + str(res.status)
        return TestConnectionResult.from_api_response(res, message)

    def extra_tasks(
        self, scope: GerritProject, config: ScopeConfig, connection: GerritConnection
    ) -> list[PipelineTask]:
        if DomainType.CODE in config.domain_types:
            url = urlparse(scope.url)
            if connection.username and connection.password:
                url = url._replace(
                    netloc=f"{connection.username}:{connection.password.get_secret_value()}@{url.netloc}"
                )
            yield gitextractor(url.geturl(), scope.name, scope.domain_id(), connection.proxy)

    def extra_stages(
        self,
        scope_config_pairs: list[tuple[GerritProject, ScopeConfig]],
        connection: GerritConnection,
    ) -> list[list[PipelineTask]]:
        for scope, config in scope_config_pairs:
            if DomainType.CODE in config.domain_types:
                yield [refdiff(scope.id, config.refdiff)]

    @property
    def streams(self) -> list[Stream]:
        return [
            GerritChanges,
            GerritChangeCommits,
        ]


if __name__ == "__main__":
    GerritPlugin.start()
