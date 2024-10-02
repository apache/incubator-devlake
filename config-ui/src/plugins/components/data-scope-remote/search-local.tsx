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

import { useState, useReducer } from 'react';
import { CheckCircleFilled, SearchOutlined } from '@ant-design/icons';
import { Space, Tag, Button, Input, Modal } from 'antd';
import { useDebounce } from '@mints/hooks';
import type { IDType } from '@mints/miller-columns';
import { MillerColumns } from '@mints/miller-columns';

import API from '@/api';
import { Block, Loading, Message } from '@/components';
import type { IPluginConfig } from '@/types';

import * as S from './styled';

type StateType = {
  status: string;
  scope: any[];
};

const reducer = (state: StateType, action: { type: string; payload?: Pick<Partial<StateType>, 'scope'> }) => {
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

  const [{ status, scope }, dispatch] = useReducer(reducer, {
    status: 'idle',
    scope: [],
  });

  const searchDebounce = useDebounce(search, { wait: 500 });

  const request = async (groupId?: string | number, params?: any) => {
    const res = await API.scope.remote(plugin, connectionId, {
      groupId: groupId ?? null,
      pageToken: params?.nextPageToken,
    });

    return {
      data: res.children.map((it) => ({
        parentId: it.parentId,
        id: it.id,
        title: it.name ?? it.fullName,
        canExpand: it.type === 'group',
        original: it,
      })),
      hasMore: !!res.nextPageToken,
      params: {
        nextPageToken: res.nextPageToken,
      },
    };
  };

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
        original: it,
      }));
      dispatch({ type: 'APPEND', payload: { scope: data } });

      if (res.nextPageToken) {
        await getData(groupId, res.nextPageToken);
      }

      await Promise.all(data.filter((it) => it.canExpand).map((it) => getData(it.id)));
    };

    await getData();

    dispatch({ type: 'DONE' });
  };

  const millerColumnsProps = {
    bordered: true,
    theme: {
      colorPrimary: '#7497f7',
      borderColor: '#dbe4fd',
    },
    columnHeight: 300,
    mode,
    renderTitle: (id?: IDType) =>
      !id &&
      config.millerColumn?.firstColumnTitle && <S.ColumnTitle>{config.millerColumn.firstColumnTitle}</S.ColumnTitle>,
    renderLoading: () => <Loading size={20} style={{ padding: '4px 12px' }} />,
    selectable: true,
    disabledIds: disabledScope.map((it) => it.id),
    selectedIds: selectedScope.map((it) => it.id),
    onSelectedIds: (_: IDType[], data?: any) => onChange(data ?? []),
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
        {status === 'loading' && (
          <S.JobLoad>
            <Loading style={{ marginRight: 8 }} size={20} />
            Loading: <span className="count">{scope.filter((sc) => !sc.canExpand).length}</span> scopes found
          </S.JobLoad>
        )}

        {status === 'done' && (
          <S.JobLoad>
            <CheckCircleFilled style={{ color: '#4DB764' }} />
            <span className="count">{scope.filter((sc) => !sc.canExpand).length}</span> scopes found
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
        {status === 'idle' ? (
          <MillerColumns
            {...millerColumnsProps}
            request={request}
            columnCount={config.millerColumn?.columnCount ?? 1}
          />
        ) : (
          <>
            <Input prefix={<SearchOutlined />} value={search} onChange={(e) => setSearch(e.target.value)} />
            <MillerColumns
              {...millerColumnsProps}
              loading={status === 'loading'}
              items={
                searchDebounce
                  ? scope
                      .filter((it) => it.title.includes(searchDebounce) && !it.canExpand)
                      .map((it) => ({
                        ...it,
                        parentId: null,
                      }))
                  : scope
              }
            />
          </>
        )}
      </Block>
      <Modal open={open} centered onOk={handleRequestAll} onCancel={() => setOpen(false)}>
        <Message content={`This operation may take a long time, as it iterates through all the ${config.title}.`} />
      </Modal>
    </>
  );
};
