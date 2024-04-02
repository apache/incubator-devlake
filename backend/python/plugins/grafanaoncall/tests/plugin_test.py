import os
import pytest

from pydevlake.testing import assert_valid_plugin, assert_plugin_run

from grafanaoncall.main import GrafanaOncallPlugin
from grafanaoncall.models import GrafanaOncallConnection, GrafanaOncallScopeConfig


def test_valid_plugin():
    assert_valid_plugin(GrafanaOncallPlugin())


def test_valid_plugin_and_connection(endpoint, token):
    plugin = GrafanaOncallPlugin()
    connection = GrafanaOncallConnection(id=1, name='test_connection', endpoint=endpoint, token=token)
    scope_config = GrafanaOncallScopeConfig(id=1, name='test_config')

    assert_plugin_run(plugin, connection, scope_config)
