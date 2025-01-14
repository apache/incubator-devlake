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

import { createBrowserRouter, Navigate } from 'react-router-dom';

import {
  DBMigrate,
  Onboard,
  Error,
  Layout,
  Connections,
  Connection,
  ProjectHomePage,
  ProjectLayout,
  ProjectGeneralSettings,
  ProjectWebhook,
  ProjectAdditionalSettings,
  BlueprintConnectionDetailPage,
  Pipelines,
  ApiKeys,
  NotFound,
} from '@/routes';

import { App } from '../App';

const PATH_PREFIX = import.meta.env.DEVLAKE_PATH_PREFIX ?? '/';

export const router = createBrowserRouter(
  [
    {
      path: 'db-migrate',
      element: <DBMigrate />,
    },
    {
      path: '/',
      element: <App />,
      errorElement: <Error />,
      children: [
        {
          index: true,
          element: <Navigate to="projects" />,
        },
        {
          path: 'onboard',
          element: <Onboard />,
        },
        {
          path: 'projects/:pname',
          element: <ProjectLayout />,
          children: [
            {
              index: true,
              element: <Navigate to="general-settings" />,
            },
            {
              path: 'general-settings',
              element: <ProjectGeneralSettings />,
            },
            {
              path: 'general-settings/:unique',
              element: <BlueprintConnectionDetailPage />,
            },
            {
              path: 'webhooks',
              element: <ProjectWebhook />,
            },
            {
              path: 'additional-settings',
              element: <ProjectAdditionalSettings />,
            },
          ],
        },
        {
          path: '',
          element: <Layout />,
          children: [
            {
              index: true,
              element: <Navigate to="projects" />,
            },
            {
              path: 'projects',
              element: <ProjectHomePage />,
            },
            {
              path: 'connections',
              element: <Connections />,
            },
            {
              path: 'connections/:plugin/:id',
              element: <Connection />,
            },
            {
              path: 'advanced',
              children: [
                {
                  path: 'keys',
                  element: <ApiKeys />,
                },
                {
                  path: 'pipelines',
                  element: <Pipelines />,
                },
              ],
            },
          ],
        },
      ],
    },
    {
      path: '*',
      element: <NotFound />,
    },
  ],
  {
    basename: PATH_PREFIX,
  },
);
