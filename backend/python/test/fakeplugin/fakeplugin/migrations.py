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

from fakeplugin.models import FakeConnection, FakePipeline, FakeProject, FakeScopeConfig
from pydevlake.migration import migration, MigrationScriptBuilder, Dialect


@migration(20230530000001, name="init schemas")
def init_schemas(b: MigrationScriptBuilder):
    b.create_tables(FakeConnection, FakePipeline, FakeProject, FakeScopeConfig)


# test migration
@migration(20230630000001, name="populated _raw_data_table column for fakeproject")
def add_raw_data_params_table_to_scope(b: MigrationScriptBuilder):
    b.execute(f'UPDATE {FakeProject.__tablename__} SET _raw_data_table = "_raw_fakeproject_scopes" WHERE 1=1', Dialect.MYSQL) #mysql only
    b.execute(f'''UPDATE {FakeProject.__tablename__} SET _raw_data_table = '_raw_fakeproject_scopes' WHERE 1=1''', Dialect.POSTGRESQL) #mysql and postgres