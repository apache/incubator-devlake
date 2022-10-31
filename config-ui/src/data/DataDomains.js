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

/**
 * @type {Record<string, string>}
 */
const DataDomainTypes = {
  CODE: 'CODE',
  TICKET: 'TICKET',
  CODE_REVIEW: 'CODEREVIEW',
  CROSSDOMAIN: 'CROSS',
  DEVOPS: 'CICD'
  // USER: 'user',
}

const DataDomains = [
  DataDomainTypes.CODE,
  DataDomainTypes.TICKET,
  DataDomainTypes.CODE_REVIEW,
  DataDomainTypes.CROSSDOMAIN,
  DataDomainTypes.DEVOPS
  // ScopeEntityTypes.USER,
]

/**
 * @type {Array<Object>}
 */
const ALL_DATA_DOMAINS = [
  {
    id: 1,
    name: 'source-code-management',
    title: 'Source Code Management',
    value: DataDomainTypes.CODE
  },
  {
    id: 2,
    name: 'issue-tracking',
    title: 'Issue Tracking',
    value: DataDomainTypes.TICKET
  },
  {
    id: 3,
    name: 'code-review',
    title: 'Code Review',
    value: DataDomainTypes.CODE_REVIEW
  },
  {
    id: 4,
    name: 'cross-domain',
    title: 'Crossdomain',
    value: DataDomainTypes.CROSSDOMAIN
  },
  { id: 5, name: 'ci-cd', title: 'CI/CD', value: DataDomainTypes.DEVOPS }
]

export { DataDomainTypes, DataDomains, ALL_DATA_DOMAINS }
