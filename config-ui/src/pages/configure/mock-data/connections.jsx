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
import React, { useEffect, useState } from 'react'
const connectionsData = [
  {
    id: 0,
    name: 'Development Server',
    endpoint: 'https://jira-test-a345vf.merico.dev',
    status: 1,
    errors: []
  },
  {
    id: 1,
    name: 'Staging Server',
    endpoint: 'https://jira-staging-93xt5a.merico.dev',
    status: 2,
    errors: []
  },
  {
    id: 2,
    name: 'Production Server',
    endpoint: 'https://jira-prod-z51gox.merico.dev',
    status: 0,
    errors: []
  },
  {
    id: 3,
    name: 'Demo Instance 591',
    endpoint: 'https://jira-demo-591.merico.dev',
    status: 0,
    errors: []
  },
  {
    id: 4,
    name: 'Demo Instance 142',
    endpoint: 'https://jira-demo-142.merico.dev',
    status: 0,
    errors: []
  },
  {
    id: 5,
    name: 'Demo Instance 111',
    endpoint: 'https://jira-demo-111.merico.dev',
    status: 0,
    errors: []
  },
  {
    id: 6,
    name: 'Demo Instance 784',
    endpoint: 'https://jira-demo-784.merico.dev',
    status: 3,
    errors: []
  },
]

export {
  connectionsData
}
