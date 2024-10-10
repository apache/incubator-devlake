from typing import Optional

from enum import Enum

from sqlalchemy import BigInteger, Column

from pydevlake import Connection, ScopeConfig, ToolScope

# needed to be able to run migrations
from grafanaoncall.migrations import *


class User(ToolModel, NoPKModel, table=True):
    __tablename__ = '_tool_grafanaoncall_users'

    id: str = Field(primary_key=True, auto_increment=False)
    username: str
    email: str


class AlertGroup(ToolModel, NoPKModel, table=True):
    class State(Enum):
        Firing = "firing"
        Acknowledged = "acknowledged"
        Resolved = "resolved"

    __tablename__ = "_tool_grafanaoncall_alert_groups"

    id: str = Field(primary_key=True, auto_increment=False)
    created_at: datetime
    resolved_at: datetime
    resolved_by: str
    acknowledged_at: datetime
    acknowledged_by: str
    title: str
    url: str = Field(source='/permalinks/web')
    state: State


class GrafanaOncallConnection(Connection):
    __tablename__ = "_tool_grafanaoncall_connections"

    endpoint: str
    token: SecretStr
    grafana_endpoint: str = ''
    verify: Optional[str]


class GrafanaOncallScopeConfig(ScopeConfig):
    __tablename__ = "_tool_grafanaoncall_scope_configs"


class GrafanaOncallToolScope(ToolScope, table=True):
    __tablename__ = "_tool_grafanaoncall_tool_scopes"

    connection_id: int = Field(sa_column=Column(BigInteger()))
    url: str

    @property
    def description(self):
        return self.name
