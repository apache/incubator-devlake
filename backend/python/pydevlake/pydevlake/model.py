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


import os
from typing import Optional
from inspect import getmodule
from datetime import datetime

import inflect
from pydantic import AnyUrl, validator
from sqlalchemy import Column, DateTime, func
from sqlalchemy.orm import declared_attr
from sqlalchemy.inspection import inspect
from sqlmodel import SQLModel, Field


inflect_engine = inflect.engine()


class Model(SQLModel):
    id: Optional[int] = Field(primary_key=True)
    created_at: Optional[datetime] = Field(
        sa_column=Column(DateTime(), default=func.now())
    )
    updated_at: Optional[datetime] = Field(
        sa_column=Column(DateTime(), default=func.now(), onupdate=func.now())
    )

class ToolTable(Model):
    @declared_attr
    def __tablename__(cls) -> str:
        plugin_name = _get_plugin_name(cls)
        plural_entity = inflect_engine.plural_noun(cls.__name__.lower())
        return f'_tool_{plugin_name}_{plural_entity}'


class Connection(ToolTable):
    name: str
    proxy: Optional[AnyUrl]

    @validator('proxy', pre=True)
    def allow_empty_proxy(cls, proxy):
        if proxy == "":
            return None
        return proxy


class TransformationRule(ToolTable):
    name: str


class RawModel(SQLModel):
    id: int = Field(primary_key=True)
    params: str = b''
    data: bytes
    url: str = b''
    input: bytes = b''
    created_at: datetime = Field(default_factory=datetime.now)


class RawDataOrigin(SQLModel):
    # SQLModel doesn't like attributes starting with _
    # so we change the names of the columns.
    raw_data_params: Optional[str] = Field(sa_column_kwargs={'name':'_raw_data_params'})
    raw_data_table: Optional[str] = Field(sa_column_kwargs={'name':'_raw_data_table'})
    raw_data_id: Optional[str] = Field(sa_column_kwargs={'name':'_raw_data_id'})
    raw_data_remark: Optional[str] = Field(sa_column_kwargs={'name':'_raw_data_remark'})

    def set_origin(self, raw: RawModel):
        self.raw_data_id = raw.id
        self.raw_data_params = raw.params
        self.raw_data_table = raw.__tablename__


class NoPKModel(RawDataOrigin):
    created_at: Optional[datetime] = Field(
        sa_column=Column(DateTime(), default=func.now())
    )
    updated_at: Optional[datetime] = Field(
        sa_column=Column(DateTime(), default=func.now(), onupdate=func.now())
    )


class ToolModel(ToolTable, NoPKModel):
    @declared_attr
    def __tablename__(cls) -> str:
        plugin_name = _get_plugin_name(cls)
        plural_entity = inflect_engine.plural_noun(cls.__name__.lower())
        return f'_tool_{plugin_name}_{plural_entity}'


class DomainModel(NoPKModel):
    id: str = Field(primary_key=True)


class ToolScope(ToolModel):
    id: str = Field(primary_key=True)
    name: str


class DomainScope(DomainModel):
    pass


def generate_domain_id(tool_model: ToolModel, connection_id: str):
    """
    Generate an identifier for a domain entity
    from the tool entity it originates from.
    """
    model_type = type(tool_model)
    segments = [_get_plugin_name(model_type), model_type.__name__, str(connection_id)]
    mapper = inspect(model_type)
    for primary_key_column in mapper.primary_key:
        prop = mapper.get_property_by_column(primary_key_column)
        attr_val = getattr(tool_model, prop.key)
        segments.append(str(attr_val))
    return ':'.join(segments)


def _get_plugin_name(cls):
    """
    Get the plugin name from a class by looking into
    the file path of its module.
    """
    module = getmodule(cls)
    path_segments = module.__file__.split(os.sep)
    # Finds the name of the first enclosing folder
    # that is not a python module
    depth = len(module.__name__.split('.')) + 1
    return path_segments[-depth]
