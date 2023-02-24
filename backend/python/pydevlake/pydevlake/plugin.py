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


from typing import Type, Union, Iterable
import sys
from abc import ABC, abstractmethod
import requests

import fire

import pydevlake.message as msg
from pydevlake.subtasks import Subtask
from pydevlake.docgen import generate_doc
from pydevlake.ipc import PluginCommands
from pydevlake.context import Context
from pydevlake.stream import Stream
from pydevlake.model import ToolScope, DomainScope, Connection, TransformationRule


class Plugin(ABC):
    def __init__(self):
        self._streams = dict()
        for stream in self.streams:
            if isinstance(stream, type):
                stream = stream(self.name)
            self._streams[stream.name] = stream

    @property
    def name(self) -> str:
        """
        The name of the plugin, defaults to the class name lowercased.
        """
        return type(self).__name__.lower().removesuffix('plugin')

    @property
    def description(self) -> str:
        return f"{self.name} plugin"

    @property
    @abstractmethod
    def connection_type(self) -> Type[Connection]:
        pass

    @property
    @abstractmethod
    def tool_scope_type(self) -> Type[ToolScope]:
        pass

    @property
    def transformation_rule_type(self) -> Type[TransformationRule]:
        return None

    @abstractmethod
    def test_connection(self, connection: Connection):
        """
        Test if the the connection with the datasource can be established with the given connection.
        Must raise an exception if the connection can't be established.
        """
        pass

    @property
    def subtasks(self) -> list[Subtask]:
        return [subtask for stream in self._streams.values() for subtask in stream.subtasks]

    @abstractmethod
    def get_domain_scopes(self, scope_name: str, connection: Connection) -> Iterable[DomainScope]:
        pass

    @abstractmethod
    def remote_scopes(self, connection: Connection, query: str = ''):
        pass

    @property
    def streams(self) -> list[Union[Stream, Type[Stream]]]:
        pass

    def collect(self, ctx: Context, stream: str):
        yield from self.get_stream(stream).collector.run(ctx)

    def extract(self, ctx: Context, stream: str):
        yield from self.get_stream(stream).extractor.run(ctx)

    def convert(self, ctx: Context, stream: str):
        yield from self.get_stream(stream).convertor.run(ctx)

    def run_migrations(self, force: bool):
        # TODO: Create tables
        pass

    def make_pipeline(self, ctx: Context, scopes: list[msg.BlueprintScope]):
        """
        Make a simple pipeline using the scopes declared by the plugin.
        """
        stages = [
            msg.PipelineStage(
                tasks=[
                    msg.PipelineTask(
                        self.name,
                        skipOnFail=False,
                        subtasks=[t.name for t in self.subtasks],
                        options={
                            "scopeId": scope.id,
                            "scopeName": scope.name}
                    )
                ]
            )
            for scope in scopes
        ]

        plan = msg.PipelinePlan(stages=stages)
        yield plan

        scopes = [
            scope
            for bp_scope in scopes
            for scope in self.get_scopes(bp_scope.name, ctx.connection)
        ]
        yield scopes

    def get_stream(self, stream_name: str):
        stream = self._streams.get(stream_name)
        if stream is None:
            raise Exception(f'Unkown stream {stream_name}')
        return stream

    def startup(self, endpoint: str):
        details = msg.PluginDetails(
            plugin_info=self.plugin_info(),
            swagger=msg.SwaggerDoc(
                name=self.name,
                resource=self.name,
                spec=generate_doc(self.name, self.connection_type, self.transformation_rule_type)
            )
        )
        resp = requests.post(f"{endpoint}/plugins/register", data=details.json())
        if resp.status_code != 200:
            raise Exception(f"unexpected http status code {resp.status_code}: {resp.content}")

    def plugin_info(self) -> msg.PluginInfo:
        subtask_metas = [
            msg.SubtaskMeta(
                name=subtask.name,
                entry_point_name=subtask.verb,
                arguments=[subtask.stream.name],
                required=True,
                enabled_by_default=True,
                description=subtask.description,
                domain_types=[dm.value for dm in subtask.stream.domain_types]
            )
            for subtask in self.subtasks
        ]

        if self.transformation_rule_type:
            tx_rule_model_info = msg.DynamicModelInfo.from_model(self.transformation_rule_type)
        else:
            tx_rule_model_info = None

        return msg.PluginInfo(
            name=self.name,
            description=self.description,
            plugin_path=self._plugin_path(),
            extension="datasource",
            connection_model_info=msg.DynamicModelInfo.from_model(self.connection_type),
            transformation_rule_model_info=tx_rule_model_info,
            scope_model_info=msg.DynamicModelInfo.from_model(self.tool_scope_type),
            subtask_metas=subtask_metas
        )

    def _plugin_path(self):
        module_name = type(self).__module__
        module = sys.modules[module_name]
        return module.__file__

    @classmethod
    def start(cls):
        plugin = cls()
        fire.Fire(PluginCommands(plugin))
