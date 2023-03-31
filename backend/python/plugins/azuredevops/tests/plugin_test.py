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
import pytest

from pydevlake.testing import assert_valid_plugin, assert_plugin_run

from azuredevops.models import AzureDevOpsConnection, AzureDevOpsTransformationRule
from azuredevops.main import AzureDevOpsPlugin


def test_valid_plugin():
    assert_valid_plugin(AzureDevOpsPlugin())


def test_valid_plugin_and_connection():
    # TODO: Set AZURE_DEVOPS_TOKEN env variable in CI
    token = os.environ.get('AZURE_DEVOPS_TOKEN')
    if not(token):
        pytest.skip("No Azure DevOps token provided")

    plugin = AzureDevOpsPlugin()
    connection = AzureDevOpsConnection(id=1, name='test_connection', token=token)
    tx_rule = AzureDevOpsTransformationRule(id=1, name='test_rule')

    assert_plugin_run(plugin, connection, tx_rule)
