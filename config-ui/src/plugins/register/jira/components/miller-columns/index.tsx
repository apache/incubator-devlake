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
import type { McsItem } from 'miller-columns-select';
import MillerColumnsSelect from 'miller-columns-select';

import { Loading } from '@/components';

import { ScopeItemType } from '../../types';

import type { UseMillerColumnsProps } from './use-miller-columns';
import { useMillerColumns } from './use-miller-columns';

interface Props extends UseMillerColumnsProps {
  selectedItems: ScopeItemType[];
  onChangeItems: (selectedItems: ScopeItemType[]) => void;
}

export const MillerColumns = ({ connectionId, selectedItems, onChangeItems }: Props) => {
  const [selectedIds, setSelectedIds] = useState<ID[]>([]);

  const { items, getHasMore, onScroll } = useMillerColumns({
    connectionId,
  });

  useEffect(() => {
    setSelectedIds(selectedItems.map((it) => it.boardId));
  }, [selectedItems]);

  const handleChangeItems = (selectedIds: ID[]) => {
    const result = selectedIds.map((id) => {
      const selectedItem = selectedItems.find((it) => it.boardId === id);
      if (selectedItem) {
        return selectedItem;
      }

      const item = items.find((it) => it.id === id) as McsItem<ScopeItemType>;
      return {
        connectionId,
        boardId: item.boardId,
        name: item.name,
        self: item.self,
        type: item.type,
        projectId: item.projectId,
      };
    });

    onChangeItems(result);
  };

  const renderLoading = () => {
    return <Loading size={20} style={{ padding: '4px 12px' }} />;
  };

  return (
    <MillerColumnsSelect
      items={items}
      getHasMore={getHasMore}
      onScroll={onScroll}
      columnCount={1}
      columnHeight={300}
      renderLoading={renderLoading}
      selectedIds={selectedIds}
      onSelectItemIds={handleChangeItems}
    />
  );
};
