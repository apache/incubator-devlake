from datetime import datetime

from pydantic import SecretStr

from pydevlake import Connection as BaseConnection, Field
from pydevlake.migration import Dialect, migration, MigrationScriptBuilder
from pydevlake.model import Model, NoPKModel, ToolModel, ToolTable


@migration(20240131000001, name="initialize schemas for Grafana OnCall")
def init_raw_schemas(b: MigrationScriptBuilder):
    class Connection(BaseConnection):
        endpoint: str
        token: SecretStr

    class ToolScope(ToolModel):
        __tablename__ = "_tool_grafanaoncall_tool_scopes"

        id: str = Field(primary_key=True, auto_increment=False)
        name: str

    class ScopeConfig(ToolTable, Model):
        __tablename__ = "_tool_grafanaoncall_scope_configs"

        connection_id: int
        name: str = Field(default="default")

    class User(ToolModel):
        __tablename__ = "_tool_grafanaoncall_users"
        id: str = Field(primary_key=True, auto_increment=False)
        username: str
        email: str

    class AlertGroup(ToolModel):
        __tablename__ = "_tool_grafanaoncall_alert_groups"

        id: str = Field(primary_key=True, auto_increment=False)
        resolved_at: datetime
        resolved_by: str
        acknowledged_at: datetime
        acknowledged_by: str
        title: str
        url: str
        state: str

    b.create_tables(
        Connection,
        ToolScope,
        ScopeConfig,
        User,
        AlertGroup
    )


@migration(20240412112000, name="add verify and grafana_endpoint fields to _tool_grafanaoncall_connections")
def add_verify_and_grafana_endpoint_fields_to_tool_grafanaoncall_connections(b: MigrationScriptBuilder):
    table = '_tool_grafanaoncall_connections'
    b.execute(f'ALTER TABLE {table} ADD COLUMN verify varchar', Dialect.POSTGRESQL)
    b.execute(f'ALTER TABLE {table} ADD COLUMN verify varchar(255)', Dialect.MYSQL)
    b.execute(f'ALTER TABLE {table} ADD COLUMN grafana_endpoint varchar', Dialect.POSTGRESQL)
    b.execute(f'ALTER TABLE {table} ADD COLUMN grafana_endpoint varchar(255)', Dialect.MYSQL)


@migration(20240412112100, name="add url field to _tool_grafanaoncall_tool_scopes")
def add_url_field_to_tool_grafanaoncall_tool_scopes(b: MigrationScriptBuilder):
    table = '_tool_grafanaoncall_tool_scopes'
    b.execute(f'ALTER TABLE {table} ADD COLUMN url varchar', Dialect.POSTGRESQL)
    b.execute(f'ALTER TABLE {table} ADD COLUMN url varchar(255)', Dialect.MYSQL)
