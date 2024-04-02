import os

import pytest

from pydevlake.migration import migration, MigrationScriptBuilder
from pydevlake.model import SubtaskRun

@pytest.fixture
def endpoint():
    return os.environ.get('GRAFANA_ONCALL_ENDPOINT') or pytest.skip("No Grafana OnCall endpoint provided")


@pytest.fixture
def token():
    return os.environ.get('GRAFANA_ONCALL_TOKEN') or pytest.skip("No Grafana OnCall token provided")
