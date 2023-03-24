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


import os
import json
from functools import wraps
from typing import Generator, TextIO, Optional, Union

from pydevlake.context import Context
from pydevlake.message import Message


def plugin_method(func):
    def open_send_channel() -> TextIO:
        fd = 3
        return os.fdopen(fd, 'w')

    def send_output(send_ch: TextIO, obj: object):
        if not isinstance(obj, Message):
            raise Exception(f"Not a message: {obj}")
        send_ch.write(obj.json(exclude_unset=True))
        send_ch.write('\n')
        send_ch.flush()

    @wraps(func)
    def wrapper(self, *args):
        ret = func(self, *args)
        if ret is not None:
            with open_send_channel() as send_ch:
                if isinstance(ret, Generator):
                    for each in ret:
                        send_output(send_ch, each)
                else:
                    send_output(send_ch, ret)
        return None

    return wrapper


class PluginCommands:
    def __init__(self, plugin):
        self._plugin = plugin

    @plugin_method
    def collect(self, ctx: dict, stream: str):
        yield from self._plugin.collect(self._mk_context(ctx), stream)

    @plugin_method
    def extract(self, ctx: dict, stream: str):
        yield from self._plugin.extract(self._mk_context(ctx), stream)

    @plugin_method
    def convert(self, ctx: dict, stream: str):
        yield from self._plugin.convert(self._mk_context(ctx), stream)

    @plugin_method
    def test_connection(self, connection: dict):
        connection = self._parse(connection)
        if "name" not in connection:
            connection["name"] = "Test connection"
        connection = self._plugin.connection_type(**connection)
        self._plugin.test_connection(connection)

    @plugin_method
    def make_pipeline(self, scope_tx_rule_pairs: list[tuple[dict, dict]], entities: list[str], connection: dict):
        connection = self._plugin.connection_type(**self._parse(connection))
        entities = self._parse(entities)
        scope_tx_rule_pairs = [
            (
                self._plugin.tool_scope_type(**self._parse(raw_scope)),
                self._plugin.transformation_rule_type(**self._parse(raw_tx_rule)) if raw_tx_rule else None
            )
            for raw_scope, raw_tx_rule in scope_tx_rule_pairs
        ]
        return self._plugin.make_pipeline(scope_tx_rule_pairs, entities, connection['id'])

    @plugin_method
    def run_migrations(self, force: bool):
        self._plugin.run_migrations(force)

    @plugin_method
    def plugin_info(self):
        return self._plugin.plugin_info()

    @plugin_method
    def remote_scopes(self, connection: dict, group_id: Optional[str] = None):
        connection = self._parse(connection)
        c = self._plugin.connection_type(**connection)
        return self._plugin.make_remote_scopes(c, group_id)

    def startup(self, endpoint: str):
        self._plugin.startup(endpoint)

    def _mk_context(self, data: dict):
        data = self._parse(data)
        db_url = data['db_url']
        scope_dict = self._parse(data['scope'])
        scope = self._plugin.tool_scope_type(**scope_dict)
        connection_dict = self._parse(data['connection'])
        connection = self._plugin.connection_type(**connection_dict)
        if self._plugin.transformation_rule_type:
            transformation_rule_dict = self._parse(data['transformation_rule'])
            transformation_rule = self._plugin.transformation_rule_type(**transformation_rule_dict)
        else:
            transformation_rule = None
        options = data.get('options', {})
        return Context(db_url, scope, connection, transformation_rule, options)

    def _parse(self, data: Union[str, dict, list]) -> Union[dict, list]:
        print(data)
        if isinstance(data, (dict, list)):
            return data
        if isinstance(data, str):
            try:
                return json.loads(data)
            except json.JSONDecodeError as e:
                raise Exception(f"Invalid JSON: {e.msg}")
        raise Exception(f"Invalid argument type: {type(data)}")
