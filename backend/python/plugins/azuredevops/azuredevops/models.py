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

from sqlmodel import Field

from pydevlake import Connection
from pydevlake.model import ToolModel, ToolScope

default_date = datetime.datetime.fromisoformat("1970-01-01")


class AzureDevOpsConnection(Connection):
    token: str


class Project(ToolModel, table=True):
    id: str = Field(primary_key=True)
    name: str
    url: str


class GitRepository(ToolScope, table=True):
    url: str
    sshUrl: str
    remoteUrl: str
    defaultBranch: Optional[str]
    project_id: str  # = Field(foreign_key=Project.id)
    org_id: str
    size: int
    isDisabled: bool
    isInMaintenance: bool
    isFork: Optional[bool]
    parentRepositoryUrl: Optional[str]


# https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests?view=azure-devops-rest-7.1&tabs=HTTP#identityrefwithvote
class GitPullRequest(ToolModel, table=True):
    class Status(Enum):
        Abandoned = "abandoned"
        Active = "active"
        All = "all"
        Completed = "completed"
        NotSet = "notSet"

    id: int = Field(primary_key=True)
    project_id: Optional[str]
    description: Optional[str]
    code_review_id: int = 0
    repo_id: Optional[str]
    status: Status
    created_by_id: Optional[str]
    created_by_name: Optional[str]
    creation_date: datetime.datetime = default_date
    closed_date: datetime.datetime = default_date
    source_commit_sha: Optional[str]  # lastmergesourcecommit #base
    target_commit_sha: Optional[str]  # lastmergetargetcommit #head
    merge_commit_sha: Optional[str]
    url: Optional[str]
    type: Optional[str]
    title: Optional[str]
    target_ref_name: Optional[str]
    source_ref_name: Optional[str]
    fork_repo_id: Optional[str]


class GitCommit(ToolModel, table=True):
    commit_sha: str = Field(primary_key=True)
    project_id: str
    repo_id: str
    committer_name: str = ""
    committer_email: str = ""
    commit_date: datetime.datetime = default_date
    author_name: str = ""
    author_email: str = ""
    authored_date: datetime.datetime = default_date
    comment: str = ""
    url: str = ""
    additions: int = 0
    deletions: int = 0


class Account(ToolModel, table=True):
    class Type(Enum):
        Organization = "organization"
        Personal = "personal"

    class Status(Enum):
        Deleted = "deleted"
        Disabled = "disabled"
        Enabled = "enabled"
        Moved = "moved"
        Non = "none"

    account_id: str = Field(primary_key=True)
    account_name: str
    account_owner: str
    account_type: Type
    account_status: Status
    organization_name: str
    namespace_id: str


class Build(ToolModel, table=True):
    class Status(Enum):
        All = "all"
        Cancelling = "cancelling"
        Completed = "completed"
        InProgress = "inProgress"
        Non = "none"
        NotStarted = "notStarted"
        Postponed = "postponed"

    class Priority(Enum):
        AboveNormal = "aboveNormal"
        BelowNormal = "belowNormal"
        High = "high"
        Low = "low"
        Normal = "normal"

    class Result(Enum):
        Canceled = "canceled"
        Failed = "failed"
        Non = "none"
        PartiallySucceeded = "partiallySucceeded"
        Succeeded = "succeeded"

    id: int = Field(primary_key=True)
    project_id: str
    repo_id: str
    repo_type: str
    build_number: str
    build_number_revision: Optional[str]
    controller_id: Optional[str]
    definition_id: Optional[str]
    deleted: Optional[bool]
    start_time: Optional[datetime.datetime]
    finish_time: Optional[datetime.datetime]
    status: Status
    tags: list[str] = []
    priority: Priority
    build_result: Result
    source_branch: str
    source_version: str


class Job(ToolModel, table=True):
    class Type(Enum):
        Task = "Task"
        Job = "Job"
        Checkpoint = "Checkpoint"
        Stage = "Stage"
        Phase = "Phase"

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
    build_id: int
    repo_id: str
    parentId: Optional[str]
    type: Optional[Type]
    name: str
    startTime: datetime.datetime
    finishTime: datetime.datetime
    lastModified: datetime.datetime
    currentOperation: Optional[int]
    percentComplete: Optional[int]
    state: State
    result: Result
    resultCode: Optional[int]
    changeId: Optional[int]
    workerName: Optional[str]
    order: Optional[int]
    errorCount: Optional[int]
    warningCount: Optional[int]
