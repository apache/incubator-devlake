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

from typing import Optional
from pydevlake import ToolModel, Field
from pydevlake.extractor import autoextract
import pytest


def test_autoextract_from_json_pointer():
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
    assert model.name == 'test'
    assert model.version == '1.0.0'


def test_autoextract_optional_field_with_missing_value():
    class DummyModel(ToolModel):
        name: str
        version: Optional[str] = Field(source='/version/number')

    json = {
        'name': 'test',
        'version': {
            # missing 'number'
            'build_date': '2023-04-19'
        }
    }
    model = autoextract(json, DummyModel)
    assert model.name == 'test'
    assert model.version == None


def test_autoextract_optional_field_with_default_with_missing_value():
    class DummyModel(ToolModel):
        name: str
        version: Optional[str] = Field(source='/version/number', default='1.0.0')

    json = {
        'name': 'test',
        'version': {
            # missing 'number'
            'build_date': '2023-04-19'
        }
    }
    model = autoextract(json, DummyModel)
    assert model.name == 'test'
    assert model.version == '1.0.0'



def test_autoextract_required_field_and_missing_field():
    class DummyModel(ToolModel):
        name: str
        version: str = Field(source='/version/number')

    json = {
        'name': 'test',
        'version': {
            # missing 'number'
            'build_date': '2023-04-19'
        }
    }

    pytest.raises(ValueError, autoextract, json, DummyModel)
