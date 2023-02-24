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
from functools import wraps
from typing import Generator, TextIO

from pydevlake.context import Context
from pydevlake.message import Message


def plugin_method(func):
    def open_send_channel() -> TextIO:
        fd = 3
        return os.fdopen(fd, 'w')

    def send_output(send_ch: TextIO, obj: object):
        if not isinstance(obj, Message):
            return
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
        connection = self._plugin.connection_type(**connection)
        self._plugin.test_connection(connection)

    @plugin_method
    def make_pipeline(self, ctx: dict, scopes: list[dict]):
        yield from self._plugin.make_pipeline(self._mk_context(ctx), scopes)

    @plugin_method
    def run_migrations(self, force: bool):
        self._plugin.run_migrations(force)

    @plugin_method
    def plugin_info(self):
        return self._plugin.plugin_info()

    @plugin_method
    def remote_scopes(self, connection: dict, query: str = ''):
        c = self._plugin.connection_type(**connection)
        self.plugin.remote_scopes(c, query)

    def startup(self, endpoint: str):
        self._plugin.startup(endpoint)

    def _mk_context(self, data: dict):
        db_url = data['db_url']
        scope_id = data['scope_id']
        connection_id = data['connection_id']
        connection = self._plugin.connection_type(**data['connection'])
        if self._plugin.transformation_rule_type:
            transformation_rule = self._plugin.transformation_rule_type(**data['transformation_rule'])
        else:
            transformation_rule = None
        options = data.get('options', {})
        return Context(db_url, scope_id, connection_id, connection, transformation_rule, options)
