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

import { useState } from 'react';
import { useDebounce } from 'ahooks';

import { MultiSelector } from '@/components';
import { useRefreshData } from '@/hooks';

import type { ItemType } from './types';
import * as API from './api';

interface Props {
  plugin: string;
  connectionId: ID;
  disabledItems?: any[];
  selectedItems?: any[];
  onChangeItems?: (selectedItems: any[]) => void;
}

export const DataScopeSearch = ({ plugin, connectionId, disabledItems, selectedItems, onChangeItems }: Props) => {
  const [query, setQuery] = useState('');

  const search = useDebounce(query, { wait: 500 });

  const { ready, data } = useRefreshData<{ children: ItemType[] }>(async () => {
    if (!search) return [];
    return API.searchScope(plugin, connectionId, {
      search,
      page: 1,
      pageSize: 50,
    });
  }, [search]);

  const getKey = (it: ItemType) => it.id;

  const getName = (it: ItemType) => it.fullName;

  const handleChangeItems = (selectedItems: ItemType[]) => onChangeItems?.(selectedItems);

  return (
    <MultiSelector
      placeholder="Search Repositories..."
      items={data?.children ?? []}
      getKey={getKey}
      getName={getName}
      disabledItems={disabledItems}
      selectedItems={selectedItems}
      onChangeItems={handleChangeItems}
      loading={!ready}
      noResult="No Repositories Available."
      onQueryChange={(query) => setQuery(query)}
    />
  );
};
