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

import { useState, useEffect } from 'react';
import type { McsItem } from 'miller-columns-select';
import MillerColumnsSelect from 'miller-columns-select';

import { Dialog, FormItem, Loading } from '@/components';

import * as T from '../types';
import * as API from '../api';
import * as S from '../styled';

interface Props {
  isOpen: boolean;
  saving: boolean;
  onCancel: () => void;
  onSubmit: (items: T.WebhookItemType[]) => void;
}

export const SelectorDialog = ({ isOpen, saving, onCancel, onSubmit }: Props) => {
  const [items, setItems] = useState<McsItem<T.WebhookItemType>[]>([]);
  const [selectedItems, setSelectedItems] = useState<T.WebhookItemType[]>([]);
  const [isLast, setIsLast] = useState(false);

  const updateItems = (arr: any) =>
    arr.map((it: any) => ({
      parentId: null,
      id: it.id,
      title: it.name,
      name: it.name,
    }));

  useEffect(() => {
    (async () => {
      const res = await API.getConnections();
      setItems([...updateItems(res)]);
      setIsLast(true);
    })();
  }, []);

  const handleSubmit = () => onSubmit(selectedItems);

  return (
    <Dialog
      isOpen={isOpen}
      title="Select Existing Webhooks"
      style={{
        width: 820,
      }}
      okText="Confrim"
      okLoading={saving}
      okDisabled={!selectedItems.length}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      <S.Wrapper>
        <FormItem label="Webhooks" subLabel="Select an existing Webhook to import to the current project.">
          <MillerColumnsSelect
            columnCount={1}
            columnHeight={160}
            getHasMore={() => !isLast}
            renderLoading={() => <Loading size={20} style={{ padding: '4px 12px' }} />}
            items={items}
            selectedIds={selectedItems.map((it) => it.id)}
            onSelectItemIds={(seletedIds: ID[]) =>
              setSelectedItems(
                items
                  .filter((it) => seletedIds.includes(it.id))
                  .map((it) => ({
                    id: it.id,
                    name: it.name,
                  })),
              )
            }
          />
        </FormItem>
      </S.Wrapper>
    </Dialog>
  );
};

export default SelectorDialog;
