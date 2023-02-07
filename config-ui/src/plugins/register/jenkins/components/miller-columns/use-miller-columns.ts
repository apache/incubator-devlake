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

import { useState, useEffect, useMemo } from 'react';
import type { McsID, McsItem } from 'miller-columns-select';

import { useProxyPrefix } from '@/hooks';

import type { ScopeItemType } from '../../types';
import * as API from '../../api';

export type ExtraType = {
  type: 'folder' | 'file';
} & ScopeItemType;

export interface UseMillerColumnsProps {
  connectionId: ID;
}

export const useMillerColumns = ({ connectionId }: UseMillerColumnsProps) => {
  const [items, setItems] = useState<McsItem<ExtraType>[]>([]);
  const [loadedIds, setLoadedIds] = useState<ID[]>([]);

  const prefix = useProxyPrefix({
    plugin: 'jenkins',
    connectionId,
  });

  const formatJobs = (jobs: any, parentId: ID | null = null) =>
    jobs.map((it: any) => ({
      parentId,
      id: parentId ? `${parentId}/${it.name}` : it.name,
      title: it.name,
      type: it.jobs ? 'folder' : 'file',
    }));

  useEffect(() => {
    (async () => {
      const res = await API.getJobs(prefix);
      setItems(formatJobs(res.jobs));
      setLoadedIds(['root']);
    })();
  }, [prefix]);

  const onExpand = async (id: McsID) => {
    const res = await API.getJobChildJobs(prefix, (id as string).split('/').join('/job/'));

    setLoadedIds([...loadedIds, id]);
    setItems([...items, ...formatJobs(res.jobs, id)]);
  };

  return useMemo(
    () => ({
      items,
      getHasMore(id: McsID | null) {
        return !loadedIds.includes(id ?? 'root');
      },
      onExpand,
    }),
    [items, loadedIds],
  );
};
