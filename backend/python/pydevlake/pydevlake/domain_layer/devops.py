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
from enum import Enum

from sqlmodel import Field

from pydevlake.model import DomainModel, NoPKModel, DomainScope


class CICDResult(Enum):
    SUCCESS = "SUCCESS"
    FAILURE = "FAILURE"
    ABORT = "ABORT"
    MANUAL = "MANUAL"


class CICDStatus(Enum):
    IN_PROGRESS = "IN_PROGRESS"
    DONE = "DONE"


class CICDType(Enum):
    TEST = "TEST"
    LINT = "LINT"
    BUILD = "BUILD"
    DEPLOYMENT = "DEPLOYMENT"


class CICDEnvironment(Enum):
    PRODUCTION = "PRODUCTION"
    STAGING = "STAGING"
    TESTING = "TESTING"


class CICDPipeline(DomainModel, table=True):
    __tablename__ = 'cicd_pipelines'
    name: str
    status: Optional[CICDStatus]
    created_date: Optional[datetime]
    finished_date: Optional[datetime]
    result: Optional[CICDResult]
    duration_sec: Optional[int]
    environment: Optional[str]
    type: Optional[CICDType]
    cicd_scope_id: Optional[str]


class CiCDPipelineCommit(NoPKModel, table=True):
    __tablename__ = 'cicd_pipeline_commits'
    pipeline_id: str = Field(primary_key=True)
    commit_sha: str = Field(primary_key=True)
    branch: str
    repo_id: str
    repo_url: str


class CicdScope(DomainScope):
    __tablename__ = 'cicd_scopes'
    name: str
    description: Optional[str]
    url: Optional[str]
    createdDate: Optional[datetime]
    updatedDate: Optional[datetime]


class CICDTask(DomainModel, table=True):
    __tablename__ = 'cicd_tasks'
    name: str
    pipeline_id: str
    result: Optional[CICDResult]
    status: Optional[CICDStatus]
    type: Optional[CICDType]
    environment: Optional[CICDEnvironment]
    duration_sec: int
    started_date: Optional[datetime]
    finished_date: Optional[datetime]
    cicd_scope_id: str
