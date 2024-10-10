from copy import deepcopy
from typing import Iterable

from grafanaoncall.oncall_api import GrafanaOncallAPI
from grafanaoncall.models import AlertGroup, GrafanaOncallToolScope
from pydevlake import Context, DomainType, Stream
from pydevlake.domain_layer.issue_tracking import BoardIssue, Issue


def _route_id_to_escalation_chain_id_mapping(api: GrafanaOncallAPI) -> dict[str, str]:
    result = {}
    for item in api.routes():
        result[item['id']] = item['escalation_chain_id']
    return result


class AlertGroups(Stream):
    tool_model = AlertGroup
    domain_types = [DomainType.CROSS]

    def collect(self, state, context: Context) -> Iterable[tuple[object, dict]]:
        api = GrafanaOncallAPI(context.connection)
        scope: GrafanaOncallToolScope = context.scope

        route_id_to_escalation_chain_id_mapping = _route_id_to_escalation_chain_id_mapping(api)

        # we know that alert groups are sorted by created_at desc in the grafana oncall api
        previous_max_created_at = state.get('max_created_at')
        new_max_created_at = None

        for alert_group in api.alert_groups():
            if previous_max_created_at and alert_group['created_at'] <= previous_max_created_at:
                return

            if route_id_to_escalation_chain_id_mapping.get(alert_group['route_id']) == scope.id:
                if not new_max_created_at:
                    new_max_created_at = alert_group['created_at']

                yield alert_group, {**state, 'max_created_at': new_max_created_at}

    def convert(self, ag: AlertGroup, context: Context):
        connection_id = context.connection.id
        tool_scope_id = context.scope.id

        assignee = ag.resolved_by or ag.acknowledged_by

        i = Issue(
            url=ag.url,
            title=ag.title,
            type=Issue.Type.Incident.value,
            status=(Issue.Status.Done if ag.state == AlertGroup.State.Resolved else Issue.Status.InProgress if ag.state == AlertGroup.State.Acknowledged else Issue.Status.ToDo).value,
            original_status=ag.state,
            resolution_date=ag.resolved_at,
            created_date=ag.created_at,
            updated_date=ag.updated_at,
            # TODO: No field to track when Alert Group was acknowledged
            lead_time_minutes=round((ag.resolved_at - ag.created_at).total_seconds() / 60) if ag.state == AlertGroup.State.Resolved else None,
            assignee_id=assignee and f'grafanaoncall:Account:{connection_id}:{assignee}',
        )
        yield i

        bi = BoardIssue(
            board_id=f'grafanaoncall:GrafanaOncallToolScope:{connection_id}:{tool_scope_id}',
            issue_id=i.id
        )
        yield bi
