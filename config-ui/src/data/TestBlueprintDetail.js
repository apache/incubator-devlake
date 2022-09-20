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
import React from 'react'
import { Intent, Icon, Colors, Spinner } from '@blueprintjs/core'

import { NullBlueprint } from '@/data/NullBlueprint'
import { Providers, ProviderIcons } from '@/data/Providers'
import {
  StageStatus,
  TaskStatus,
  TaskStatusLabels,
  StatusColors,
  StatusBgColors
} from '@/data/Task'

const EMPTY_RUN = {
  id: null,
  status: TaskStatus.CREATED,
  statusLabel: TaskStatusLabels[TaskStatus.RUNNING],
  icon: null,
  startedAt: Date.now(),
  duration: '0 min',
  stage: 'Stage 1',
  tasksTotal: 0,
  tasksFinished: 0,
  error: null
}

const TEST_BLUEPRINT = {
  ...NullBlueprint,
  id: 1,
  name: 'DevLake Daily Blueprint',
  createdAt: new Date().toLocaleString(),
  updatedAt: new Date().toLocaleString()
}

const TEST_CONNECTIONS = [
  {
    id: 0,
    provider: Providers.GITHUB,
    name: 'Merico GitHub',
    dataScope: 'merico-dev/ake, merico-dev/lake-website',
    dataEntities: ['code', 'ticket', 'user']
  },
  {
    id: 0,
    provider: Providers.JIRA,
    name: 'Merico JIRA',
    dataScope: 'Sprint Dev Board, DevLake Sync Board ',
    dataEntities: ['ticket']
  }
]

// eslint-disable-next-line no-unused-vars
const TEST_BLUEPRINT_API_RESPONSE = {
  name: 'DEVLAKE (Hourly)',
  mode: 'NORMAL',
  plan: [
    [
      {
        plugin: 'github',
        subtasks: [
          'collectApiRepo',
          'extractApiRepo',
          'collectApiIssues',
          'extractApiIssues',
          'collectApiPullRequests',
          'extractApiPullRequests',
          'collectApiComments',
          'extractApiComments',
          'collectApiEvents',
          'extractApiEvents',
          'collectApiPullRequestCommits',
          'extractApiPullRequestCommits',
          'collectApiPullRequestReviews',
          'extractApiPullRequestReviewers',
          'collectApiCommits',
          'extractApiCommits',
          'collectApiCommitStats',
          'extractApiCommitStats',
          'enrichPullRequestIssues',
          'convertRepo',
          'convertIssues',
          'convertCommits',
          'convertIssueLabels',
          'convertPullRequestCommits',
          'convertPullRequests',
          'convertPullRequestLabels',
          'convertPullRequestIssues',
          'convertIssueComments',
          'convertPullRequestComments'
        ],
        options: {
          connectionId: 1,
          owner: 'e2corporation',
          repo: 'incubator-devlake',
          transformationRules: {
            issueComponent: '',
            issuePriority: '',
            issueSeverity: '',
            issueTypeBug: '',
            issueTypeIncident: '',
            issueTypeRequirement: '',
            prComponent: '',
            prType: ''
          }
        }
      },
      {
        plugin: 'gitextractor',
        subtasks: null,
        options: {
          repoId: 'github:GithubRepo:1:506830252',
          url: 'https://git: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@github.com/e2corporation/incubator-devlake.git'
        }
      }
    ],
    [
      {
        plugin: 'refdiff',
        subtasks: null,
        options: {
          tagsLimit: 10,
          tagsOrder: '',
          tagsPattern: ''
        }
      }
    ]
  ],
  enable: true,
  cronConfig: '0 0 * * *',
  isManual: false,
  settings: {
    version: '1.0.0',
    connections: [
      {
        connectionId: 1,
        plugin: 'github',
        scope: [
          {
            entities: ['CODE', 'TICKET'],
            options: {
              owner: 'e2corporation',
              repo: 'incubator-devlake'
            },
            transformation: {
              prType: '',
              prComponent: '',
              issueSeverity: '',
              issueComponent: '',
              issuePriority: '',
              issueTypeRequirement: '',
              issueTypeBug: '',
              issueTypeIncident: '',
              refdiff: {
                tagsOrder: '',
                tagsPattern: '',
                tagsLimit: 10
              }
            }
          }
        ]
      }
    ]
  },
  id: 1,
  createdAt: '2022-07-11T10:23:38.908-04:00',
  updatedAt: '2022-07-11T10:23:38.908-04:00'
}

