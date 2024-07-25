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

import json
from datetime import datetime
from typing import Optional

from pydevlake import Field
from pydevlake.model import ScopeConfig, ToolScope, ToolModel, Connection
from pydantic import SecretStr
from pydevlake.pipeline_tasks import RefDiffOptions

# needed to be able to run migrations
from gerrit.migrations import *  # noqa


class GerritConnection(Connection):
    endpoint: str
    username: Optional[str]
    password: Optional[SecretStr]
    pattern: Optional[str]

    @property
    def url(self):
        if self.endpoint.endswith("/"):
            return self.endpoint
        return self.endpoint + "/"


class GerritProjectConfig(ScopeConfig):
    refdiff: Optional[RefDiffOptions]


class GerritProject(ToolScope, table=True):
    name: str
    url: str


class GerritChange(ToolModel, table=True):
    id: str = Field(primary_key=True)
    change_id: str
    change_number: int = Field(source="/_number")
    subject: str
    status: str
    branch: str
    created_date: datetime
    merged_date: Optional[datetime]
    closed_date: Optional[datetime]
    current_revision: Optional[str]
    owner_name: Optional[str] = Field(source="/owner/name")
    owner_email: Optional[str] = Field(source="/owner/email")
    revisions_json: Optional[str] = Field(source="/revisions_json")

    @property
    def revisions(self):
        return json.loads(self.revisions_json)


class GerritChangeCommit(ToolModel, table=True):
    commit_id: str = Field(primary_key=True)
    pull_request_id: str
    author_name: str
    author_email: str
    author_date: datetime
