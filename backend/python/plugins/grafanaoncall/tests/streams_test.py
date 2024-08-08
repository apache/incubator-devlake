import pytest

from grafanaoncall.main import GrafanaOncallPlugin
from pydevlake.testing import assert_stream_run, ContextBuilder


@pytest.fixture
def context(endpoint, token):
    return (
        ContextBuilder(GrafanaOncallPlugin())
        .with_connection(endpoint=endpoint, token=token)
        .with_scope(id='FANGXUH8ZYJ74', name='at-test')
        .build()
    )


def test_issues_stream(context):
    stream = GrafanaOncallPlugin().get_stream('alertgroups')

    connection = context.connection
    scope = context.scope
    scope_config = context.scope_config

    assert_stream_run(stream, connection, scope, scope_config)
