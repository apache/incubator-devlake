/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

const starrocks = [
  [
    {
      plugin: 'starrocks',
      options: {
        source_type: '', // mysql or postgres
        source_dsn: '', // gorm dsn
        update_column: '', // update column
        host: '127.0.0.1',
        port: 9030,
        user: 'root',
        password: '',
        database: 'lake',
        be_host: '',
        be_port: 8040,
        tables: ['_tool_.*'], // support regexp
        batch_size: 10000,
        order_by: {},
        extra: {}, // will append to create table sql
        domain_layer: '', // priority over tables
      },
    },
  ],
];

export default starrocks;
