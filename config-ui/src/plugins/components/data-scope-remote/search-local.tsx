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

import { useState, useReducer, useCallback } from 'react';
import { CheckCircleFilled, SearchOutlined } from '@ant-design/icons';
import { Space, Tag, Button, Input, Modal } from 'antd';
import { MillerColumns } from '@mints/miller-columns';
import { useDebounce } from '@mints/hooks';

import API from '@/api';
import { Block, Loading, Message } from '@/components';
import type { IPluginConfig } from '@/types';

import * as S from './styled';

type StateType = {
  status: string;
  scope: any[];
  originData: any[];
};

const reducer = (
  state: StateType,
  action: { type: string; payload?: Pick<Partial<StateType>, 'scope' | 'originData'> },
) => {
  switch (action.type) {
    case 'LOADING':
      return {
        ...state,
        status: 'loading',
      };
    case 'APPEND':
      return {
        ...state,
        scope: [...state.scope, ...(action.payload?.scope ?? [])],
        originData: [...state.originData, ...(action.payload?.originData ?? [])],
      };
    case 'DONE':
      return {
        ...state,
        status: 'done',
      };
    default:
      return state;
  }
};

interface Props {
  mode: 'single' | 'multiple';
  plugin: string;
  connectionId: ID;
  config: IPluginConfig['dataScope'];
  disabledScope: any[];
  selectedScope: any[];
  onChange: (selectedScope: any[]) => void;
}

export const SearchLocal = ({ mode, plugin, connectionId, config, disabledScope, selectedScope, onChange }: Props) => {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState('');

  const [{ status, scope, originData }, dispatch] = useReducer(reducer, {
    status: 'idle',
    scope: [],
    originData: [],
  });

  const searchDebounce = useDebounce(search, { wait: 500 });

  const request = useCallback(
    async (groupId?: string | number, params?: any) => {
      if (scope.length) {
        return {
          data: searchDebounce
            ? scope
                .filter((it) => it.title.includes(searchDebounce) && !it.canExpand)
                .map((it) => ({ ...it, parentId: null }))
            : scope.filter((it) => it.parentId === (groupId ?? null)),
          hasMore: status === 'loading' ? true : false,
          originData,
        };
      }

      const res = await API.scope.remote(plugin, connectionId, {
        groupId: groupId ?? null,
        pageToken: params?.nextPageToken,
      });

      const data = res.children.map((it) => ({
        parentId: it.parentId,
        id: it.id,
        title: it.name ?? it.fullName,
        canExpand: it.type === 'group',
      }));

      return {
        data,
        hasMore: !!res.nextPageToken,
        params: {
          nextPageToken: res.nextPageToken,
        },
        originData: res.children,
      };
    },
    [plugin, connectionId, scope, status, searchDebounce],
  );

  const handleRequestAll = async () => {
    setOpen(false);
    dispatch({ type: 'LOADING' });

    const getData = async (groupId?: string | number, currentPageToken?: string) => {
      const res = await API.scope.remote(plugin, connectionId, {
        groupId: groupId ?? null,
        pageToken: currentPageToken,
      });

      const data = res.children.map((it) => ({
        parentId: it.parentId,
        id: it.id,
        title: it.name ?? it.fullName,
        canExpand: it.type === 'group',
      }));

      dispatch({ type: 'APPEND', payload: { scope: data, originData: res.children } });

      if (res.nextPageToken) {
        await getData(groupId, res.nextPageToken);
      }

      await Promise.all(data.filter((it) => it.canExpand).map((it) => getData(it.id)));
    };

    await getData();

    dispatch({ type: 'DONE' });
  };

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
                {sc.fullName ?? sc.name}
              </Tag>
            ))
          ) : (
            <span>Please select scope...</span>
          )}
        </Space>
      </Block>
      <Block>
        {(status === 'loading' || status === 'cancel') && (
          <S.JobLoad>
            <Loading style={{ marginRight: 8 }} size={20} />
            Loading: <span className="count">{scope.length}</span> scopes found
          </S.JobLoad>
        )}

        {status === 'done' && (
          <S.JobLoad>
            <CheckCircleFilled style={{ color: '#4DB764' }} />
            <span className="count">{scope.length}</span> scopes found
          </S.JobLoad>
        )}

        {status === 'idle' && (
          <S.JobLoad>
            <Button type="primary" onClick={() => setOpen(true)}>
              Load all scopes to search by keywords
            </Button>
          </S.JobLoad>
        )}
      </Block>
      <Block>
        {status === 'done' && (
          <Input prefix={<SearchOutlined />} value={search} onChange={(e) => setSearch(e.target.value)} />
        )}
        <MillerColumns
          bordered
          theme={{
            colorPrimary: '#7497f7',
            borderColor: '#dbe4fd',
          }}
          request={request}
          columnCount={search ? 1 : config.millerColumn?.columnCount ?? 1}
          columnHeight={300}
          mode={mode}
          renderTitle={(id) =>
            !id &&
            config.millerColumn?.firstColumnTitle && (
              <S.ColumnTitle>{config.millerColumn.firstColumnTitle}</S.ColumnTitle>
            )
          }
          renderLoading={() => <Loading size={20} style={{ padding: '4px 12px' }} />}
          renderError={() => <span style={{ color: 'red' }}>Something Error</span>}
          selectable
          disabledIds={(disabledScope ?? []).map((it) => it.id)}
          selectedIds={selectedScope.map((it) => it.id)}
          onSelectedIds={(ids, data) => onChange((data ?? []).filter((it) => ids.includes(it.id)))}
        />
      </Block>
      <Modal open={open} centered onOk={handleRequestAll} onCancel={() => setOpen(false)}>
        <Message content={`This operation may take a long time, as it iterates through all the ${config.title}.`} />
      </Modal>
    </>
  );
};
