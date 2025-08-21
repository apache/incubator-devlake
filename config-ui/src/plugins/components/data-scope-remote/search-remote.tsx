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
import { SearchOutlined } from '@ant-design/icons';
import { Space, Tag, Input, message } from 'antd';
import type { McsID, McsItem, McsColumn } from 'miller-columns-select';
import MillerColumnsSelect from 'miller-columns-select';
import { useDebounce } from 'ahooks';
import { uniqBy } from 'lodash';

import API from '@/api';
import { Loading, Block } from '@/components';
import { IPluginConfig } from '@/types';

import * as T from './types';
import * as S from './styled';

interface Props {
  mode: 'single' | 'multiple';
  plugin: string;
  connectionId: ID;
  config: IPluginConfig['dataScope'];
  disabledScope: any[];
  selectedScope: any[];
  onChange: (selectedScope: any[]) => void;
}

export const SearchRemote = ({ mode, plugin, connectionId, config, disabledScope, selectedScope, onChange }: Props) => {
  const [miller, setMiller] = useState<{
    items: McsItem<T.ResItem>[];
    loadedIds: ID[];
    errorId?: ID | null;
    nextTokenMap: Record<ID, string>;
  }>({
    items: [],
    loadedIds: [],
    nextTokenMap: {},
  });

  const [search, setSearch] = useState<{
    loading: boolean;
    items: McsItem<T.ResItem>[];
    currentItems: McsItem<T.ResItem>[];
    query: string;
    page: number;
    total: number;
  }>({
    loading: true,
    items: [],
    currentItems: [],
    query: '',
    page: 1,
    total: 0,
  });

  const searchDebounce = useDebounce(search.query, { wait: 500 });

  const allItems = useMemo(
    () =>
      uniqBy(
        [...miller.items, ...search.items].filter((it) => it.type === 'scope'),
        'id',
      ),
    [miller.items, search.items],
  );

  const getItems = async (groupId: ID | null, currentPageToken?: string) => {
    let newItems: McsItem<T.ResItem>[] = [];
    let nextPageToken = '';
    let errorId: ID | null;

    try {
      const res = await API.scope.remote(plugin, connectionId, {
        groupId,
        pageToken: currentPageToken,
      });

      newItems = (res.children ?? []).map((it) => ({
        ...it,
        title: it.name,
      }));

      nextPageToken = res.nextPageToken;
    } catch (err: any) {
      errorId = groupId;
      message.error(err.response.data.message);
    }

    if (nextPageToken && newItems.length) {
      setMiller((m) => ({
        ...m,
        items: [...m.items, ...newItems],
        nextTokenMap: {
          ...m.nextTokenMap,
          [`${groupId ? groupId : 'root'}`]: nextPageToken,
        },
      }));
    } else {
      setMiller((m) => ({
        ...m,
        items: [...m.items, ...newItems],
        loadedIds: [...m.loadedIds, groupId ?? 'root'],
        errorId,
      }));
    }
  };

  useEffect(() => {
    getItems(null);
  }, []);

  const searchItems = async () => {
    if (!searchDebounce) return;

    const res = await API.scope.searchRemote(plugin, connectionId, {
      search: searchDebounce,
      page: search.page,
      pageSize: 20,
    });

    const newItems = (res.children ?? []).map((it) => ({
      ...it,
      title: it.fullName ?? it.name,
    }));

    setSearch((s) => ({
      ...s,
      loading: false,
      items: [...allItems, ...newItems],
      currentItems: newItems,
      total: res.count,
    }));
  };

  useEffect(() => {
    searchItems();
  }, [searchDebounce, search.page]);

  return (
    <>
      <Block title={config.title} required>
        <Space wrap>
          {selectedScope.length ? (
            selectedScope.map((sc) => (
              <Tag
                key={sc.id}
                color="blue"
                closable
                onClose={() => onChange(selectedScope.filter((it) => it.id !== sc.id))}
              >
                {sc.fullName}
              </Tag>
            ))
          ) : (
            <span>Please select scope...</span>
          )}
        </Space>
      </Block>
      <Block>
        <Input
          prefix={<SearchOutlined />}
          placeholder={config.searchPlaceholder ?? 'Search'}
          value={search.query}
          onChange={(e) => setSearch({ ...search, query: e.target.value, loading: true, currentItems: [] })}
        />
        {!searchDebounce ? (
          <MillerColumnsSelect
            mode={mode}
            items={miller.items}
            columnCount={config.millerColumn?.columnCount ?? 1}
            columnHeight={300}
            getCanExpand={(it) => it.type === 'group'}
            getHasMore={(id) => !miller.loadedIds.includes(id ?? 'root')}
            getHasError={(id) => id === miller.errorId}
            onExpand={(id: McsID) => getItems(id, miller.nextTokenMap[id])}
            onScroll={(id: McsID | null) => getItems(id, miller.nextTokenMap[id ?? 'root'])}
            renderTitle={(column: McsColumn) =>
              !column.parentId &&
              config.millerColumn?.firstColumnTitle && (
                <S.ColumnTitle>{config.millerColumn.firstColumnTitle}</S.ColumnTitle>
              )
            }
            renderLoading={() => <Loading size={20} style={{ padding: '4px 12px' }} />}
            renderError={() => <span style={{ color: 'red' }}>Something Error</span>}
            disabledIds={(disabledScope ?? []).map((it) => it.id)}
            selectedIds={selectedScope.map((it) => it.id)}
            onSelectItemIds={(selectedIds: ID[]) => onChange(allItems.filter((it) => selectedIds.includes(it.id)))}
          />
        ) : (
          <MillerColumnsSelect
            mode={mode}
            items={search.currentItems}
            columnCount={1}
            columnHeight={300}
            getCanExpand={() => false}
            getHasMore={() => search.loading}
            onScroll={() => setSearch({ ...search, page: search.page + 1 })}
            renderLoading={() => <Loading size={20} style={{ padding: '4px 12px' }} />}
            disabledIds={(disabledScope ?? []).map((it) => it.id)}
            selectedIds={selectedScope.map((it) => it.id)}
            onSelectItemIds={(selectedIds: ID[]) => onChange(allItems.filter((it) => selectedIds.includes(it.id)))}
          />
        )}
      </Block>
    </>
  );
};
