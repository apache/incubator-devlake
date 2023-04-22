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


from typing import Type

from jsonpointer import resolve_pointer, JsonPointerException

from pydevlake import ToolModel


def autoextract(json: dict, model_cls: Type[ToolModel]) -> ToolModel:
    """
    Automatically extract a tool model from a json object.
    The tool model class can define fields with a source argument to specify the JSON pointer (RFC 6901) to the value.

    Example:
        class DummyModel(ToolModel):
            name: str
            version: str = Field(source='/version/number')

        json = {
            'name': 'test',
            'version': {
                'number': '1.0.0',
                'build_date': '2023-04-19'
            }
        }

        model = autoextract(json, DummyModel)
    """
    attributes = {}
    for field in model_cls.__fields__.values():
        pointer = field.field_info.extra.get('source')

        if pointer:
            if field.required:
                try:
                    value = resolve_pointer(json, pointer)
                except JsonPointerException:
                    raise ValueError(f"Missing required value for field {field.name} at {pointer}")
            else:
                value = resolve_pointer(json, pointer, field.default)
        else:
            value = json.get(field.name) or json.get(field.alias)
        attributes[field.name] = value
    return model_cls(**attributes)
