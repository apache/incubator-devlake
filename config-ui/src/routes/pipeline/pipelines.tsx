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
import { useState, useMemo } from 'react';

import { PageHeader } from '@/components';
import { useRefreshData } from '@/hooks';

import { PipelineTable } from './components';
import * as API from './api';

export const Pipelines = () => {
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);

  const { ready, data } = useRefreshData(() => API.getPipelines());

  const [dataSource, total] = useMemo(() => [(data?.pipelines ?? []).map((it) => it), data?.count ?? 0], [data]);

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Advanced', path: '/blueprints' },
        { name: 'Pipelines', path: '/pipelines' },
      ]}
    >
      <PipelineTable
        loading={!ready}
        dataSource={dataSource}
        pagination={{
          total,
          page,
          pageSize,
          onChange: setPage,
        }}
        noData={{
          text: 'Add new projects to see engineering metrics based on projects.',
        }}
      />
    </PageHeader>
  );
};
