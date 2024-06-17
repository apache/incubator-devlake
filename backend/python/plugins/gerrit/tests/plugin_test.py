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


from os import environ

from pydevlake.testing import assert_valid_plugin
from pydevlake.testing.testing import assert_plugin_run
from gerrit.models import GerritConnection, GerritProjectConfig
from gerrit.main import GerritPlugin


def test_valid_plugin():
    assert_valid_plugin(GerritPlugin())


def test_valid_plugin_and_connection():
    connection_name = "test_connection"
    connection_url = environ.get("GERRIT_URL", "https://gerrit.onap.org/r/")
    connection_username = environ.get("GERRIT_USERNAME", "")
    connection_password = environ.get("GERRIT_PASSWORD", "")
    plugin = GerritPlugin()
    connection = GerritConnection(
        name=connection_name,
        endpoint=connection_url,
        username=connection_username,
        password=connection_password,
    )
    scope_config = GerritProjectConfig(id=1, name="test_config")
    assert_plugin_run(plugin, connection, scope_config)
