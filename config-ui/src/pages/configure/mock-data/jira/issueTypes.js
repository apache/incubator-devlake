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
const issueTypesData = [
  {
    self: 'https://your-domain.atlassian.net/rest/api/3/issueType/1',
    id: 1,
    description: 'A small, distinct piece of work.',
    iconUrl: null,
    name: 'Task',
    title: 'Task',
    value: 'Task',
    subtask: false,
    avatarId: 1,
    hierarchyLevel: 0
  },
  {
    self: 'https://your-domain.atlassian.net/rest/api/3/issueType/2',
    id: 2,
    description: 'A problem or error.',
    iconUrl: null,
    name: 'Bug',
    title: 'Bug',
    value: 'Bug',
    subtask: false,
    avatarId: 10002,
    entityId: '9d7dd6f7-e8b6-4247-954b-7b2c9b2a5ba2',
    hierarchyLevel: 0,
    scope: {
      type: 'PROJECT',
      project: {
        id: 10000,
        key: 'KEY',
        name: 'Next Gen Project'
      }
    }
  },
  {
    self: 'https://your-domain.atlassian.net/rest/api/3/issueType/3',
    id: 3,
    description: 'A product requirement or feature',
    iconUrl: null,
    name: 'Requirement',
    title: 'Requirement',
    value: 'Requirement',
    subtask: false,
    avatarId: 1,
    hierarchyLevel: 0
  },
  {
    self: 'https://your-domain.atlassian.net/rest/api/3/issueType/4',
    id: 4,
    description: 'A big user story that needs to be broken down',
    iconUrl: null,
    name: 'Epic',
    title: 'Epic',
    value: 'Epic',
    subtask: false,
    avatarId: 1,
    hierarchyLevel: 0
  },
  {
    self: 'https://your-domain.atlassian.net/rest/api/3/issueType/5',
    id: 5,
    description: 'A sub-task/child-task',
    iconUrl: null,
    name: 'Sub-task',
    title: 'Sub-task',
    value: 'Sub-task',
    subtask: false,
    avatarId: 1,
    hierarchyLevel: 0
  },
  {
    self: 'https://your-domain.atlassian.net/rest/api/3/issueType/5',
    id: 6,
    description: 'A P0 Incident/Event',
    iconUrl: null,
    name: 'P0',
    title: 'P0',
    value: 'P0',
    subtask: false,
    avatarId: 1,
    hierarchyLevel: 0
  },
  {
    self: 'https://your-domain.atlassian.net/rest/api/3/issueType/5',
    id: 7,
    description: 'A P1 Incident/Event',
    iconUrl: null,
    name: 'P1',
    title: 'P1',
    value: 'P1',
    subtask: false,
    avatarId: 1,
    hierarchyLevel: 0
  },
  {
    self: 'https://your-domain.atlassian.net/rest/api/3/issueType/5',
    id: 8,
    description: 'A P2 Incident/Event',
    iconUrl: null,
    name: 'P2',
    title: 'P2',
    value: 'P2',
    subtask: false,
    avatarId: 1,
    hierarchyLevel: 0
  },
]

export {
  issueTypesData
}
