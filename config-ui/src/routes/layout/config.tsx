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

import {
  AppstoreOutlined,
  ProjectOutlined,
  ExperimentOutlined,
  KeyOutlined,
  DashboardOutlined,
  FileSearchOutlined,
  ApiOutlined,
  GithubOutlined,
  SlackOutlined,
} from '@ant-design/icons';

import { DOC_URL } from '@/release';

const PATH_PREFIX = import.meta.env.DEVLAKE_PATH_PREFIX ?? '';

type MenuItem = {
  key: string;
  label: string;
  icon?: React.ReactNode;
  children?: MenuItem[];
};

export const menuItems: MenuItem[] = [
  {
    key: `${PATH_PREFIX}/projects`,
    label: 'Projects',
    icon: <ProjectOutlined />,
  },
  {
    key: `${PATH_PREFIX}/connections`,
    label: 'Connections',
    icon: <AppstoreOutlined />,
  },
  {
    key: `${PATH_PREFIX}/advanced`,
    label: 'Advanced',
    icon: <ExperimentOutlined />,
    children: [
      {
        key: `${PATH_PREFIX}/advanced/blueprints`,
        label: 'Blueprints',
      },
      {
        key: `${PATH_PREFIX}/advanced/pipelines`,
        label: 'Pipelines',
      },
    ],
  },
  {
    key: `${PATH_PREFIX}/keys`,
    label: 'API Keys',
    icon: <KeyOutlined />,
  },
];

const getMenuMatchs = (items: MenuItem[], parentKey?: string) => {
  return items.reduce((pre, cur) => {
    pre[cur.key] = {
      ...cur,
      parentKey,
    };

    if (cur.children) {
      pre = { ...pre, ...getMenuMatchs(cur.children, cur.key) };
    }

    return pre;
  }, {} as Record<string, MenuItem & { parentKey?: string }>);
};

export const menuItemsMatch = getMenuMatchs(menuItems);

export const headerItems = [
  {
    link: import.meta.env.DEV ? `${window.location.protocol}//${window.location.hostname}:3002` : `/grafana`,
    label: 'Dashboards',
    icon: <DashboardOutlined />,
  },
  {
    link: DOC_URL.TUTORIAL,
    label: 'Docs',
    icon: <FileSearchOutlined />,
  },
  {
    link: '/api/swagger/index.html',
    label: 'API',
    icon: <ApiOutlined />,
  },
  {
    link: 'https://github.com/apache/incubator-devlake',
    label: 'GitHub',
    icon: <GithubOutlined />,
  },
  {
    link: 'https://join.slack.com/t/devlake-io/shared_invite/zt-26ulybksw-IDrJYuqY1FrdjlMMJhs53Q',
    label: 'Slack',
    icon: <SlackOutlined />,
  },
];
