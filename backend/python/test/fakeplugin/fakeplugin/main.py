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

from enum import Enum
from datetime import datetime
from typing import Optional
import json

from sqlmodel import Field

from pydevlake import Plugin, Connection, TransformationRule, Stream, ToolModel, ToolScope, RemoteScopeGroup, DomainType
from pydevlake.domain_layer.devops import CicdScope, CICDPipeline, CICDStatus, CICDResult, CICDType


VALID_TOKEN = "this_is_a_valid_token"


class FakePipeline(ToolModel, table=True):
    class State(Enum):
        PENDING = "pending"
        RUNNING = "running"
        FAILURE = "failure"
        SUCCESS = "success"

    id: str = Field(primary_key=True)
    started_at: Optional[datetime]
    finished_at: Optional[datetime]
    state: State


class FakeStream(Stream):
    tool_model = FakePipeline
    domain_types = [DomainType.CICD]

    fake_pipelines = [
        FakePipeline(id=1, state=FakePipeline.State.SUCCESS, started_at=datetime(2023, 1, 10, 11, 0, 0), finished_at=datetime(2023, 1, 10, 11, 3, 0)),
        FakePipeline(id=2, state=FakePipeline.State.FAILURE, started_at=datetime(2023, 1, 10, 12, 0, 0), finished_at=datetime(2023, 1, 10, 12, 1, 30)),
        FakePipeline(id=1, state=FakePipeline.State.PENDING),
    ]

    def collect(self, state, context):
        for p in self.fake_pipelines:
            yield json.loads(p.json()), {}

    def convert(self, pipeline: FakePipeline, ctx):
        if ctx.transformation_rule:
            env = ctx.transformation_rule.env
        else:
            env = "unknown"
        yield CICDPipeline(
            name=pipeline.id,
            status=self.convert_status(pipeline.state),
            finished_date=pipeline.finished_at,
            result=self.convert_result(pipeline.state),
            duration_sec=self.duration(pipeline),
            environment=env,
            type=CICDType.BUILD
        )

    def convert_status(self, state: FakePipeline.State):
        match state:
            case FakePipeline.State.FAILURE | FakePipeline.State.SUCCESS:
                return CICDStatus.DONE
            case _:
                return CICDStatus.IN_PROGRESS

    def convert_result(self, state: FakePipeline.State):
        match state:
            case FakePipeline.State.SUCCESS:
                return CICDResult.SUCCESS
            case FakePipeline.State.FAILURE:
                return CICDResult.FAILURE
            case _:
                return None

    def duration(self, pipeline: FakePipeline):
        if pipeline.finished_at:
            return (pipeline.finished_at - pipeline.started_at).seconds
        return None


class FakeConnection(Connection):
    token: str


class FakeProject(ToolScope, table=True):
    pass


class FakeTransformationRule(TransformationRule):
    env: str


class FakePlugin(Plugin):
    @property
    def connection_type(self):
        return FakeConnection

    @property
    def tool_scope_type(self):
        return FakeProject

    @property
    def transformation_rule_type(self):
        return FakeTransformationRule

    def domain_scopes(self, project: FakeProject):
        yield CicdScope(
            id=1,
            name=project.name,
            url=f"http://fake.org/api/project/{project.name}"
        )

    def remote_scopes(self, connection: FakeConnection, group_id: str):
        if group_id == 'group1':
            return [
                FakeProject(
                    id='p1',
                    name='Project 1'
                )
            ]
        else:
            return []

    def remote_scope_groups(self, connection: FakeConnection):
        return [
            RemoteScopeGroup(
                id='group1',
                name='Group 1'
            )
        ]

    def test_connection(self, connection: FakeConnection):
        if connection.token != VALID_TOKEN:
            raise Exception("Invalid token")

    @property
    def streams(self):
        return [
            FakeStream
        ]


if __name__ == '__main__':
    FakePlugin.start()
