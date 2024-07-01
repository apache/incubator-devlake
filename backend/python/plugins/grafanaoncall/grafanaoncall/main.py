from typing import Iterable

from grafanaoncall.oncall_api import GrafanaOncallAPI
from grafanaoncall.models import GrafanaOncallConnection, GrafanaOncallScopeConfig, GrafanaOncallToolScope
from grafanaoncall.streams.alert_groups import AlertGroups
from grafanaoncall.streams.users import Users
from pydevlake import DomainScope, Plugin, RemoteScopeGroup, TestConnectionResult
from pydevlake.api import APIException
from pydevlake.domain_layer.issue_tracking import Board


class GrafanaOncallPlugin(Plugin):
    connection_type = GrafanaOncallConnection
    tool_scope_type = GrafanaOncallToolScope
    scope_config_type = GrafanaOncallScopeConfig
    streams = [AlertGroups, Users]

    def domain_scopes(self, tool_scope: GrafanaOncallToolScope) -> Iterable[DomainScope]:
        yield Board(
            id=tool_scope.id,
            name=tool_scope.name,
            description=tool_scope.description,
            url=tool_scope.url,
        )

    def remote_scope_groups(self, connection: GrafanaOncallConnection) -> Iterable[RemoteScopeGroup]:
        yield RemoteScopeGroup(
            id='oncall',
            name='OnCall'
        )

    def remote_scopes(self, connection: GrafanaOncallConnection, group_id: str) -> Iterable[GrafanaOncallToolScope]:
        api = GrafanaOncallAPI(connection)
        for item in api.escalation_chains():
            yield GrafanaOncallToolScope(
                id=item['id'],
                name=item['name'],
                url=f'{connection.grafana_endpoint}/a/grafana-oncall-app/escalations/{item["id"]}'
            )

    def test_connection(self, connection: GrafanaOncallConnection) -> TestConnectionResult:
        api = GrafanaOncallAPI(connection)
        message = None
        try:
            res = api.escalation_chains()
        except APIException as e:
            res = e.response
            if res.status == 403:
                message = f"Invalid token."
        return TestConnectionResult.from_api_response(res, message)


if __name__ == '__main__':
    GrafanaOncallPlugin.start()
