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
from pathlib import Path
from string import Template
import json

from pydevlake.message import Connection, TransformationRule


# TODO: Move swagger documentation generation to GO side along with API implementation
TEMPLATE_PATH = str(Path(__file__).parent / 'doc.template.json')

def generate_doc(plugin_name: str, 
                 connection_type: Type[Connection], 
                 transformation_rule_type: Type[TransformationRule]):
    with open(TEMPLATE_PATH, 'r') as f:
        doc_template = Template(f.read())
        connection_schema = connection_type.schema_json()
        transformation_rule_schema = transformation_rule_type.schema_json()
        doc = doc_template.substitute(
            plugin_name=plugin_name, 
            connection_schema=connection_schema,
            transformation_rule_schema=transformation_rule_schema)
        return json.loads(doc)
