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


from sqlalchemy.engine import Engine

from pydevlake.model import Connection, TransformationRule, ToolScope


class Context:
    def __init__(self,
                 engine: Engine,
                 scope: ToolScope,
                 connection: Connection,
                 transformation_rule: TransformationRule = None,
                 options: dict = None):
        self.engine = engine
        self.scope = scope
        self.connection = connection
        self.transformation_rule = transformation_rule
        self.options = options or {}
        self._engine = None

    @property
    def incremental(self) -> bool:
        return self.options.get('incremental') is True
