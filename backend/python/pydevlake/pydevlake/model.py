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
from sqlalchemy import Column, String, DateTime, func
from sqlalchemy.orm import declared_attr
from sqlalchemy.inspection import inspect
from sqlmodel import SQLModel, Field
import inflect

inflect_engine = inflect.engine()


def get_plugin_name(cls):
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
    raw_data_params: str = Field(sa_column_kwargs={'name':'_raw_data_params'})
    raw_data_table: str = Field(sa_column_kwargs={'name':'_raw_data_table'})
    raw_data_id: Optional[str] = Field(sa_column_kwargs={'name':'_raw_data_id'})
    raw_data_remark: Optional[str] = Field(sa_column_kwargs={'name':'_raw_data_remark'})

    def set_origin(self, raw: RawModel):
        self.raw_data_id = raw.id
        self.raw_data_params = raw.params
        self.raw_data_table = raw.__tablename__


class ToolModel(RawDataOrigin):
    @declared_attr
    def __tablename__(cls) -> str:
        plugin_name = get_plugin_name(cls)
        plural_entity = inflect_engine.plural_noun(cls.__name__.lower())
        return f'_tool_{plugin_name}_{plural_entity}'


class NoPKModel(SQLModel):
    created_at: datetime = Field(
        sa_column=Column(DateTime(), default=func.now())
    )
    updated_at: datetime = Field(
        sa_column=Column(DateTime(), default=func.now(), onupdate=func.now())
    )


class DomainModel(NoPKModel):
    id: str = Field(primary_key=True)


def generate_domain_id(tool_model: ToolModel, connection_id: str):
    """
    Generate an identifier for a domain entity
    from the tool entity it originates from.
    """
    model_type = type(tool_model)
    segments = [get_plugin_name(model_type), model_type.__name__, str(connection_id)]
    mapper = inspect(model_type)
    for primary_key_column in mapper.primary_key:
        prop = mapper.get_property_by_column(primary_key_column)
        attr_val = getattr(tool_model, prop.key)
        segments.append(str(attr_val))
    return ':'.join(segments)
