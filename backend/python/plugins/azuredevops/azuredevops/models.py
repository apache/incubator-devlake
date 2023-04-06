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

from sqlmodel import Field

from pydevlake import Connection, TransformationRule
from pydevlake.model import ToolModel, ToolScope
from pydevlake.pipeline_tasks import RefDiffOptions


class AzureDevOpsConnection(Connection):
    token: str
    organization: Optional[str]


class AzureDevOpsTransformationRule(TransformationRule):
    refdiff: Optional[RefDiffOptions]
    deployment_pattern: Optional[re.Pattern]
    production_pattern: Optional[re.Pattern]


class GitRepository(ToolScope, table=True):
    url: str
    remoteUrl: str
    default_branch: Optional[str]
    project_id: str
    org_id: str
    parent_repository_url: Optional[str] # parentRepository?/url


class GitPullRequest(ToolModel, table=True):
    class Status(Enum):
        Abandoned = "abandoned"
        Active = "active"
        Completed = "completed"

    pull_request_id: int = Field(primary_key=True)
    description: Optional[str]
    status: Status
    created_by_id: Optional[str]
    created_by_name: Optional[str]
    creation_date: datetime.datetime
    closed_date: Optional[datetime.datetime]
    source_commit_sha: Optional[str] # lastMergeSourceCommit/commitId
    target_commit_sha: Optional[str] # lastMergeTargetCommit/commitId
    merge_commit_sha: Optional[str] # lastMergeCommit/commitId
    url: Optional[str]
    type: Optional[str]
    title: Optional[str]
    target_ref_name: Optional[str]
    source_ref_name: Optional[str]
    fork_repo_id: Optional[str] # forkSource?/repository/id


class GitPullRequestCommit(ToolModel, table=True):
    commit_id: str = Field(primary_key=True)
    pull_request_id: str


class Build(ToolModel, table=True):
    class Status(Enum):
        Cancelling = "cancelling"
        Completed = "completed"
        InProgress = "inProgress"
        NotStarted = "notStarted"
        Postponed = "postponed"

    class Result(Enum):
        Canceled = "canceled"
        Failed = "failed"
        Non = "none"
        PartiallySucceeded = "partiallySucceeded"
        Succeeded = "succeeded"

    id: int = Field(primary_key=True)
    name: str # /definition/name
    start_time: Optional[datetime.datetime]
    finish_time: Optional[datetime.datetime]
    status: Status
    result: Result
    source_branch: str
    source_version: str


class Job(ToolModel, table=True):
    class State(Enum):
        Completed = "completed"
        InProgress = "inProgress"
        Pending = "pending"

    class Result(Enum):
        Abandoned = "abandoned"
        Canceled = "canceled"
        Failed = "failed"
        Skipped = "skipped"
        Succeeded = "succeeded"
        SucceededWithIssues = "succeededWithIssues"

    id: str = Field(primary_key=True)
    build_id: str
    name: str
    startTime: datetime.datetime
    finishTime: datetime.datetime
    state: State
    result: Result
