from typing import Iterable

from grafanaoncall.oncall_api import GrafanaOncallAPI
from grafanaoncall.models import AlertGroup, GrafanaOncallToolScope, User
from pydevlake import Context, DomainType, Stream
from pydevlake.domain_layer.crossdomain import Account
from pydevlake.domain_layer.issue_tracking import Issue


class Users(Stream):
    tool_model = User
    domain_types = [DomainType.CROSS]

    def collect(self, state, context: Context) -> Iterable[tuple[object, dict]]:
        api = GrafanaOncallAPI(context.connection)

        for user in api.users():
            yield user, state

    def convert(self, u: User, context: Context):
        yield Account(
            email=u.email,
            user_name=u.username,
            id=f'grafanaoncall:Account:{context.connection.id}:{u.id}',
        )
