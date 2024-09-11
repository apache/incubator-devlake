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

import { useState, useCallback } from 'react';
import { SearchOutlined } from '@ant-design/icons';
import { Space, Tag, Input } from 'antd';
import { MillerColumns } from '@mints/miller-columns';
import { useDebounce } from '@mints/hooks';

import API from '@/api';
import { Block, Loading } from '@/components';
import type { IPluginConfig } from '@/types';

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
  const [search, setSearch] = useState('');

  const searchDebounce = useDebounce(search, { wait: 500 });

  const request = useCallback(
    async (groupId?: string | number, params?: any) => {
      let data = [];
      let hasMore = false;
      let newParams = {};
      let originData = [];

      if (!searchDebounce) {
        const res = await API.scope.remote(plugin, connectionId, {
          groupId: groupId ?? null,
          pageToken: params?.pageToken,
        });

        data = res.children.map((it) => ({
          parentId: it.parentId,
          id: it.id,
          title: it.name ?? it.fullName,
          canExpand: it.type === 'group',
        }));

        hasMore = !!res.nextPageToken;
        newParams = {
          pageToken: res.nextPageToken,
        };
        originData = res.children;
      } else {
        const res = await API.scope.searchRemote(plugin, connectionId, {
          search: searchDebounce,
          page: params?.page ?? 1,
          pageSize: 20,
        });

        data = res.children.map((it) => ({
          parentId: it.parentId,
          id: it.id,
          title: it.fullName ?? it.name,
          canExpand: it.type === 'group',
        }));

        hasMore = res.children.length === res.pageSize;
        newParams = {
          page: (params?.page ?? 0) + 1,
          count: (params?.count ?? 0) + res.children.length,
        };
        originData = res.children;
      }

      return {
        data,
        hasMore,
        params: newParams,
        originData,
      };
    },
    [plugin, connectionId, searchDebounce],
  );

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
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <MillerColumns
          bordered
          theme={{
            colorPrimary: '#7497f7',
            borderColor: '#dbe4fd',
          }}
          request={request}
          columnCount={searchDebounce ? 1 : config.millerColumn?.columnCount ?? 1}
          columnHeight={300}
          mode={mode}
          renderTitle={(id?) =>
            !id &&
            config.millerColumn?.firstColumnTitle && (
              <S.ColumnTitle>{config.millerColumn.firstColumnTitle}</S.ColumnTitle>
            )
          }
          renderLoading={() => <Loading size={20} style={{ padding: '4px 12px' }} />}
          selectable
          disabledIds={disabledScope.map((it) => it.id)}
          selectedIds={selectedScope.map((it) => it.id)}
          onSelectedIds={(ids, data) => onChange((data ?? []).filter((it) => ids.includes(it.id)))}
        />
      </Block>
    </>
  );
};
