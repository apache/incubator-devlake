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

import customize from './customize';
import dbt from './dbt';
import feishu from './feishu';
import general from './general';
import gitextractor from './gitextractor';
import github from './github';
import gitlab from './gitlab';
import jenkins from './jenkins';
import jira from './jira';
import refdiff from './refdiff';
import starrocks from './starrocks';
import tapd from './tapd';
import zentao from './zentao';

export const EXAMPLE_CONFIG = [
  {
    id: 'general',
    name: 'Load General Configuration',
    config: general,
  },
  {
    id: 'refdiff',
    name: 'Load RefDiff Configuration',
    config: refdiff,
  },
  {
    id: 'gitextractor',
    name: 'Load GitExtractor Configuration',
    config: gitextractor,
  },
  {
    id: 'github',
    name: 'Load GitHub Configuration',
    config: github,
  },
  {
    id: 'gitlab',
    name: 'Load GitLab Configuration',
    config: gitlab,
  },
  {
    id: 'jira',
    name: 'Load Jira Configuration',
    config: jira,
  },
  {
    id: 'jenkins',
    name: 'Load Jenkins Configuration',
    config: jenkins,
  },
  {
    id: 'feishu',
    name: 'Load Feishu Configuration',
    config: feishu,
  },
  {
    id: 'dbt',
    name: 'Load DBT Configuration',
    config: dbt,
  },
  {
    id: 'tapd',
    name: 'Load TAPD Configuration',
    config: tapd,
  },
  {
    id: 'zentao',
    name: 'Load ZENTAO Configuration',
    config: zentao,
  },
  {
    id: 'starrocks',
    name: 'Load StarRocks Configuration',
    config: starrocks,
  },
  {
    id: 'customize',
    name: 'Load Customize Configuration',
    config: customize,
  },
];
