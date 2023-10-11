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
import { Button, InputGroup, Icon, Intent } from '@blueprintjs/core';
import type { McsID, McsItem, McsColumn } from 'miller-columns-select';
import { MillerColumnsSelect } from 'miller-columns-select';
import { useDebounce } from 'ahooks';

import { FormItem, MultiSelector, Loading, Dialog, Message } from '@/components';
import { PluginConfigType } from '@/plugins';

import * as T from './types';
import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  config: PluginConfigType['dataScope'];
  disabledScope: any[];
  selectedScope: any[];
  onChange: (selectedScope: any[]) => void;
}

let canceling = false;

export const SearchLocal = ({ plugin, connectionId, config, disabledScope, selectedScope, onChange }: Props) => {
  const [miller, setMiller] = useState<{
    items: McsItem<T.ResItem>[];
    loadedIds: ID[];
    expandedIds: ID[];
    nextTokenMap: Record<ID, string>;
  }>({
    items: [],
    loadedIds: [],
    expandedIds: [],
    nextTokenMap: {},
  });

  const [isOpen, setIsOpen] = useState(false);
  const [status, setStatus] = useState('init');

  const [query, setQuery] = useState('');
  const search = useDebounce(query, { wait: 500 });

  const scopes = useMemo(
    () =>
      search
        ? miller.items
            .filter((it) => it.name.toLocaleLowerCase().includes(search.toLocaleLowerCase()))
            .filter((it) => it.type !== 'group')
            .map((it) => ({
              ...it,
              parentId: null,
            }))
        : miller.items,
    [search, miller.items],
  );

  const getItems = async ({
    groupId,
    currentPageToken,
    loadAll,
  }: {
    groupId: ID | null;
    currentPageToken?: string;
    loadAll?: boolean;
  }) => {
    if (canceling) {
      canceling = false;
      setStatus('init');
      return;
    }

    const res = await API.getRemoteScope(plugin, connectionId, {
      groupId,
      pageToken: currentPageToken,
    });

    const newItems = (res.children ?? []).map((it) => ({
      ...it,
      title: it.name,
    }));

    if (res.nextPageToken) {
      setMiller((m) => ({
        ...m,
        items: [...m.items, ...newItems],
        expandedIds: [...m.expandedIds, groupId ?? 'root'],
        nextTokenMap: {
          ...m.nextTokenMap,
          [`${groupId ? groupId : 'root'}`]: res.nextPageToken,
        },
      }));

      if (loadAll) {
        await getItems({ groupId, currentPageToken: res.nextPageToken, loadAll });
      }
    } else {
      setMiller((m) => ({
        ...m,
        items: [...m.items, ...newItems],
        expandedIds: [...m.expandedIds, groupId ?? 'root'],
        loadedIds: [...m.loadedIds, groupId ?? 'root'],
      }));

      const groupItems = newItems.filter((it) => it.type === 'group');

      if (loadAll && groupItems.length) {
        groupItems.forEach(async (it) => await getItems({ groupId: it.id, loadAll: true }));
      }
    }
  };

  useEffect(() => {
    getItems({ groupId: null });
  }, []);

  useEffect(() => {
    if (
      miller.items.length &&
      !miller.items.filter((it) => it.type === 'group' && !miller.loadedIds.includes(it.id)).length
    ) {
      setStatus('loaded');
    }
  }, [miller]);

  const handleLoadAllScopes = async () => {
    setIsOpen(false);
    setStatus('loading');

    if (!miller.loadedIds.includes('root')) {
      await getItems({
        groupId: null,
        currentPageToken: miller.nextTokenMap['root'],
        loadAll: true,
      });
    }

    const noLoadedItems = miller.items.filter((it) => it.type === 'group' && !miller.loadedIds.includes(it.id));
    if (noLoadedItems.length) {
      noLoadedItems.forEach(async (it) => {
        await getItems({
          groupId: it.id,
          currentPageToken: miller.nextTokenMap[it.id],
          loadAll: true,
        });
      });
    }
  };

  const handleCancelLoadAllScopes = () => {
    setStatus('cancel');
    canceling = true;
  };

  return (
    <S.Wrapper>
      <FormItem label={config.title} required>
        <MultiSelector
          disabled
          items={selectedScope}
          getKey={(it) => it.id}
          getName={(it) => it.fullName}
          selectedItems={selectedScope}
        />
      </FormItem>
      <FormItem>
        {(status === 'loading' || status === 'cancel') && (
          <S.JobLoad>
            <Loading style={{ marginRight: 8 }} size={20} />
            Loading: <span className="count">{miller.items.length}</span> scopes found
            <Button
              style={{ marginLeft: 8 }}
              loading={status === 'cancel'}
              small
              text="Cancel"
              onClick={handleCancelLoadAllScopes}
            />
          </S.JobLoad>
        )}

        {status === 'loaded' && (
          <S.JobLoad>
            <Icon icon="endorsed" style={{ color: '#4DB764' }} />
            <span className="count">{miller.items.length}</span> scopes found
          </S.JobLoad>
        )}

        {status === 'init' && (
          <S.JobLoad>
            <Button
              disabled={!miller.items.length}
              intent={Intent.PRIMARY}
              text="Load all scopes to search by keywords"
              onClick={() => setIsOpen(true)}
            />
          </S.JobLoad>
        )}
      </FormItem>
      <FormItem>
        {status === 'loaded' && (
          <InputGroup leftIcon="search" value={query} onChange={(e) => setQuery(e.target.value)} />
        )}
        <MillerColumnsSelect
          items={scopes}
          columnCount={search ? 1 : config.millerColumn?.columnCount ?? 1}
          columnHeight={300}
          getCanExpand={(it) => it.type === 'group'}
          getHasMore={(id) => !miller.loadedIds.includes(id ?? 'root')}
          onExpand={(id: McsID) => getItems({ groupId: id })}
          onScroll={(id: McsID | null) =>
            getItems({ groupId: id, currentPageToken: miller.nextTokenMap[id ?? 'root'] })
          }
          renderTitle={(column: McsColumn) =>
            !column.parentId &&
            config.millerColumn?.firstColumnTitle && (
              <S.ColumnTitle>{config.millerColumn.firstColumnTitle}</S.ColumnTitle>
            )
          }
          renderLoading={() => <Loading size={20} style={{ padding: '4px 12px' }} />}
          disabledIds={(disabledScope ?? []).map((it) => it.id)}
          selectedIds={selectedScope.map((it) => it.id)}
          onSelectItemIds={(selectedIds: ID[]) => onChange(miller.items.filter((it) => selectedIds.includes(it.id)))}
          expandedIds={miller.expandedIds}
        />
      </FormItem>
      <Dialog isOpen={isOpen} okText="Load" onCancel={() => setIsOpen(false)} onOk={handleLoadAllScopes}>
        <Message content={`This operation may take a long time, as it iterates through all the ${config.title}.`} />
      </Dialog>
    </S.Wrapper>
  );
};
