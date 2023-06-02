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

import datetime
from enum import Enum
from typing import Optional
import re

from pydantic import SecretStr

from pydevlake import Field, Connection, TransformationRule
from pydevlake.model import ToolModel, ToolScope
from pydevlake.pipeline_tasks import RefDiffOptions
from pydevlake.migration import migration, MigrationScriptBuilder, Dialect


class AzureDevOpsConnection(Connection):
    token: SecretStr
    organization: Optional[str]


class AzureDevOpsTransformationRule(TransformationRule):
    refdiff: Optional[RefDiffOptions]
    deployment_pattern: Optional[re.Pattern]
    production_pattern: Optional[re.Pattern]


class GitRepository(ToolScope, table=True):
    url: str
    remote_url: str
    default_branch: Optional[str]
    project_id: str
    org_id: str
    parent_repository_url: Optional[str] = Field(source='parentRepository/url')


class GitPullRequest(ToolModel, table=True):
    class PRStatus(Enum):
        Abandoned = "abandoned"
        Active = "active"
        Completed = "completed"

    pull_request_id: int = Field(primary_key=True)
    description: Optional[str]
    status: PRStatus
    created_by_id: str = Field(source='/createdBy/id')
    created_by_name: str = Field(source='/createdBy/displayName')
    creation_date: datetime.datetime
    closed_date: Optional[datetime.datetime]
    source_commit_sha: str = Field(source='/lastMergeSourceCommit/commitId')
    target_commit_sha: str = Field(source='/lastMergeTargetCommit/commitId')
    merge_commit_sha: Optional[str] = Field(source='/lastMergeCommit/commitId')
    url: Optional[str]
    type: Optional[str] = Field(source='/labels/0/name') # TODO: get this off transformation rules regex
    title: Optional[str]
    target_ref_name: Optional[str]
    source_ref_name: Optional[str]
    fork_repo_id: Optional[str] = Field(source='/forkSource/repository/id')


class GitPullRequestCommit(ToolModel, table=True):
    commit_id: str = Field(primary_key=True)
    pull_request_id: str
    author_name: str = Field(source='/author/name')
    author_email: str = Field(source='/author/email')
    author_date: datetime.datetime = Field(source='/author/date')


class Build(ToolModel, table=True):
    class BuildStatus(Enum):
        Cancelling = "cancelling"
        Completed = "completed"
        InProgress = "inProgress"
        NotStarted = "notStarted"
        Postponed = "postponed"

    class BuildResult(Enum):
        Canceled = "canceled"
        Failed = "failed"
        Non = "none"
        PartiallySucceeded = "partiallySucceeded"
        Succeeded = "succeeded"

    id: int = Field(primary_key=True)
    name: str = Field(source='/definition/name')
    start_time: Optional[datetime.datetime]
    finish_time: Optional[datetime.datetime]
    status: BuildStatus
    result: Optional[BuildResult]
    source_branch: str
    source_version: str


class Job(ToolModel, table=True):
    class JobState(Enum):
        Completed = "completed"
        InProgress = "inProgress"
        Pending = "pending"

    class JobResult(Enum):
        Abandoned = "abandoned"
        Canceled = "canceled"
        Failed = "failed"
        Skipped = "skipped"
        Succeeded = "succeeded"
        SucceededWithIssues = "succeededWithIssues"

    id: str = Field(primary_key=True)
    build_id: str = Field(primary_key=True)
    name: str
    start_time: Optional[datetime.datetime]
    finish_time: Optional[datetime.datetime]
    state: JobState
    result: Optional[JobResult]


@migration(20230524181430)
def add_build_id_as_job_primary_key(b: MigrationScriptBuilder):
    # NOTE: We can't add a column to the primary key of an existing table
    # so we have to drop the primary key constraint first,
    # which is done differently in MySQL and PostgreSQL,
    # and then add the new composite primary key.
    table = Job.__tablename__
    b.execute(f'ALTER TABLE {table} DROP PRIMARY KEY', Dialect.MYSQL)
    b.execute(f'ALTER TABLE {table} DROP CONSTRAINT {table}_pkey', Dialect.POSTGRESQL)
    b.execute(f'ALTER TABLE {table} ADD PRIMARY KEY (id, build_id)')
