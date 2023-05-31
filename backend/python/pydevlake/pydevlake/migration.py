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


from typing import List, Literal, Optional, Union, Annotated
from enum import Enum
from datetime import datetime

from pydantic import BaseModel, Field


MIGRATION_SCRIPTS = []

class Dialect(Enum):
    MYSQL = "mysql"
    POSTGRESQL = "postgres"


class Execute(BaseModel):
    type: Literal["execute"] = "execute"
    sql: str
    dialect: Optional[Dialect] = None


class DropColumn(BaseModel):
    type: Literal["drop_column"] = "drop_column"
    table: str
    column: str


class DropTable(BaseModel):
    type: Literal["drop_table"] = "drop_table"
    table: str


Operation = Annotated[
    Union[Execute, DropColumn, DropTable],
    Field(discriminator="type")
]


class MigrationScript(BaseModel):
    operations: List[Operation]
    version: int
    name: str


class MigrationScriptBuilder:
    def __init__(self):
        self.operations = []

    def execute(self, sql: str, dialect: Optional[Dialect] = None):
        """
        Executes a raw SQL statement.
        If dialect is specified the statement will be executed only if the db dialect matches.
        """
        self.operations.append(Execute(sql=sql, dialect=dialect))

    def drop_column(self, table: str, column: str):
        """
        Drops a column from a table.
        """
        self.operations.append(DropColumn(table=table, column=column))

    def drop_table(self, table: str):
        """
        Drops a table.
        """
        self.operations.append(DropTable(table=table))


def migration(version: int, name: Optional[str] = None):
    """
    Builds a migration script from a function.

    Usage:

    @migration(20230511)
    def change_description_type(b: MigrationScriptBuilder):
        b.exec('ALTER TABLE my_table ...')
    """
    _validate_version(version)

    def wrapper(fn):
        builder = MigrationScriptBuilder()
        fn(builder)
        script = MigrationScript(operations=builder.operations, version=version, name=name or fn.__name__)
        MIGRATION_SCRIPTS.append(script)
        return script
    return wrapper


def _validate_version(version: int):
    str_version = str(version)
    err = ValueError(f"Invalid version {version}, must be in YYYYMMDDhhmmss format")
    if len(str_version) != 14:
        raise err
    try:
        datetime.strptime(str_version, "%Y%m%d%H%M%S")
    except ValueError:
        raise  err
