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

import { useEffect, useMemo, useState } from 'react';
import { Button, Intent, InputGroup } from '@blueprintjs/core';
import type { McsID, McsItem, McsColumn } from 'miller-columns-select';
import MillerColumnsSelect from 'miller-columns-select';
import { useDebounce } from 'ahooks';
import { uniqBy } from 'lodash';

import { FormItem, MultiSelector, Loading, Buttons } from '@/components';
import { getPluginConfig, getPluginScopeId } from '@/plugins';
import { operator } from '@/utils';

import * as T from './types';
import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  disabledScope?: any[];
  onCancel: () => void;
  onSubmit: (origin: any) => void;
}

export const DataScopeSelectRemote = ({ plugin, connectionId, disabledScope, onCancel, onSubmit }: Props) => {
  const [operating, setOperating] = useState(false);
  const [selectedScope, setSelectedScope] = useState<T.ResItem[]>([]);

  const config = useMemo(() => getPluginConfig(plugin).dataScope, [plugin]);

  const handleSubmit = async () => {
    const [success, res] = await operator(
      () => API.updateDataScope(plugin, connectionId, { data: selectedScope.map((it) => it.data) }),
      {
        setOperating,
        formatMessage: () => 'Add data scope successful.',
      },
    );

    if (success) {
      onSubmit(res);
    }
  };

  return (
    <>
      {config.render ? (
        config.render({
          connectionId,
          disabledItems: disabledScope?.map((it) => ({ id: getPluginScopeId(plugin, it) })),
          selectedItems: selectedScope,
          onChangeSelectedItems: setSelectedScope,
        })
      ) : (
        <SelectRemote
          plugin={plugin}
          connectionId={connectionId}
          config={config}
          disabledScope={disabledScope}
          selectedScope={selectedScope}
          onChangeSelectedScope={setSelectedScope}
        />
      )}
      <Buttons position="bottom" align="right">
        <Button outlined intent={Intent.PRIMARY} text="Cancel" disabled={operating} onClick={onCancel} />
        <Button
          outlined
          intent={Intent.PRIMARY}
          text="Save"
          loading={operating}
          disabled={!selectedScope.length}
          onClick={handleSubmit}
        />
      </Buttons>
    </>
  );
};

const SelectRemote = ({
  plugin,
  connectionId,
  config,
  disabledScope,
  selectedScope,
  onChangeSelectedScope,
}: Omit<Props, 'onCancel' | 'onSubmit'> & {
  config: any;
  selectedScope: any[];
  onChangeSelectedScope: (selectedScope: any[]) => void;
}) => {
  // miller columns
  const [items, setItems] = useState<McsItem<T.ResItem>[]>([]);
  const [loadedIds, setLoadedIds] = useState<ID[]>([]);
  const [nextTokenMap, setNextTokenMap] = useState<Record<ID, string>>({});

  // search
  const [query, setQuery] = useState('');
  const [items2, setItems2] = useState<McsItem<T.ResItem>[]>([]);
  const [page, setPage] = useState(1);
  const [total] = useState(0);

  const search = useDebounce(query, { wait: 500 });

  const allItems = useMemo(() => uniqBy([...items, ...items2], 'id'), [items, items2]);

  const getItems = async (groupId: ID | null, currentPageToken?: string) => {
    const res = await API.getRemoteScope(plugin, connectionId, {
      groupId,
      pageToken: currentPageToken,
    });

    setItems([
      ...items,
      ...(res.children ?? []).map((it: any) => ({
        ...it,
        title: it.name,
      })),
    ]);

    if (!res.nextPageToken) {
      setLoadedIds([...loadedIds, groupId ? groupId : 'root']);
    } else {
      setNextTokenMap({
        ...nextTokenMap,
        [`${groupId ? groupId : 'root'}`]: res.nextPageToken,
      });
    }
  };

  useEffect(() => {
    getItems(null);
  }, []);

  const searchItems = async () => {
    if (!search) return;

    const res = await API.searchRemoteScope(plugin, connectionId, {
      search,
      page,
      pageSize: 50,
    });

    setItems2(
      res.children.map((it) => ({
        ...it,
        title: it.fullName,
      })),
    );

    if (page === 1) {
      // setTotal(res.count);
    }
  };

  useEffect(() => {
    searchItems();
  }, [search, page]);

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
        <InputGroup leftIcon="search" value={query} onChange={(e) => setQuery(e.target.value)} />
        {!search ? (
          <MillerColumnsSelect
            items={items}
            columnCount={config.millerColumnCount ?? 1}
            columnHeight={300}
            getCanExpand={(it) => it.type === 'group'}
            getHasMore={(id) => !loadedIds.includes(id ?? 'root')}
            onExpand={(id: McsID) => getItems(id, nextTokenMap[id])}
            onScroll={(id: McsID | null) => getItems(id, nextTokenMap[id ?? 'root'])}
            renderTitle={(column: McsColumn) =>
              !column.parentId && config.millerFirstTitle && <S.ColumnTitle>{config.millerFirstTitle}</S.ColumnTitle>
            }
            renderLoading={() => <Loading size={20} style={{ padding: '4px 12px' }} />}
            disabledIds={(disabledScope ?? []).map((it) => getPluginScopeId(plugin, it))}
            selectedIds={selectedScope.map((it) => it.id)}
            onSelectItemIds={(selectedIds: ID[]) =>
              onChangeSelectedScope(allItems.filter((it) => selectedIds.includes(it.id)))
            }
          />
        ) : (
          <MillerColumnsSelect
            items={items2}
            columnCount={1}
            columnHeight={300}
            getCanExpand={() => false}
            getHasMore={() => total === 0}
            onScroll={() => setPage(page + 1)}
            renderLoading={() => <Loading size={20} style={{ padding: '4px 12px' }} />}
            disabledIds={(disabledScope ?? []).map((it) => getPluginScopeId(plugin, it))}
            selectedIds={selectedScope.map((it) => it.id)}
            onSelectItemIds={(selectedIds: ID[]) =>
              onChangeSelectedScope(allItems.filter((it) => selectedIds.includes(it.id)))
            }
          />
        )}
      </FormItem>
    </S.Wrapper>
  );
};
