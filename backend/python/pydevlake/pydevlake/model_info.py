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

import jsonref
from pydantic import BaseModel


class DynamicModelInfo(BaseModel):
    json_schema: dict
    table_name: str

    class Config:
        allow_population_by_field_name = True

    @staticmethod
    def from_model(model_class):
        schema = model_class.schema(by_alias=True)
        if 'definitions' in schema:
            # Replace $ref with actual schema
            schema = jsonref.replace_refs(schema, proxies=False)
            del schema['definitions']
        # Pydantic forgets to put type in enums
        for prop in schema['properties'].values():
            if 'type' not in prop and 'enum' in prop:
                prop['type'] = 'string'
        return DynamicModelInfo(
            json_schema=schema,
            table_name=model_class.__tablename__
        )
