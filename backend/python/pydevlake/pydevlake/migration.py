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

from pydevlake.model_info import DynamicModelInfo

MIGRATION_SCRIPTS = []


class Dialect(Enum):
    MYSQL = "mysql"
    POSTGRESQL = "postgres"


class Execute(BaseModel):
    type: Literal["execute"] = "execute"
    sql: str
    dialect: Optional[Dialect] = None
    ignore_error: bool = False


class AddColumn(BaseModel):
    type: Literal["add_column"] = "add_column"
    table: str
    column: str
    column_type: str


class DropColumn(BaseModel):
    type: Literal["drop_column"] = "drop_column"
    table: str
    column: str


class DropTable(BaseModel):
    type: Literal["drop_table"] = "drop_table"
    table: str


class RenameColumn(BaseModel):
    type: Literal["rename_column"] = "rename_column"
    table: str
    old_name: str
    new_name: str


class RenameTable(BaseModel):
    type: Literal["rename_table"] = "rename_table"
    old_name: str
    new_name: str


class CreateTable(BaseModel):
    type: Literal["create_table"] = "create_table"
    model_info: DynamicModelInfo


Operation = Annotated[
    Union[Execute, AddColumn, DropColumn, RenameColumn, DropTable, RenameTable, CreateTable],
    Field(discriminator="type")
]


class MigrationScript(BaseModel):
    operations: List[Operation]
    version: int
    name: str


class MigrationScriptBuilder:
    def __init__(self):
        self.operations = []

    def execute(self, sql: str, dialect: Optional[Dialect] = None, ignore_error = False):
        """
        Executes a raw SQL statement.
        If dialect is specified the statement will be executed only if the db dialect matches.
        """
        self.operations.append(Execute(sql=sql, dialect=dialect, ignore_error=ignore_error))

    def add_column(self, table: str, column: str, type: str):
        """
        Adds a column to a table if it does not exist.
        """
        self.operations.append(AddColumn(table=table, column=column, column_type=type))

    def drop_column(self, table: str, column: str):
        """
        Drops a column from a table if it exists.
        """
        self.operations.append(DropColumn(table=table, column=column))

    def rename_column(self, table: str, old_name: str, new_name: str):
        """
        Renames a column in a table.
        """
        self.operations.append(RenameColumn(table=table, old_name=old_name, new_name=new_name))

    def drop_table(self, table: str):
        """
        Drops a table if it exists.
        """
        self.operations.append(DropTable(table=table))

    def rename_table(self, old_name: str, new_name: str):
        """
        Renames a table if it exists and the new name is not already taken.
        """
        self.operations.append(RenameTable(old_name=old_name, new_name=new_name))

    def create_tables(self, *model_classes):
        """
        Creates a table if it doesn't exist based on the object's fields.
        """
        for model_class in model_classes:
            self.operations.append(CreateTable(model_info=DynamicModelInfo.from_model(model_class)))


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
        raise err
