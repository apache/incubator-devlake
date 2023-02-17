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

import React, { useState, useEffect } from 'react';

import { MultiSelector } from '@/components';

import type { ItemType } from './types';
import * as API from './api';

interface Props {
  plugin: string;
  connectionId: ID;
  selectedItems?: any[];
  onChangeItems?: (selectedItems: any[]) => void;
}

export const DataScopeSearch = ({ plugin, connectionId, selectedItems, onChangeItems }: Props) => {
  const [loading, setLoading] = useState(false);
  const [items, setItems] = useState<ItemType[]>([]);
  const [search, setSearch] = useState('');

  useEffect(() => {
    if (!search) return;
    setItems([]);
    setLoading(true);

    const timer = setTimeout(async () => {
      try {
        const res = await API.searchScope(plugin, connectionId, {
          search,
          page: 1,
          pageSize: 50,
        });
        setItems(res.children ?? []);
      } finally {
        setLoading(false);
      }
    }, 1000);

    return () => clearTimeout(timer);
  }, [search]);

  const getKey = (it: ItemType) => it.id;

  const getName = (it: ItemType) => it.name;

  const handleChangeItems = (selectedItems: ItemType[]) => onChangeItems?.(selectedItems.map((it) => it.data));

  return (
    <MultiSelector
      placeholder="Search Repositories..."
      items={items}
      getKey={getKey}
      getName={getName}
      selectedItems={selectedItems}
      onChangeItems={handleChangeItems}
      loading={loading}
      noResult="No Repositories Available."
      onQueryChange={(s) => setSearch(s)}
    />
  );
};
