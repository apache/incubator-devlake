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

from pydevlake import ScopeConfig, Field
from pydevlake.model import ToolScope, ToolModel, Connection
from pydevlake.pipeline_tasks import RefDiffOptions

# needed to be able to run migrations
from azuredevops.migrations import *


class AzureDevOpsConnection(Connection):
    token: SecretStr
    organization: Optional[str]


class GitRepositoryConfig(ScopeConfig):
    refdiff: Optional[RefDiffOptions]
    deployment_pattern: Optional[re.Pattern]
    production_pattern: Optional[re.Pattern]


class GitRepository(ToolScope, table=True):
    url: str
    remote_url: Optional[str]
    default_branch: Optional[str]
    project_id: str
    org_id: str
    parent_repository_url: Optional[str] = Field(source='/parentRepository/url')
    provider: Optional[str]
    updated_date: datetime.datetime = Field(source='/project/lastUpdateTime')

    def is_external(self):
        return bool(self.provider)


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
    type: Optional[str] = Field(source='/labels/0/name')  # TODO: Add regex to scope config
    title: Optional[str]
    target_ref_name: Optional[str]
    source_ref_name: Optional[str]
    fork_repo_id: Optional[str] = Field(source='/forkSource/repository/id')


class GitPullRequestCommit(ToolModel, table=True):
    commit_id: str = Field(primary_key=True)
    pull_request_id: str
    author_name: str = Field(source='/author/name')
    author_email: Optional[str] = Field(source='/author/email')
    author_date: datetime.datetime = Field(source='/author/date')


class Build(ToolModel, table=True):
    class BuildStatus(Enum):
        Cancelling = "cancelling"
        Completed = "completed"
        InProgress = "inProgress"
        NotStarted = "notStarted"
        Postponed = "postponed"

        def __str__(self) -> str:
            return self.name

    class BuildResult(Enum):
        Canceled = "canceled"
        Failed = "failed"
        Non = "none"
        PartiallySucceeded = "partiallySucceeded"
        Succeeded = "succeeded"

        def __str__(self) -> str:
            return self.name

    id: int = Field(primary_key=True)
    name: str = Field(source='/definition/name')
    queue_time: Optional[datetime.datetime] = Field(source='/queueTime')
    start_time: Optional[datetime.datetime]
    finish_time: Optional[datetime.datetime]
    status: BuildStatus
    result: Optional[BuildResult]
    source_branch: str
    source_version: str
    display_title: Optional[str] = Field(source='/triggerInfo/ci.message')
    url: Optional[str] = Field(source='/_links/web/href')


class Job(ToolModel, table=True):
    class JobState(Enum):
        Completed = "completed"
        InProgress = "inProgress"
        Pending = "pending"

        def __str__(self) -> str:
            return self.name

    class JobResult(Enum):
        Abandoned = "abandoned"
        Canceled = "canceled"
        Failed = "failed"
        Skipped = "skipped"
        Succeeded = "succeeded"
        SucceededWithIssues = "succeededWithIssues"

        def __str__(self) -> str:
            return self.name

    id: str = Field(primary_key=True)
    build_id: str = Field(primary_key=True)
    name: str
    start_time: Optional[datetime.datetime]
    finish_time: Optional[datetime.datetime]
    state: JobState
    result: Optional[JobResult]
