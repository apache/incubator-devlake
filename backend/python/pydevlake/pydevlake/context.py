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


from urllib.parse import urlparse, parse_qsl
from sqlmodel import SQLModel, create_engine

from pydevlake.model import Connection, TransformationRule


class Context:
    def __init__(self,
                 db_url: str,
                 scope_id: str,
                 connection: Connection,
                 transformation_rule: TransformationRule = None,
                 options: dict = None):
        self.db_url = db_url
        self.scope_id = scope_id
        self.connection = connection
        self.transformation_rule = transformation_rule
        self.options = options or {}
        self._engine = None

    @property
    def engine(self):
        if not self._engine:
            db_url = self.db_url
            if not db_url:
                raise Exception("Missing db_url setting")

            # `parseTime` parameter is not understood by MySQL driver
            # so we have to parse query args to remove it
            connect_args = dict(parse_qsl(urlparse(self.db_url).query))
            db_url = self.db_url.split('?')[0]

            if 'parseTime' in connect_args:
                del connect_args['parseTime']

            self._engine = create_engine(db_url, connect_args=connect_args)
            SQLModel.metadata.create_all(self._engine)
        return self._engine

    @property
    def incremental(self):
        return self.options.get('incremental') is True
