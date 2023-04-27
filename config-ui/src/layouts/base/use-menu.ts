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

import { useMemo } from 'react';
import { IconName } from '@blueprintjs/core';

import { PluginConfig, PluginType } from '@/plugins';

export type MenuItemType = {
  key: string;
  title: string;
  icon?: IconName;
  iconUrl?: string;
  path: string;
  children?: MenuItemType[];
  target?: boolean;
  isBeta?: boolean;
  disabled?: boolean;
};

export const useMenu = () => {
  return useMemo(
    () =>
      [
        {
          key: 'connection',
          title: 'Connections',
          icon: 'data-connection',
          path: '/connections',
          children: PluginConfig.filter((p) =>
            [PluginType.Connection, PluginType.Incoming_Connection].includes(p.type),
          ).map((it) => ({
            key: it.plugin,
            title: it.name,
            iconUrl: it.icon,
            path: `/connections/${it.plugin}`,
            isBeta: it.isBeta,
          })),
        },
        {
          key: 'project',
          title: 'Projects',
          icon: 'home',
          path: '/projects',
        },
        {
          key: 'advanced',
          title: 'Advanced',
          icon: 'pulse',
          // path: '/advanced',
          children: [
            {
              key: 'blueprints',
              title: 'Blueprints',
              icon: '',
              path: '/blueprints',
            },
            {
              key: 'pipelines',
              title: 'Pipelines',
              icon: '',
              path: '/pipelines',
              disabled: true,
            },
          ],
        },
      ] as MenuItemType[],
    [],
  );
};
