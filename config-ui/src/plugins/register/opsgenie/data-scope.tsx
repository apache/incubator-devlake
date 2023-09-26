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
import type { McsID, McsItem } from 'miller-columns-select';
import { MillerColumnsSelect } from 'miller-columns-select';
import { useDebounce } from 'ahooks';

import { FormItem, MultiSelector, Loading, Dialog, Message } from '@/components';
import * as T from '@/plugins/components/data-scope-select-remote/types';
import * as API from '@/plugins/components/data-scope-select-remote/api';

import * as S from './styled';

interface Props {
  connectionId: ID;
  disabledItems: T.ResItem[];
  selectedItems: T.ResItem[];
  onChangeSelectedItems: (items: T.ResItem[]) => void;
}

let canceling = false;

export const DataScope = ({ connectionId, selectedItems, onChangeSelectedItems }: Props) => {
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

  const jobs = useMemo(
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

    const res = await API.getRemoteScope('jenkins', connectionId, {
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

  const handleLoadAllJobs = async () => {
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

  const handleCancelLoadAllJobs = () => {
    setStatus('cancel');
    canceling = true;
  };

  return (
    <S.DataScope>
      <FormItem label="Jobs" required>
        <MultiSelector
          disabled
          items={selectedItems}
          getKey={(it) => it.id}
          getName={(it) => it.fullName}
          selectedItems={selectedItems}
        />
      </FormItem>
      <FormItem>
        {(status === 'loading' || status === 'cancel') && (
          <S.JobLoad>
            <Loading style={{ marginRight: 8 }} size={20} />
            Loading: <span className="count">{miller.items.length}</span> jobs found
            <Button
              style={{ marginLeft: 8 }}
              loading={status === 'cancel'}
              small
              text="Cancel"
              onClick={handleCancelLoadAllJobs}
            />
          </S.JobLoad>
        )}

        {status === 'loaded' && (
          <S.JobLoad>
            <Icon icon="endorsed" style={{ color: '#4DB764' }} />
            <span className="count">{miller.items.length}</span> jobs found
          </S.JobLoad>
        )}

        {status === 'init' && (
          <S.JobLoad>
            <Button
              disabled={!miller.items.length}
              intent={Intent.PRIMARY}
              text="Load all jobs to search by keywords"
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
          items={jobs}
          columnCount={search ? 1 : 2.5}
          columnHeight={300}
          getCanExpand={(it) => it.type === 'group'}
          getHasMore={(id) => !miller.loadedIds.includes(id ?? 'root')}
          onExpand={(id: McsID) => getItems({ groupId: id })}
          onScroll={(id: McsID | null) =>
            getItems({ groupId: id, currentPageToken: miller.nextTokenMap[id ?? 'root'] })
          }
          renderLoading={() => <Loading size={20} style={{ padding: '4px 12px' }} />}
          selectedIds={selectedItems.map((it) => it.id)}
          onSelectItemIds={(selectedIds: ID[]) =>
            onChangeSelectedItems(miller.items.filter((it) => selectedIds.includes(it.id)))
          }
          expandedIds={miller.expandedIds}
          // onChangeExpandedIds={(expandedIds: ID[]) => setExpandedIds(expandedIds)}
        />
      </FormItem>
      <Dialog isOpen={isOpen} okText="Load" onCancel={() => setIsOpen(false)} onOk={handleLoadAllJobs}>
        <Message content="This operation may take a long time, as it iterates through all the Jenkins Jobs." />
      </Dialog>
    </S.DataScope>
  );
};
