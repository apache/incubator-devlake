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

import pytest

from typing import Union, Type, Iterable

from pydevlake.plugin import Plugin
from pydevlake.model import DomainModel


def assert_convert(plugin: Union[Plugin, Type[Plugin]], stream_name: str, raw: dict, expected: Union[DomainModel, Iterable[DomainModel]]):
    if isinstance(plugin, type):
        plugin = plugin()
    stream = plugin.get_stream(stream_name)
    tool_model = stream.extract(raw)
    domain_models = stream.convert(tool_model, None)
    if not isinstance(expected, Iterable):
        expected = [expected]
    if not isinstance(domain_models, Iterable):
        domain_models = [domain_models]
    for res, exp in zip(domain_models, expected):
        assert res == exp
