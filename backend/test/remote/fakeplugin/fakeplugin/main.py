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

from sqlmodel import Field

from pydevlake import Plugin, Connection, TransformationRule, Stream, ToolModel, RemoteScope, DomainType
from pydevlake.domain_layer.devops import CicdScope, CICDPipeline


VALID_TOKEN = "this_is_a_valid_token"
VALID_PROJECT = "this_is_a_valid_project"


class FakePipeline(ToolModel, table=True):
    class State(Enum):
        PENDING = "pending"
        RUNNING = "running"
        FAILURE = "failure"
        SUCCESS = "success"

    id: str = Field(primary_key=True)
    project: str
    started_at: Optional[datetime]
    finished_at: Optional[datetime]
    state: State


class FakeStream(Stream):
    tool_model = FakePipeline
    domain_types = [DomainType.CICD]

    fake_pipelines = [
        FakePipeline(id=1, project=VALID_PROJECT, state=FakePipeline.State.SUCCESS, started_at=datetime(2023, 1, 10, 11, 0, 0), finished_at=datetime(2023, 1, 10, 11, 3, 0)),
        FakePipeline(id=2, project=VALID_PROJECT, state=FakePipeline.State.FAILURE, started_at=datetime(2023, 1, 10, 12, 0, 0), finished_at=datetime(2023, 1, 10, 12, 1, 30)),
        FakePipeline(id=1, project=VALID_PROJECT, state=FakePipeline.State.PENDING),
    ]

    def collect(self, state, context):
        project = context.options['project']
        if project == VALID_PROJECT:
            for p in self.fake_pipelines:
                yield dict(p)

    def convert(self, pipeline: FakePipeline):
        yield CICDPipeline(
            name=pipeline.id,
            status=self.convert_status(pipeline.state),
            finished_date=pipeline.finished_at,
            result=self.convert_result(pipeline.state),
            duration_sec=self.duration(pipeline),
            environment=[],
            type=CICDPipeline.Type.CI
        )

    def convert_status(self, state: FakePipeline.State):
        match state:
            case FakePipeline.State.FAILURE | FakePipeline.State.SUCCESS:
                return CICDPipeline.Status.DONE
            case _:
                return CICDPipeline.Status.IN_PROGRESS

    def convert_result(self, state: FakePipeline.State):
        match state:
            case FakePipeline.State.SUCCESS:
                return CICDPipeline.Result.SUCCESS
            case FakePipeline.State.FAILURE:
                return CICDPipeline.Status.FAILURE
            case _:
                return None

    def duration(self, pipeline: FakePipeline):
        if pipeline.finished_at:
            return (pipeline.finished_at - pipeline.started_at).seconds
        return None


class FakeConnection(Connection):
    token: str


class FakeTransformationRule(TransformationRule):
    tx1: str


class FakePlugin(Plugin):
    @property
    def connection_type(self):
        return FakeConnection

    @property
    def tool_scope_type(self):
        return CicdScope

    def get_domain_scopes(self, scope_name: str, connection: FakeConnection):
        yield CicdScope(
            id=1,
            name=scope_name,
            url=f"http://fake.org/api/project/{scope_name}"
        )

    def remote_scopes(self, connection: FakeConnection, query: str = ''):
        yield RemoteScope(
            id='test',
            name='Not a real scope'
        )

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
