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
 * @type {object}
 */
const DataEntityTypes = {
  CODE: 'CODE',
  TICKET: 'TICKET',
  CODE_REVIEW: 'CODEREVIEW',
  CROSSDOMAIN: 'CROSS',
  DEVOPS: 'CICD'
  // USER: 'user',
}

/**
 * @type {<Array<Object>>}
 */
const DataEntityList = [
  {
    id: 1,
    name: 'source-code-management',
    title: 'Source Code Management',
    value: DataEntityTypes.CODE
  },
  {
    id: 2,
    name: 'issue-tracking',
    title: 'Issue Tracking',
    value: DataEntityTypes.TICKET
  },
  {
    id: 3,
    name: 'code-review',
    title: 'Code Review',
    value: DataEntityTypes.CODE_REVIEW
  },
  {
    id: 4,
    name: 'cross-domain',
    title: 'Crossdomain',
    value: DataEntityTypes.CROSSDOMAIN
  },
  {
    id: 5,
    name: 'ci-cd',
    title: 'CI/CD',
    value: DataEntityTypes.DEVOPS
  }
]

/**
 * @typedef {object} DataEntity
 * @property {number} id
 * @property {string} name
 * @property {string?} title
 * @property {string} value
 * @property {'CODE'|'TICKET'|'CODEREVIEW'|'CROSS'|'CICD'} type
 */
class DataEntity {
  constructor(data = {}) {
    this.id = data?.type
      ? DataEntityList.find((e) => e.value === data?.type)?.id
      : 0
    this.name = data?.type
      ? DataEntityList.find((e) => e.value === data?.type)?.name
      : DataEntityList[0]?.name
    this.title = data?.type
      ? DataEntityList.find((e) => e.value === data?.type)?.title
      : DataEntityList[0]?.title
    this.value = data?.type
      ? DataEntityList.find((e) => e.value === data?.type)?.value
      : DataEntityList[0]?.type
    this.type = data?.type || DataEntityList[0]?.type
  }

  get(property) {
    return this[property]
  }

  set(property, value) {
    this[property] = value
    return this.property
  }
}

export default DataEntity
