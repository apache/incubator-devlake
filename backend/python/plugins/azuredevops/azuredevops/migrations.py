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
import re
from enum import Enum
from typing import Optional

from pydantic import SecretStr

from pydevlake import ToolModel, Connection, Field
from pydevlake.migration import migration, Dialect, MigrationScriptBuilder
from pydevlake.model import Model, ToolTable
from pydevlake.pipeline_tasks import RefDiffOptions


@migration(20230501000001, name="initialize schemas for Azure Devops")
def init_schemas(b: MigrationScriptBuilder):
    class AzureDevOpsConnection(Connection):
        token: SecretStr
        organization: Optional[str]

    class AzureDevopsTransformationRule(ToolTable, Model):
        name: str = Field(default="default")
        refdiff: Optional[RefDiffOptions]
        deployment_pattern: Optional[re.Pattern]
        production_pattern: Optional[re.Pattern]

    class GitRepository(ToolModel):
        id: str = Field(primary_key=True)
        name: str
        url: str
        remote_url: Optional[str]
        default_branch: Optional[str]
        project_id: str
        org_id: str
        parent_repository_url: Optional[str]
        provider: Optional[str]

    class GitPullRequest(ToolModel):
        class PRStatus(Enum):
            Abandoned = "abandoned"
            Active = "active"
            Completed = "completed"

        pull_request_id: int = Field(primary_key=True)
        description: Optional[str]
        status: PRStatus
        created_by_id: str
        created_by_name: str
        creation_date: datetime.datetime
        closed_date: Optional[datetime.datetime]
        source_commit_sha: str
        target_commit_sha: str
        merge_commit_sha: Optional[str]
        url: Optional[str]
        type: Optional[str]
        title: Optional[str]
        target_ref_name: Optional[str]
        source_ref_name: Optional[str]
        fork_repo_id: Optional[str]

    class GitPullRequestCommit(ToolModel):
        commit_id: str = Field(primary_key=True)
        pull_request_id: str
        author_name: str
        author_email: str
        author_date: datetime.datetime

    class Build(ToolModel):
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
        name: str
        start_time: Optional[datetime.datetime]
        finish_time: Optional[datetime.datetime]
        status: BuildStatus
        result: Optional[BuildResult]
        source_branch: str
        source_version: str

    class Job(ToolModel):
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
        build_id: str
        name: str
        start_time: Optional[datetime.datetime]
        finish_time: Optional[datetime.datetime]
        state: JobState
        result: Optional[JobResult]

    b.create_tables(
        AzureDevOpsConnection,
        AzureDevopsTransformationRule,
        GitRepository,
        GitPullRequestCommit,
        GitPullRequest,
        Build,
        Job,
    )


@migration(20230524181430)
def add_build_id_as_job_primary_key(b: MigrationScriptBuilder):
    # NOTE: We can't add a column to the primary key of an existing table
    # so we have to drop the primary key constraint first,
    # which is done differently in MySQL and PostgreSQL,
    # and then add the new composite primary key.
    table = '_tool_azuredevops_jobs'
    b.execute(f'ALTER TABLE {table} MODIFY COLUMN build_id VARCHAR(255)', Dialect.MYSQL)
    b.execute(f'ALTER TABLE {table} ALTER COLUMN build_id TYPE VARCHAR(255)', Dialect.POSTGRESQL)
    b.execute(f'ALTER TABLE {table} DROP PRIMARY KEY', Dialect.MYSQL)
    b.execute(f'ALTER TABLE {table} DROP CONSTRAINT {table}_pkey', Dialect.POSTGRESQL)
    b.execute(f'ALTER TABLE {table} ADD PRIMARY KEY (id, build_id)')


@migration(20230606165630)
def rename_tx_rule_table_to_scope_config(b: MigrationScriptBuilder):
    b.rename_table('_tool_azuredevops_azuredevopstransformationrules', '_tool_azuredevops_gitrepositoryconfigs')


@migration(20230607165630, name="add entities column to gitrepositoryconfig table")
def add_entities_column_to_scope_config(b: MigrationScriptBuilder):
    b.add_column('_tool_azuredevops_gitrepositoryconfigs', 'entities', 'json')


@migration(20230630000001, name="populated _raw_data_table column for azuredevops git repos")
def add_raw_data_params_table_to_scope(b: MigrationScriptBuilder):
    b.execute(f'''UPDATE _tool_azuredevops_gitrepositories SET _raw_data_table = '_raw_azuredevops_scopes' WHERE 1=1''')
