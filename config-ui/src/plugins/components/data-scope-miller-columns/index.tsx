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
import type { McsID, McsItem, McsColumn } from 'miller-columns-select';
import MillerColumnsSelect from 'miller-columns-select';

import { Loading } from '@/components';

import type { ExtraType } from './types';
import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  selectedItems?: any;
  onChangeItems?: (selectedItems: any) => void;
}

export const DataScopeMillerColumns = ({ plugin, connectionId, ...props }: Props) => {
  const [items, setItems] = useState<McsItem<ExtraType>[]>([]);
  const [selectedIds, setSelectedIds] = useState<ID[]>([]);
  const [loadedIds, setLoadedIds] = useState<ID[]>([]);
  const [nextTokenMap, setNextTokenMap] = useState<Record<ID, string>>({});

  useEffect(() => {
    setSelectedIds((props.selectedItems ?? []).map((it: any) => it.id));
  }, [props.selectedItems]);

  const getItems = async (groupId: ID | null, pageToken?: string) => {
    const res = await API.getScope(plugin, connectionId, {
      groupId,
      pageToken,
    });

    setItems([
      ...items,
      ...res.children.map((it: any) => ({
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

  const handleChangeItems = (selectedIds: ID[]) => {
    const result = items.filter((it) => (selectedIds.length ? selectedIds.includes(it.id) : false));
    props.onChangeItems ? props.onChangeItems(result.map((it) => it.data)) : setSelectedIds(selectedIds);
  };

  const handleExpand = (id: McsID) => getItems(id, nextTokenMap[id]);

  const handleScroll = (id: McsID | null) => getItems(id, nextTokenMap[id ?? 'root']);

  const renderTitle = (column: McsColumn) => {
    return !column.parentId && <S.ColumnTitle>Subgroups/Projects</S.ColumnTitle>;
  };

  const renderLoading = () => {
    return <Loading size={20} style={{ padding: '4px 12px' }} />;
  };

  return (
    <MillerColumnsSelect
      items={items}
      getCanExpand={(it) => it.type === 'group'}
      getHasMore={(id) => !loadedIds.includes(id ?? 'root')}
      onExpand={handleExpand}
      onScroll={handleScroll}
      columnCount={2.5}
      columnHeight={300}
      renderTitle={renderTitle}
      renderLoading={renderLoading}
      selectedIds={selectedIds}
      onSelectItemIds={handleChangeItems}
    />
  );
};