const TEST_STAGES = [
  {
    id: 1,
    name: 'stage-1',
    title: 'Stage 1',
    status: StageStatus.COMPLETED,
    icon: <Icon icon='tick-circle' size={14} color={StatusColors.COMPLETE} />,
    tasks: [
      {
        id: 0,
        provider: 'jira',
        icon: ProviderIcons[Providers.JIRA](14, 14),
        title: 'JIRA',
        caption: 'STREAM Board',
        duration: '4 min',
        subTasksCompleted: 25,
        recordsFinished: 1234,
        message: 'All 25 subtasks completed',
        status: TaskStatus.COMPLETE
      },
      {
        id: 0,
        provider: 'jira',
        icon: ProviderIcons[Providers.JIRA](14, 14),
        title: 'JIRA',
        caption: 'LAKE Board',
        duration: '4 min',
        subTasksCompleted: 25,
        recordsFinished: 1234,
        message: 'All 25 subtasks completed',
        status: TaskStatus.COMPLETE
      }
    ],
    stageHeaderClassName: 'complete'
  },
  {
    id: 2,
    name: 'stage-2',
    title: 'Stage 2',
    status: StageStatus.PENDING,
    icon: <Spinner size={14} intent={Intent.PRIMARY} />,
    tasks: [
      {
        id: 0,
        provider: 'jira',
        icon: ProviderIcons[Providers.JIRA](14, 14),
        title: 'JIRA',
        caption: 'EE Board',
        duration: '5 min',
        subTasksCompleted: 25,
        recordsFinished: 1234,
        message: 'Subtask 5/25: Extracting Issues',
        status: TaskStatus.ACTIVE
      },
      {
        id: 0,
        provider: 'jira',
        icon: ProviderIcons[Providers.JIRA](14, 14),
        title: 'JIRA',
        caption: 'EE Bugs Board',
        duration: '0 min',
        subTasksCompleted: 0,
        recordsFinished: 0,
        message: 'Invalid Board ID',
        status: TaskStatus.FAILED
      }
    ],
    stageHeaderClassName: 'active'
  },
  {
    id: 3,
    name: 'stage-3',
    title: 'Stage 3',
    status: StageStatus.PENDING,
    icon: null,
    tasks: [
      {
        id: 0,
        provider: 'github',
        icon: ProviderIcons[Providers.GITHUB](14, 14),
        title: 'GITHUB',
        caption: 'merico-dev/lake',
        duration: null,
        subTasksCompleted: 0,
        recordsFinished: 0,
        message: 'Subtasks pending',
        status: TaskStatus.CREATED
      }
    ],
    stageHeaderClassName: 'pending'
  },
  {
    id: 4,
    name: 'stage-4',
    title: 'Stage 4',
    status: StageStatus.PENDING,
    icon: null,
    tasks: [
      {
        id: 0,
        providr: 'github',
        icon: ProviderIcons[Providers.GITHUB](14, 14),
        title: 'GITHUB',
        caption: 'merico-dev/lake',
        duration: null,
        subTasksCompleted: 0,
        recordsFinished: 0,
        message: 'Subtasks pending',
        status: TaskStatus.CREATED
      }
    ],
    stageHeaderClassName: 'pending'
  }
]

const TEST_HISTORICAL_RUNS = [
  {
    id: 0,
    status: 'TASK_COMPLETED',
    statusLabel: 'Completed',
    statusIcon: <Icon icon='tick-circle' size={14} color={Colors.GREEN5} />,
    startedAt: '05/25/2022 0:00 AM',
    completedAt: '05/25/2022 0:15 AM',
    duration: '15 min'
  },
  {
    id: 1,
    status: 'TASK_COMPLETED',
    statusLabel: 'Completed',
    statusIcon: <Icon icon='tick-circle' size={14} color={Colors.GREEN5} />,
    startedAt: '05/25/2022 0:00 AM',
    completedAt: '05/25/2022 0:15 AM',
    duration: '15 min'
  },
  {
    id: 2,
    status: 'TASK_FAILED',
    statusLabel: 'Failed',
    statusIcon: <Icon icon='delete' size={14} color={Colors.RED5} />,
    startedAt: '05/25/2022 0:00 AM',
    completedAt: '05/25/2022 0:00 AM',
    duration: '0 min'
  },
  {
    id: 3,
    status: 'TASK_COMPLETED',
    statusLabel: 'Completed',
    statusIcon: <Icon icon='tick-circle' size={14} color={Colors.GREEN5} />,
    startedAt: '05/25/2022 0:00 AM',
    completedAt: '05/25/2022 0:15 AM',
    duration: '15 min'
  },
  {
    id: 4,
    status: 'TASK_COMPLETED',
    statusLabel: 'Completed',
    statusIcon: <Icon icon='tick-circle' size={14} color={Colors.GREEN5} />,
    startedAt: '05/25/2022 0:00 AM',
    completedAt: '05/25/2022 0:15 AM',
    duration: '15 min'
  },
  {
    id: 5,
    status: 'TASK_FAILED',
    statusLabel: 'Failed',
    statusIcon: <Icon icon='delete' size={14} color={Colors.RED5} />,
    startedAt: '05/25/2022 0:00 AM',
    completedAt: '05/25/2022 0:00 AM',
    duration: '0 min'
  }
]

const TEST_RUN = {
  id: null,
  status: TaskStatus.RUNNING,
  statusLabel: TaskStatusLabels[TaskStatus.RUNNING],
  icon: <Spinner size={18} intent={Intent.PRIMARY} />,
  startedAt: '7/7/2022, 5:31:33 PM',
  duration: '1 min',
  stage: 'Stage 1',
  tasksTotal: 5,
  tasksFinished: 8,
  // totalTasks: 13,
  error: null
}

export {
  EMPTY_RUN,
  TEST_RUN,
  TEST_BLUEPRINT,
  TEST_CONNECTIONS,
  TEST_HISTORICAL_RUNS,
  TEST_BLUEPRINT_API_RESPONSE,
  TEST_STAGES
}
