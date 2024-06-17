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
from datetime import datetime

from pydantic import SecretStr

from pydevlake import ToolModel, Connection, Field
from pydevlake.migration import migration, MigrationScriptBuilder
from pydevlake.model import ScopeConfig, ToolScope
from pydevlake.pipeline_tasks import RefDiffOptions


@migration(20240108000001, name="initialize schemas for gerrit")
def init_schemas(b: MigrationScriptBuilder):
    class GerritConnection(Connection):
        endpoint: str
        username: Optional[str]
        password: Optional[SecretStr]
        pattern: Optional[str]

    class GerritProject(ToolScope):
        name: str
        url: str

    class GerritProjectConfig(ScopeConfig):
        refdiff: Optional[RefDiffOptions]

    class GerritChange(ToolModel):
        id: str = Field(primary_key=True)
        change_id: str
        change_number: int
        subject: str
        status: str
        branch: str
        created_date: datetime
        merged_date: Optional[datetime]
        closed_date: Optional[datetime]
        current_revision: Optional[str]
        owner_name: Optional[str]
        owner_email: Optional[str]
        revisions_json: Optional[str]

    class GerritChangeCommit(ToolModel):
        commit_id: str = Field(primary_key=True)
        pull_request_id: str
        author_name: str
        author_email: str
        author_date: datetime

    b.create_tables(
        GerritConnection,
        GerritProject,
        GerritProjectConfig,
        GerritChange,
        GerritChangeCommit,
    )
