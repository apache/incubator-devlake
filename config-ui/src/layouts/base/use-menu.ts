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

import { useMemo } from 'react'
import { IconName } from '@blueprintjs/core'

import { Plugins } from '@/registry'

export type MenuItemType = {
  key: string
  title: string
  icon?: IconName
  iconUrl?: string
  path: string
  children?: MenuItemType[]
  target?: boolean
}

export const useMenu = () => {
  const getGrafanaUrl = () => {
    const suffix = '/d/0Rjxknc7z/demo-homepage?orgId=1'
    const { protocol, hostname } = window.location

    return process.env.LOCAL
      ? `${protocol}//${hostname}:3002${suffix}`
      : `/grafana/${suffix}`
  }

  return useMemo(
    () =>
      [
        {
          key: 'connection',
          title: 'Connections',
          icon: 'data-connection',
          path: '/integrations',
          children: Plugins.filter((p) => p.type === 'integration').map(
            (it) => ({
              key: it.id,
              title: it.name,
              iconUrl: `/${it.icon}`,
              path: `/integrations/${it.id}`
            })
          )
        },
        {
          key: 'blueprint',
          title: 'Blueprints',
          icon: 'home',
          path: '/blueprints',
          children: [
            {
              key: 'create-blueprint',
              title: 'Create Blueprint',
              icon: 'git-pull',
              path: '/blueprints/create'
            }
          ]
        },
        {
          key: 'dashboard',
          title: 'Dashboard',
          icon: 'dashboard',
          path: getGrafanaUrl(),
          target: true
        }
      ] as MenuItemType[],
    []
  )
}
