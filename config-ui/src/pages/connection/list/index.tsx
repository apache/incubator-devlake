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

import React, { useMemo } from 'react';
import { useParams } from 'react-router-dom';

import { PageHeader } from '@/components';
import { getPluginConfig } from '@/plugins';
import { WebHookConnection } from '@/plugins/register/webook';
import { ConnectionContextProvider } from '@/store';

import { Connection } from './connection';

export const ConnectionListPage = () => {
  const { plugin } = useParams<{ plugin: string }>();

  const config = useMemo(() => getPluginConfig(plugin), [plugin]);

  return (
    <ConnectionContextProvider plugin={plugin}>
      <PageHeader
        breadcrumbs={[
          { name: 'Connections', path: '/connections' },
          { name: config.name, path: `/connections/${plugin}` },
        ]}
      >
        {plugin === 'webhook' ? <WebHookConnection /> : <Connection plugin={plugin} />}
      </PageHeader>
    </ConnectionContextProvider>
  );
};
